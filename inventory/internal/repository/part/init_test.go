package part

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"

	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

func (s *RepositorySuite) TestInitWithDummyMainFlow() {
	// Настраиваем мок для InsertMany метода
	s.collection.On("InsertMany", s.ctx, mock.AnythingOfType("[]interface {}")).Return(&mongo.InsertManyResult{
		InsertedIDs: []interface{}{"id1", "id2", "id3"},
	}, nil)

	// Вызываем метод инициализации
	err := s.repository.InitWithDummy(s.ctx)

	// Проверяем результат
	require.NoError(s.T(), err)
}

func (s *RepositorySuite) TestInitWithDummyGeneratesValidParts() {
	// Настраиваем мок для InsertMany метода с проверкой типа данных
	s.collection.On("InsertMany", s.ctx, mock.MatchedBy(func(docs []interface{}) bool {
		// Проверяем что все документы являются валидными частями
		for _, doc := range docs {
			part, ok := doc.(*repoModel.Part)
			if !ok {
				return false
			}

			// Проверяем валидные имена
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

			nameFound := false
			for _, name := range validNames {
				if part.Name == name {
					nameFound = true
					break
				}
			}

			if !nameFound {
				return false
			}

			// Проверяем валидные UUID
			if part.Uuid == "" {
				return false
			}
		}
		return true
	})).Return(&mongo.InsertManyResult{
		InsertedIDs: []interface{}{"id1", "id2", "id3"},
	}, nil)

	// Вызываем метод инициализации
	err := s.repository.InitWithDummy(s.ctx)
	require.NoError(s.T(), err)
}

func (s *RepositorySuite) TestInitWithDummyPriceRange() {
	// Настраиваем мок для InsertMany метода с проверкой цен
	s.collection.On("InsertMany", s.ctx, mock.MatchedBy(func(docs []interface{}) bool {
		// Проверяем что все цены находятся в допустимом диапазоне (100-10000)
		for _, doc := range docs {
			part, ok := doc.(*repoModel.Part)
			if !ok {
				return false
			}

			if part.Price < 100.0 || part.Price > 10000.0 {
				return false
			}
		}
		return true
	})).Return(&mongo.InsertManyResult{
		InsertedIDs: []interface{}{"id1", "id2", "id3"},
	}, nil)

	// Вызываем метод инициализации
	err := s.repository.InitWithDummy(s.ctx)
	require.NoError(s.T(), err)
}

func (s *RepositorySuite) TestInitWithDummyStockQuantity() {
	// Настраиваем мок для InsertMany метода с проверкой количества на складе
	s.collection.On("InsertMany", s.ctx, mock.MatchedBy(func(docs []interface{}) bool {
		// Проверяем что все количества в наличии находятся в допустимом диапазоне (1-100)
		for _, doc := range docs {
			part, ok := doc.(*repoModel.Part)
			if !ok {
				return false
			}

			if part.StockQuantity < 1 || part.StockQuantity > 100 {
				return false
			}
		}
		return true
	})).Return(&mongo.InsertManyResult{
		InsertedIDs: []interface{}{"id1", "id2", "id3"},
	}, nil)

	// Вызываем метод инициализации
	err := s.repository.InitWithDummy(s.ctx)
	require.NoError(s.T(), err)
}
