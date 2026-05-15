package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/cyradin/govnocode/internal/config"
	"github.com/cyradin/govnocode/internal/container"
	"github.com/cyradin/govnocode/pkg/logger"
)

var GitCommit string = "dev"

var (
	projectRootFlag = flag.String(
		"project-root",
		".",
		"Path to the project root directory where the agent will run (must exist)",
	)

	taskFlag = flag.String(
		"task",
		"",
		"Task description for the agent (required)",
	)
)

func main() {
	flag.Parse()

	if err := validateFlags(); err != nil {
		log.Fatalf("invalid flags: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	container := container.New(GitCommit, cfg)

	if err := run(cfg, container); err != nil {
		container.Logger().Error("application error", logger.Error(err))
	}
}

func run(_ *config.Config, container *container.Container) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx = logger.WithContext(ctx, container.Logger())

	absPath, err := filepath.Abs(*projectRootFlag)
	if err != nil {
		return fmt.Errorf("get absolute path: %w", err)
	}

	errCh := make(chan error, 1)

	agent, err := container.GoAgent().WithWorkdir(absPath).Build(ctx)
	if err != nil {
		return fmt.Errorf("init go agent: %w", err)
	}

	go func() {
		logger.FromContext(ctx).Info("running code agent", slog.String("project_root", absPath))

		errCh <- agent.Start(ctx, absPath, *taskFlag)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func validateFlags() error {
	if *taskFlag == "" {
		return fmt.Errorf("flag -task is required")
	}

	info, err := os.Stat(*projectRootFlag)
	if err != nil {
		if os.IsNotExist(err) {
			//nolint:mnd
			if err := os.MkdirAll(*projectRootFlag, 0750); err != nil {
				return fmt.Errorf("unable to create directory: %w", err)
			}

			return nil
		}

		return fmt.Errorf("cannot access project root: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("project root is not a directory: %s", *projectRootFlag)
	}

	return nil
}
