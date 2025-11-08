package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *RepositorySuite) TestGetPartMainFlow() {
	// Создаем тестовую часть
	testPart := s.createTestPart()

	// Добавляем часть в репозиторий
	s.repository.Add(testPart)

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
	s.repository.Add(testPart)

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
	s.repository.Add(testPart)

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
	threadCount := 10
	uuidsPerThread := 1000
	totalUuids := threadCount * uuidsPerThread

	uuids := make([]string, totalUuids)
	for i := range uuids {
		uuids[i] = gofakeit.UUID()
	}

	cout := make(chan string, totalUuids)
	done := make(chan bool, threadCount)

	for i := range threadCount {
		go func() {
			threadUuids := make([]string, uuidsPerThread)
			copy(threadUuids, uuids[i*uuidsPerThread:(i+1)*uuidsPerThread])

			gofakeit.ShuffleStrings(threadUuids)

			for _, uuid := range threadUuids {
				for {
					v, ok := s.repository.Get(uuid)
					if ok {
						cout <- v.Uuid
						break
					}
				}
			}

			done <- true
		}()
	}

	// Запись после начала чтения
	go func() {
		time.Sleep(10 * time.Millisecond) // nolint: forbidigo
		for _, uuid := range uuids {
			testPart := s.createTestPartWithUuid(uuid)
			s.repository.Add(testPart)
		}
	}()

	// Ждём, что все uuid будут считаны
	for range totalUuids {
		uuid := <-cout
		require.True(s.T(), lo.Contains(uuids, uuid))

		uuids = lo.Filter(uuids, func(x string, _ int) bool {
			return x != uuid
		})
	}

	// Проверяем, что все uuid прошли через цикл чтения-записи
	require.Empty(s.T(), uuids)

	// Ждем завершения всех горутин
	for range threadCount {
		<-done
	}
}
