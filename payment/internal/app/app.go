package app

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Mahno9/GoMicroservicesCourse/payment/internal/config"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context, cfg *config.Config) error {
	inits := []func(ctx context.Context, cfg *config.Config) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
	}

	for _, initFunc := range inits {
		err := initFunc(ctx, cfg)
		if err != nil {
			return err
		}
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

func (a *App) initListener(ctx context.Context, cfg *config.Config) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GrpcConfig.Host(), cfg.GrpcConfig.Port()))
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		return listener.Close()
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context, cfg *config.Config) error {
	a.grpcServer = grpc.NewServer()
	closer.AddNamed("GRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	genPaymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("👂 gRPC server listening on port %s", a.diContainer.config.GrpcConfig.Port()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
