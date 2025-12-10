package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/config"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

const (
	readHeaderTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer

	server *http.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	results := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := a.runHttpServer(ctx)
		if err != nil {
			results <- err
		}
	}()

	go func() {
		err := a.runKafkaConsumer(ctx)
		if err != nil {
			results <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")

	case err := <-results:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))

		cancel()
		<-ctx.Done()

		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context, cfg *config.Config) error {
	inits := []func(ctx context.Context, cfg *config.Config) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHttpServer,
	}

	for _, init := range inits {
		err := init(ctx, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("👂 HTTP-сервер запущен на порту %s\n", a.diContainer.config.Http.Port()))
	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❗ Ошибка при запуске HTTP-сервера: \n", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) runKafkaConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Запуск консьюмера Kafka")
	return a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
}

func (a *App) initDI(ctx context.Context, cfg *config.Config) error {
	a.diContainer = NewDIContainer(cfg)
	return nil
}

func (a *App) initLogger(ctx context.Context, cfg *config.Config) error {
	return logger.Init(cfg.Logger.Level(), cfg.Logger.AsJson())
}

func (a *App) initCloser(ctx context.Context, cfg *config.Config) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHttpServer(ctx context.Context, cfg *config.Config) error {
	a.server = &http.Server{
		Addr:              cfg.Http.Address(),
		Handler:           a.diContainer.Router(ctx),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.server.Shutdown(ctx)
	})

	return nil
}
