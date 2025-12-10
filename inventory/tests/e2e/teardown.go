package e2e

import (
	"context"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	logger.Info(ctx, "🧹 Cleaning up test environment")

	cleanupTestEnvironment(ctx, env)

	logger.Info(ctx, "✅ Test environment cleaned up successfully")
}

func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "🛑 Unable to terminate application container", zap.Error(err))
		} else {
			logger.Info(ctx, "✔ Application container terminated successfully")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "🛑 Unable to terminate MongoDB container", zap.Error(err))
		} else {
			logger.Info(ctx, "✔ MongoDB container terminated successfully")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "🛑 Unable to remove network", zap.Error(err))
		} else {
			logger.Info(ctx, "✔ Network removed successfully")
		}
	}
}
