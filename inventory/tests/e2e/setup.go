package e2e

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/testcontainers"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/testcontainers/app"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/testcontainers/mongo"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/testcontainers/network"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/testcontainers/path"
)

const (
	projectName = "inventory-service"

	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	grpcPortKey = "GRPC_PORT"

	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🔰 Setting up test environment")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "🛑 Unable to create common network", zap.Error(err))
	}
	logger.Info(ctx, "✔ Created common network", zap.String("network_name", generatedNetwork.Name()))

	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "🛑 Unable to run MongoDB container", zap.Error(err))
	}
	logger.Info(ctx, "✔ MongoDB container started successfully")

	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithEnv(appEnv),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "🛑 Unable to run application container", zap.Error(err))
	}
	logger.Info(ctx, "✔ Application container started successfully")

	logger.Info(ctx, "✅ Application container started successfully")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
