package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/mocks"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

func (s *RepositorySuite) TestGetPartMainFlow() {
	// Создаем тестовую часть
	testPart := s.createTestPart()

	// Добавляем часть в репозиторий
	singleResult := mocks.NewMongoSingleResult(s.T())
	singleResult.On("Decode", mock.AnythingOfType("*model.Part")).Run(func(args mock.Arguments) {
		// Получаем указатель на структуру и заполняем его данными
		part := args.Get(0).(*repoModel.Part)
		*part = *testPart
	}).Return(nil)

	s.collection.On("FindOne", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		return v["uuid"] == testPart.Uuid
	})).Return(singleResult)

	// Вызываем метод
	result, err := s.repository.GetPart(s.ctx, testPart.Uuid)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), testPart.Uuid, result.Uuid)
	require.Equal(s.T(), testPart.Name, result.Name)
	require.Equal(s.T(), testPart.Description, result.Description)
	require.Equal(s.T(), testPart.Price, result.Price)
	require.Equal(s.T(), testPart.StockQuantity, result.StockQuantity)
	require.Equal(s.T(), model.Category(testPart.Category), result.Category)
	require.Equal(s.T(), testPart.Tags, result.Tags)
	require.Equal(s.T(), testPart.Metadata, result.Metadata)
	require.NotNil(s.T(), result.Dimensions)
	require.Equal(s.T(), testPart.Dimensions.Length, result.Dimensions.Length)
	require.Equal(s.T(), testPart.Dimensions.Width, result.Dimensions.Width)
	require.Equal(s.T(), testPart.Dimensions.Height, result.Dimensions.Height)
	require.Equal(s.T(), testPart.Dimensions.Weight, result.Dimensions.Weight)
	require.NotNil(s.T(), result.Manufacturer)
	require.Equal(s.T(), testPart.Manufacturer.Name, result.Manufacturer.Name)
	require.Equal(s.T(), testPart.Manufacturer.Country, result.Manufacturer.Country)
	require.Equal(s.T(), testPart.Manufacturer.Website, result.Manufacturer.Website)
}

func (s *RepositorySuite) TestGetPartNotFound() {
	// Используем UUID, которого нет в репозитории
	nonExistentUuid := gofakeit.UUID()

	singleResult := mocks.NewMongoSingleResult(s.T())
	singleResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments)

	s.collection.On("FindOne", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		return v["uuid"] == nonExistentUuid
	})).Return(singleResult)

	// Вызываем метод
	result, err := s.repository.GetPart(s.ctx, nonExistentUuid)

	// Проверяем результат
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrPartNotFound)
	require.Nil(s.T(), result)
}

func (s *RepositorySuite) TestGetPartEmptyUuid() {
	singleResult := mocks.NewMongoSingleResult(s.T())
	singleResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments)

	s.collection.On("FindOne", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		return v["uuid"] == ""
	})).Return(singleResult)

	// Вызываем метод с пустым UUID
	result, err := s.repository.GetPart(s.ctx, "")

	// Проверяем результат
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrPartNotFound)
	require.Nil(s.T(), result)
}

func (s *RepositorySuite) TestGetPartMinimalData() {
	// Создаем минимальную тестовую часть
	testPart := s.createMinimalTestPart()

	singleResult := mocks.NewMongoSingleResult(s.T())
	singleResult.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		part := args.Get(0).(*repoModel.Part)
		*part = *testPart
	}).Return(nil)

	s.collection.On("FindOne", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		return v["uuid"] == testPart.Uuid
	})).Return(singleResult)

	// Вызываем метод
	result, err := s.repository.GetPart(s.ctx, testPart.Uuid)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), testPart.Uuid, result.Uuid)
	require.Equal(s.T(), testPart.Name, result.Name)
	require.Equal(s.T(), testPart.Description, result.Description)
	require.Equal(s.T(), testPart.Price, result.Price)
	require.Equal(s.T(), testPart.StockQuantity, result.StockQuantity)
	require.Equal(s.T(), model.Category(testPart.Category), result.Category)
	require.Equal(s.T(), testPart.Tags, result.Tags)
	require.Equal(s.T(), testPart.Metadata, result.Metadata)
	require.Nil(s.T(), result.Dimensions)
	require.Nil(s.T(), result.Manufacturer)
}
