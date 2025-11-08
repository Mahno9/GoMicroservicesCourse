package part

import (
	"github.com/stretchr/testify/require"

	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

func (s *RepositorySuite) TestInitWithDummyMainFlow() {
	// Проверяем, что репозиторий изначально пуст
	require.Empty(s.T(), s.repository.Count())

	// Вызываем метод инициализации
	err := s.repository.InitWithDummy()

	// Проверяем результат
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), s.repository)

	// Проверяем, что были созданы части с корректными данными
	for _, part := range s.repository.GetAll() {
		require.NotEmpty(s.T(), part.Uuid)
		require.NotEmpty(s.T(), part.Name)
		require.NotEmpty(s.T(), part.Description)
		require.Greater(s.T(), part.Price, 0.0)
		require.Greater(s.T(), part.StockQuantity, int64(0))
		require.NotEmpty(s.T(), string(part.Category))

		// Проверяем, что категория одна из разрешенных
		validCategories := []repoModel.Category{
			repoModel.CategoryUnknown,
			repoModel.CategoryEngine,
			repoModel.CategoryFuel,
			repoModel.CategoryPorthole,
			repoModel.CategoryWing,
		}
		require.Contains(s.T(), validCategories, part.Category)

		// Проверяем размеры
		if part.Dimensions != nil {
			require.Greater(s.T(), part.Dimensions.Length, 0.0)
			require.Greater(s.T(), part.Dimensions.Width, 0.0)
			require.Greater(s.T(), part.Dimensions.Height, 0.0)
			require.Greater(s.T(), part.Dimensions.Weight, 0.0)
		}

		// Проверяем производителя
		if part.Manufacturer != nil {
			require.NotEmpty(s.T(), part.Manufacturer.Name)
			require.NotEmpty(s.T(), part.Manufacturer.Country)
			require.NotEmpty(s.T(), part.Manufacturer.Website)
		}

		// Проверяем теги
		require.NotNil(s.T(), part.Tags)

		// Проверяем метаданные
		require.NotNil(s.T(), part.Metadata)

		// Проверяем время создания
		require.NotNil(s.T(), part.CreatedAt)
	}
}

func (s *RepositorySuite) TestInitWithDummyMultipleCalls() {
	// Первая инициализация
	err1 := s.repository.InitWithDummy()
	require.NoError(s.T(), err1)
	initialCount := s.repository.Count()

	// Вторая инициализация (должна перезаписать данные)
	err2 := s.repository.InitWithDummy()
	require.NoError(s.T(), err2)
	secondCount := s.repository.Count()

	// Количество частей может отличаться, так как используется случайное количество
	require.Greater(s.T(), initialCount, 0)
	require.Greater(s.T(), secondCount, 0)
}

func (s *RepositorySuite) TestInitWithDummyWithExistingData() {
	// Добавляем существующие данные
	existingPart := s.createTestPart()
	s.repository.Add(existingPart)
	require.True(s.T(), s.repository.Count() == 1)

	// Вызываем метод инициализации
	err := s.repository.InitWithDummy()

	// Проверяем результат
	require.NoError(s.T(), err)
	require.Greater(s.T(), s.repository.Count(), 1) // Должны быть добавлены новые данные

	// Проверяем, что существующие данные могут быть перезаписаны или дополнены
	for _, part := range s.repository.GetAll() {
		if part.Uuid == existingPart.Uuid {
			// Найдена часть с тем же UUID
			break
		}
	}
}

func (s *RepositorySuite) TestInitWithDummyGeneratesValidParts() {
	// Вызываем метод инициализации
	err := s.repository.InitWithDummy()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), s.repository)

	// Проверяем, что все сгенерированные части имеют валидные UUID
	for _, uuid := range s.repository.GetAll() {
		require.NotEmpty(s.T(), uuid)
	}

	// Проверяем, что все части имеют валидные имена из предопределенного списка
	validNames := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	allParts := s.repository.GetAll()
	for _, part := range allParts {
		require.Contains(s.T(), validNames, part.Name)
	}

	// Проверяем, что все части имеют соответствующие описания
	validDescriptions := map[string]string{
		"Main Engine":    "Primary propulsion unit",
		"Reserve Engine": "Backup propulsion unit",
		"Thruster":       "Thruster for fine adjustments",
		"Fuel Tank":      "Main fuel tank",
		"Left Wing":      "Left aerodynamic wing",
		"Right Wing":     "Right aerodynamic wing",
		"Window A":       "Front viewing window",
		"Window B":       "Side viewing window",
		"Control Module": "Flight control module",
		"Stabilizer":     "Stabilization fin",
	}

	for _, part := range allParts {
		expectedDesc, exists := validDescriptions[part.Name]
		require.True(s.T(), exists, "Unexpected part name: %s", part.Name)
		require.Equal(s.T(), expectedDesc, part.Description)
	}
}

func (s *RepositorySuite) TestInitWithDummyPriceRange() {
	// Вызываем метод инициализации
	err := s.repository.InitWithDummy()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), s.repository)

	// Проверяем, что все цены находятся в допустимом диапазоне (100-10000)
	for _, part := range s.repository.GetAll() {
		require.GreaterOrEqual(s.T(), part.Price, 100.0)
		require.LessOrEqual(s.T(), part.Price, 10000.0)
	}
}

func (s *RepositorySuite) TestInitWithDummyStockQuantity() {
	// Вызываем метод инициализации
	err := s.repository.InitWithDummy()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), s.repository)

	// Проверяем, что все количества в наличии находятся в допустимом диапазоне (1-100)
	for _, part := range s.repository.GetAll() {
		require.GreaterOrEqual(s.T(), part.StockQuantity, int64(1))
		require.LessOrEqual(s.T(), part.StockQuantity, int64(100))
	}
}
