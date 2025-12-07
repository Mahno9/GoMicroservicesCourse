package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/config"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

type app struct {
	diContainer *diContainer
}

func New(ctx context.Context, cfg *config.Config) (*app, error) {
	a := &app{}

	err := a.initDeps(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *app) Run(ctx context.Context) error {
	result := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := a.runKafkaConsumer(ctx)
		if err != nil {
			result <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-result:
		logger.Error(ctx, "Kafka consumer crashed, shutting down", zap.Error(err))

		cancel()
		<-ctx.Done()

		return err
	}

	return nil
}

func (a *app) initDeps(ctx context.Context, cfg *config.Config) error {
	inits := []func(ctx context.Context, cfg *config.Config) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
	}

	for _, init := range inits {
		err := init(ctx, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *app) initDI(ctx context.Context, cfg *config.Config) error {
	a.diContainer = NewDIContainer(cfg)
	return nil
}

func (a *app) initLogger(ctx context.Context, cfg *config.Config) error {
	return logger.Init(cfg.Logger.Level(), cfg.Logger.AsJson())
}

func (a *app) initCloser(ctx context.Context, cfg *config.Config) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *app) runKafkaConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Запуск консьюмера Kafka")
	return a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
}
