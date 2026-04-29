package main

import (
	"context"
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

func run(_ *config.Config, container *container.Container) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx = logger.WithContext(ctx, container.Logger())

	errCh := make(chan error, 1)

	agent := container.CodingAgent()

	go func() {
		logger.FromContext(ctx).Info("app started")

		errCh <- agent.Start(ctx, "The task is to implement simple Go HTTP server with a few endpoints")
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
