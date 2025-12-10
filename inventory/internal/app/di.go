package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/api/inventory/v1"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/config"
	repository "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository"
	partRepository "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/part"
	service "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/service"
	partService "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/service/part"
	closer "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	config *config.Config

	inventoryV1API genInventoryV1.InventoryServiceServer
	partService    service.PartService
	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDIContainer(cfg *config.Config) *diContainer {
	return &diContainer{config: cfg}
}

func (c *diContainer) InventoryV1API(ctx context.Context) genInventoryV1.InventoryServiceServer {
	if c.inventoryV1API == nil {
		c.inventoryV1API = inventoryV1.NewAPI(c.PartService(ctx))
	}

	return c.inventoryV1API
}

func (c *diContainer) PartService(ctx context.Context) service.PartService {
	if c.partService == nil {
		c.partService = partService.NewService(c.PartRepository(ctx))

		// Init with dummy here
		err := c.partService.InitWithDummy(ctx)
		if err != nil {
			panic(fmt.Sprintf("❗ failed to init with dummy data: %v\n", err.Error()))
		}
	}

	return c.partService
}

func (c *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if c.partRepository == nil {
		var err error
		c.partRepository, err = partRepository.NewRepository(ctx, &repository.MongoDatabaseAdapter{Database: c.MongoDBHandle(ctx)})
		if err != nil {
			panic(fmt.Sprintf("❗ failed to create repository: %v\n", err.Error()))
		}
	}

	return c.partRepository
}

func (c *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if c.mongoDBHandle == nil {
		c.mongoDBHandle = c.MongoDBClient(ctx).Database(c.config.MongoConfig.DatabaseName())
	}
	return c.mongoDBHandle
}

func (c *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if c.mongoDBClient == nil {
		var err error
		c.mongoDBClient, err = mongo.Connect(ctx, options.Client().ApplyURI(c.config.MongoConfig.URI()))
		if err != nil {
			panic(fmt.Sprintf("❗ failed to create client: %v\n", err.Error()))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return c.mongoDBClient.Disconnect(ctx)
		})

		if err = c.mongoDBClient.Ping(ctx, nil); err != nil {
			panic(fmt.Sprintf("❗ failed to ping client: %v\n", err.Error()))
		}
	}
	return c.mongoDBClient
}
