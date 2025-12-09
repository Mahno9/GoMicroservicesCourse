package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/notification/internal/app"
	"github.com/Mahno9/GoMicroservicesCourse/notification/internal/config"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

const (
	gracefulShutdownTimeout = 5 * time.Second

	configPath = "deploy/compose/notification/.env"
)

func main() {
	config, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("❗ Failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	app, err := app.New(appCtx, config)
	if err != nil {
		panic(fmt.Errorf("❗ Failed to create app: %w", err))
	}

	err = app.Run(appCtx)
	if err != nil {
		panic(fmt.Errorf("❗ Failed to run app: %w", err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	err := closer.CloseAll(ctx)
	if err != nil {
		logger.Error(ctx, "❗ failed to close all: %v", zap.Error(err))
	}
}
