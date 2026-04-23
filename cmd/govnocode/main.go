package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/cyradin/govnocode/internal/config"
	"github.com/cyradin/govnocode/internal/container"
	"github.com/cyradin/govnocode/pkg/logger"
)

var GitCommit string = "dev"

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	container := container.New(GitCommit, cfg)

	if err := run(cfg, container); err != nil {
		container.Logger().Error("application error", logger.Error(err))
	}
}

func run(cfg *config.Config, container *container.Container) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)

	// @todo
	container.Logger().Info("app started")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
