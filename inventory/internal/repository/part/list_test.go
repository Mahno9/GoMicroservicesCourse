package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

func (s *RepositorySuite) TestListPartsMainFlow() {
	// Создаем тестовые части с разными характеристиками
	part1 := s.createTestPart()
	part1.Category = repoModel.CategoryEngine
	part1.Manufacturer.Country = "USA"
	part1.Tags = []string{"engine", "primary"}

	part2 := s.createTestPart()
	part2.Category = repoModel.CategoryFuel
	part2.Manufacturer.Country = "Germany"
	part2.Tags = []string{"fuel", "tank"}

	part3 := s.createTestPart()
	part3.Category = repoModel.CategoryWing
	part3.Manufacturer.Country = "Japan"
	part3.Tags = []string{"wing", "aerodynamic"}

	// Добавляем части в репозиторий
	s.repository.Add(part1)
	s.repository.Add(part2)
	s.repository.Add(part3)

	// Создаем фильтр - ищем части, которые соответствуют любому из UUID
	filter := &model.PartsFilter{
		Uuids: []string{
			part1.Uuid,
			part2.Uuid,
		},
	}

	// Вызываем метод
	result, err := s.repository.ListParts(s.ctx, filter)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.Len(s.T(), result, 2) // Должны вернуть только part1 и part2

	// Проверяем, что результаты соответствуют фильтру
	resultUuids := make(map[string]bool)
	for _, part := range result {
		resultUuids[part.Uuid] = true
	}
	require.True(s.T(), resultUuids[part1.Uuid])
	require.True(s.T(), resultUuids[part2.Uuid])
	require.False(s.T(), resultUuids[part3.Uuid])
}

func (s *RepositorySuite) TestListPartsEmptyFilter() {
	// Создаем тестовые части
	part1 := s.createTestPart()
	part2 := s.createTestPart()

	// Добавляем части в репозиторий
	s.repository.Add(part1)
	s.repository.Add(part2)

	// Создаем пустой фильтр
	filter := &model.PartsFilter{}

	// Вызываем метод
	result, err := s.repository.ListParts(s.ctx, filter)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.Len(s.T(), result, 2) // Должны вернуть все части
}

func (s *RepositorySuite) TestListPartsOnlyUuidsFilter() {
	// Создаем тестовые части
	part1 := s.createTestPart()
	part2 := s.createTestPart()
	part3 := s.createTestPart()

	// Добавляем части в репозиторий
	s.repository.Add(part1)
	s.repository.Add(part2)
	s.repository.Add(part3)

	// Создаем фильтр только с UUID
	filter := &model.PartsFilter{
		Uuids: []string{
			part1.Uuid,
			part3.Uuid,
		},
	}

	// Вызываем метод
	result, err := s.repository.ListParts(s.ctx, filter)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.Len(s.T(), result, 2) // Должны вернуть только part1 и part3

	// Проверяем, что результаты соответствуют фильтру
	resultUuids := make(map[string]bool)
	for _, part := range result {
		resultUuids[part.Uuid] = true
	}
	require.True(s.T(), resultUuids[part1.Uuid])
	require.False(s.T(), resultUuids[part2.Uuid])
	require.True(s.T(), resultUuids[part3.Uuid])
}

func (s *RepositorySuite) TestListPartsOnlyTagsFilter() {
	// Создаем тестовые части с разными тегами
	part1 := s.createTestPart()
	part1.Tags = []string{"engine", "primary"}

	part2 := s.createTestPart()
	part2.Tags = []string{"fuel", "tank"}

	part3 := s.createTestPart()
	part3.Tags = []string{"wing", "aerodynamic"}

	// Добавляем части в репозиторий
	s.repository.Add(part1)
	s.repository.Add(part2)
	s.repository.Add(part3)

	// Создаем фильтр только с тегами
	filter := &model.PartsFilter{
		Tags: []string{
			"primary",
			"aerodynamic",
		},
	}

	// Вызываем метод
	result, err := s.repository.ListParts(s.ctx, filter)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.Len(s.T(), result, 2) // Должны вернуть только part1 и part3

	// Проверяем, что результаты соответствуют фильтру
	resultUuids := make(map[string]bool)
	for _, part := range result {
		resultUuids[part.Uuid] = true
	}
	require.True(s.T(), resultUuids[part1.Uuid])
	require.False(s.T(), resultUuids[part2.Uuid])
	require.True(s.T(), resultUuids[part3.Uuid])
}

func (s *RepositorySuite) TestListPartsNoMatches() {
	// Создаем тестовую часть
	part1 := s.createTestPart()
	s.repository.Add(part1)

	// Создаем фильтр, который не совпадает ни с одной частью
	filter := &model.PartsFilter{
		Uuids: []string{
			gofakeit.UUID(), // Случайный UUID, которого нет в репозитории
		},
	}

	// Вызываем метод
	result, err := s.repository.ListParts(s.ctx, filter)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.Empty(s.T(), result) // Не должно быть совпадений
}

func (s *RepositorySuite) TestListPartsEmptyRepository() {
	// Репозиторий пуст (не добавляем никаких частей)

	// Создаем пустой фильтр
	filter := &model.PartsFilter{}

	// Вызываем метод
	result, err := s.repository.ListParts(s.ctx, filter)

	// Проверяем результат
	require.NoError(s.T(), err)
	require.Empty(s.T(), result) // Репозиторий пуст, результат должен быть пустым
}

func (s *RepositorySuite) TestListPartsConcurrentAccess() {
	// Создаем тестовые части
	part1 := s.createTestPart()
	part2 := s.createTestPart()

	// Добавляем части в репозиторий
	s.repository.Add(part1)
	s.repository.Add(part2)

	// Создаем фильтр
	filter := &model.PartsFilter{}

	// Запускаем несколько горутин для проверки конкурентного доступа
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			result, err := s.repository.ListParts(s.ctx, filter)
			require.NoError(s.T(), err)
			require.Len(s.T(), result, 2)
			done <- true
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}
}
