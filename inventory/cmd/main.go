package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/app"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/config"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

const (
	gracefulShutdownTimeout = 5 * time.Second

	configPath = "deploy/compose/inventory/.env"
)

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("❗ Failed to load env file: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx, cfg)
	if err != nil {
		panic(fmt.Errorf("❗ Failed to create app: %w", err))
	}

	err = a.Run(appCtx)
	if err != nil {
		panic(fmt.Errorf("❗ Failed to run app: %w", err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❗ failed to close all: %v", zap.Error(err))
	}
}
