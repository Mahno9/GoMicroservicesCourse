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

func (a *App) Run(ctx context.Context) error {
	return a.runHttpServer(ctx)
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("👂 HTTP-сервер запущен на порту %s\n", a.diContainer.config.HttpConfig.Port()))
	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❗ Ошибка при запуске HTTP-сервера: \n", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initDI(ctx context.Context, cfg *config.Config) error {
	a.diContainer = NewDIContainer(cfg)
	return nil
}

func (a *App) initLogger(ctx context.Context, cfg *config.Config) error {
	return logger.Init(cfg.LoggerConfig.Level(), cfg.LoggerConfig.AsJson())
}

func (a *App) initCloser(ctx context.Context, cfg *config.Config) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHttpServer(ctx context.Context, cfg *config.Config) error {
	a.server = &http.Server{
		Addr:              cfg.HttpConfig.Address(),
		Handler:           a.diContainer.Router(ctx),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return nil
}
