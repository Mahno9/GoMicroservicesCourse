package app

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/notification/internal/config"
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
	result := make(chan error, 3)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := a.runOrderPaidConsumer(ctx)
		if err != nil {
			result <- err
		}
	}()

	go func() {
		err := a.runShipAssembledConsumer(ctx)
		if err != nil {
			result <- err
		}
	}()

	go func() {
		a.runTelegramBotClient(ctx)
		result <- errors.New("telegram bot stopped")
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-result:
		logger.Error(ctx, "One of service crashed", zap.Error(err))

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
		a.initTelegramService,
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

func (a *app) initTelegramService(ctx context.Context, cfg *config.Config) error {
	a.diContainer.TelegramService(ctx)
	return nil
}

func (a *app) runShipAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Run consumer ShipAssembled")

	err := a.diContainer.ShipAssembledConsumerService(ctx).RunConsumer(ctx)
	return err
}

func (a *app) runOrderPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Run consumer OrderPaid")

	err := a.diContainer.OrderPaidConsumerService(ctx).RunConsumer(ctx)
	return err
}

func (a *app) runTelegramBotClient(ctx context.Context) {
	logger.Info(ctx, "🤖 Telegram bot started")

	a.diContainer.TelegramBot().Start(ctx)

	logger.Info(ctx, "Telegram bot stopped")
}
