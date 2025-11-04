package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *RepositorySuite) TestGetPartMainFlow() {
	// Создаем тестовую часть
	testPart := s.createTestPart()

	// Добавляем часть в репозиторий
	s.repository.parts[testPart.Uuid] = testPart

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
	require.Equal(s.T(), domainModel.Category(testPart.Category), result.Category)
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

	// Вызываем метод
	result, err := s.repository.GetPart(s.ctx, nonExistentUuid)

	// Проверяем результат
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, domainModel.ErrPartNotFound)
	require.Nil(s.T(), result)
}

func (s *RepositorySuite) TestGetPartEmptyUuid() {
	// Создаем тестовую часть
	testPart := s.createTestPart()
	s.repository.parts[testPart.Uuid] = testPart

	// Вызываем метод с пустым UUID
	result, err := s.repository.GetPart(s.ctx, "")

	// Проверяем результат
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, domainModel.ErrPartNotFound)
	require.Nil(s.T(), result)
}

func (s *RepositorySuite) TestGetPartMinimalData() {
	// Создаем минимальную тестовую часть
	testPart := s.createMinimalTestPart()
	s.repository.parts[testPart.Uuid] = testPart

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
	require.Equal(s.T(), domainModel.Category(testPart.Category), result.Category)
	require.Equal(s.T(), testPart.Tags, result.Tags)
	require.Equal(s.T(), testPart.Metadata, result.Metadata)
	require.Nil(s.T(), result.Dimensions)
	require.Nil(s.T(), result.Manufacturer)
}

func (s *RepositorySuite) TestGetPartConcurrentAccess() {
	// Создаем тестовую часть
	testPart := s.createTestPart()
	s.repository.parts[testPart.Uuid] = testPart

	// Запускаем несколько горутин для проверки конкурентного доступа
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			result, err := s.repository.GetPart(s.ctx, testPart.Uuid)
			require.NoError(s.T(), err)
			require.NotNil(s.T(), result)
			require.Equal(s.T(), testPart.Uuid, result.Uuid)
			done <- true
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}
}
