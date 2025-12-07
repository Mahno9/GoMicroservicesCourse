package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inventory E2E Suite")
}

var _ = BeforeSuite(func() {
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("🛑 Unable to initialize logger: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	err = godotenv.Load(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
	if err != nil {
		panic(fmt.Sprintf("🛑 Unable to load environment variables: %v", err))
	}

	logger.Info(suiteCtx, "Launching test environment...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Finishing test suite")

	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}

	suiteCancel()
})
