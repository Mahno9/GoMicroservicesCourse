package part

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/mocks"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

type RepositorySuite struct {
	suite.Suite

	ctx        context.Context
	repository *repository

	collection *mocks.MongoCollection
	db         *mocks.MongoDatabase
}

func (s *RepositorySuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.Init("test", false); err != nil {
		s.T().Errorf("Failed to initialize logger: %v", err)
	}

	s.collection = mocks.NewMongoCollection(s.T())
	s.db = mocks.NewMongoDatabase(s.T())

	// Настраиваем мок для базы данных
	s.db.On("Collection", "parts").Return(s.collection)

	// Настраиваем мок для создания индекса
	indexView := mocks.NewMongoIndexView(s.T())
	indexView.On("CreateOne", mock.Anything, mock.Anything).Return("index_name", nil)
	s.collection.On("Indexes").Return(indexView)

	var err error
	s.repository, err = NewRepository(s.ctx, s.db)
	if err != nil {
		// Игнорируем ошибку создания индекса в тестах
		s.T().Logf("Ignoring index creation error: %v", err)
	}
}

func (s *RepositorySuite) SetupTest() {
	// Создаем новый пустой репозиторий для каждого теста
	s.collection = mocks.NewMongoCollection(s.T())
	s.db = mocks.NewMongoDatabase(s.T())

	// Настраиваем мок для базы данных
	s.db.On("Collection", "parts").Return(s.collection)

	// Настраиваем мок для создания индекса
	indexView := mocks.NewMongoIndexView(s.T())
	indexView.On("CreateOne", mock.Anything, mock.Anything).Return("index_name", nil)
	s.collection.On("Indexes").Return(indexView)

	var err error
	s.repository, err = NewRepository(s.ctx, s.db)
	if err != nil {
		// Игнорируем ошибку создания индекса в тестах
		s.T().Logf("Ignoring index creation error: %v", err)
	}
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

// Вспомогательные функции для создания тестовых данных
func (s *RepositorySuite) createTestPart() *repoModel.Part {
	return s.createTestPartWithUuid(gofakeit.UUID())
}

func (s *RepositorySuite) createTestPartWithUuid(uuid string) *repoModel.Part {
	now := time.Now()
	key := gofakeit.Word()
	value := any(key)

	return &repoModel.Part{
		Uuid:          uuid,
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(5),
		Price:         gofakeit.Float64Range(10.0, 1000.0),
		StockQuantity: int64(gofakeit.Number(1, 100)),
		Category:      repoModel.CategoryEngine,
		Dimensions: &repoModel.Dimensions{
			Length: gofakeit.Float64Range(1.0, 100.0),
			Width:  gofakeit.Float64Range(1.0, 100.0),
			Height: gofakeit.Float64Range(1.0, 100.0),
			Weight: gofakeit.Float64Range(1.0, 100.0),
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{gofakeit.Word(), gofakeit.Word()},
		Metadata:  map[string]any{"key1": value},
		CreatedAt: now,
		UpdatedAt: &now,
	}
}

func (s *RepositorySuite) createMinimalTestPart() *repoModel.Part {
	return &repoModel.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   "",
		Price:         0.0,
		StockQuantity: 0,
		Category:      repoModel.CategoryUnknown,
		Tags:          []string{},
		Metadata:      make(map[string]any),
	}
}
