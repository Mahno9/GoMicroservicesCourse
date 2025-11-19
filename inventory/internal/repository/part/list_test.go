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

	// Создаем мок для курсора
	cursor := mocks.NewMongoCursor(s.T())
	cursor.On("All", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		// Получаем указатель на слайс и заполняем его данными
		parts := args.Get(1).(*[]*repoModel.Part)
		*parts = []*repoModel.Part{
			{
				Uuid:          part1.Uuid,
				Name:          part1.Name,
				Description:   part1.Description,
				Price:         part1.Price,
				StockQuantity: part1.StockQuantity,
				Category:      repoModel.Category(part1.Category),
				Tags:          part1.Tags,
				Metadata:      part1.Metadata,
				CreatedAt:     part1.CreatedAt,
				UpdatedAt:     part1.UpdatedAt,
				Dimensions: &repoModel.Dimensions{
					Length: part1.Dimensions.Length,
					Width:  part1.Dimensions.Width,
					Height: part1.Dimensions.Height,
					Weight: part1.Dimensions.Weight,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    part1.Manufacturer.Name,
					Country: part1.Manufacturer.Country,
					Website: part1.Manufacturer.Website,
				},
			},
			{
				Uuid:          part2.Uuid,
				Name:          part2.Name,
				Description:   part2.Description,
				Price:         part2.Price,
				StockQuantity: part2.StockQuantity,
				Category:      repoModel.Category(part2.Category),
				Tags:          part2.Tags,
				Metadata:      part2.Metadata,
				CreatedAt:     part2.CreatedAt,
				UpdatedAt:     part2.UpdatedAt,
				Dimensions: &repoModel.Dimensions{
					Length: part2.Dimensions.Length,
					Width:  part2.Dimensions.Width,
					Height: part2.Dimensions.Height,
					Weight: part2.Dimensions.Weight,
				},
				Manufacturer: &repoModel.Manufacturer{
					Name:    part2.Manufacturer.Name,
					Country: part2.Manufacturer.Country,
					Website: part2.Manufacturer.Website,
				},
			},
		}
	}).Return(nil)

	cursor.On("Close", s.ctx).Return(nil)

	// Настраиваем мок для Find метода
	s.collection.On("Find", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		uuids, ok := v["uuid"].(bson.M)
		if !ok {
			return false
		}
		uuidIn, ok := uuids["$in"].([]string)
		if !ok {
			return false
		}
		// Проверяем что в фильтре есть part1 и part2 UUID
		return contains(uuidIn, part1.Uuid) && contains(uuidIn, part2.Uuid)
	})).Return(cursor, nil)

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
	require.Len(s.T(), resultUuids, 2)
}

func (s *RepositorySuite) TestListPartsEmptyFilter() {
	// Создаем тестовые части
	part1 := s.createTestPart()
	part2 := s.createTestPart()

	// Создаем мок для курсора
	cursor := mocks.NewMongoCursor(s.T())
	cursor.On("All", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		parts := args.Get(1).(*[]*repoModel.Part)
		*parts = []*repoModel.Part{
			{
				Uuid:          part1.Uuid,
				Name:          part1.Name,
				Description:   part1.Description,
				Price:         part1.Price,
				StockQuantity: part1.StockQuantity,
				Category:      repoModel.Category(part1.Category),
				Tags:          part1.Tags,
				Metadata:      part1.Metadata,
				CreatedAt:     part1.CreatedAt,
				UpdatedAt:     part1.UpdatedAt,
			},
			{
				Uuid:          part2.Uuid,
				Name:          part2.Name,
				Description:   part2.Description,
				Price:         part2.Price,
				StockQuantity: part2.StockQuantity,
				Category:      repoModel.Category(part2.Category),
				Tags:          part2.Tags,
				Metadata:      part2.Metadata,
				CreatedAt:     part2.CreatedAt,
				UpdatedAt:     part2.UpdatedAt,
			},
		}
	}).Return(nil)

	cursor.On("Close", s.ctx).Return(nil)

	// Настраиваем мок для Find метода с пустым фильтром
	s.collection.On("Find", s.ctx, bson.M{}).Return(cursor, nil)

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
	part3 := s.createTestPart()

	// Создаем мок для курсора
	cursor := mocks.NewMongoCursor(s.T())
	cursor.On("All", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		parts := args.Get(1).(*[]*repoModel.Part)
		*parts = []*repoModel.Part{
			{
				Uuid:          part1.Uuid,
				Name:          part1.Name,
				Description:   part1.Description,
				Price:         part1.Price,
				StockQuantity: part1.StockQuantity,
				Category:      repoModel.Category(part1.Category),
				Tags:          part1.Tags,
				Metadata:      part1.Metadata,
				CreatedAt:     part1.CreatedAt,
				UpdatedAt:     part1.UpdatedAt,
			},
			{
				Uuid:          part3.Uuid,
				Name:          part3.Name,
				Description:   part3.Description,
				Price:         part3.Price,
				StockQuantity: part3.StockQuantity,
				Category:      repoModel.Category(part3.Category),
				Tags:          part3.Tags,
				Metadata:      part3.Metadata,
				CreatedAt:     part3.CreatedAt,
				UpdatedAt:     part3.UpdatedAt,
			},
		}
	}).Return(nil)

	cursor.On("Close", s.ctx).Return(nil)

	// Настраиваем мок для Find метода
	s.collection.On("Find", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		uuids, ok := v["uuid"].(bson.M)
		if !ok {
			return false
		}
		uuidIn, ok := uuids["$in"].([]string)
		if !ok {
			return false
		}
		// Проверяем что в фильтре есть part1 и part3 UUID
		return contains(uuidIn, part1.Uuid) && contains(uuidIn, part3.Uuid)
	})).Return(cursor, nil)

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
	require.True(s.T(), resultUuids[part3.Uuid])
	require.Len(s.T(), result, 2)
}

func (s *RepositorySuite) TestListPartsOnlyTagsFilter() {
	// Создаем тестовые части с разными тегами
	part1 := s.createTestPart()
	part1.Tags = []string{"engine", "primary"}

	part3 := s.createTestPart()
	part3.Tags = []string{"wing", "aerodynamic"}

	// Создаем мок для курсора
	cursor := mocks.NewMongoCursor(s.T())
	cursor.On("All", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		parts := args.Get(1).(*[]*repoModel.Part)
		*parts = []*repoModel.Part{
			{
				Uuid:          part1.Uuid,
				Name:          part1.Name,
				Description:   part1.Description,
				Price:         part1.Price,
				StockQuantity: part1.StockQuantity,
				Category:      repoModel.Category(part1.Category),
				Tags:          part1.Tags,
				Metadata:      part1.Metadata,
				CreatedAt:     part1.CreatedAt,
				UpdatedAt:     part1.UpdatedAt,
			},
			{
				Uuid:          part3.Uuid,
				Name:          part3.Name,
				Description:   part3.Description,
				Price:         part3.Price,
				StockQuantity: part3.StockQuantity,
				Category:      repoModel.Category(part3.Category),
				Tags:          part3.Tags,
				Metadata:      part3.Metadata,
				CreatedAt:     part3.CreatedAt,
				UpdatedAt:     part3.UpdatedAt,
			},
		}
	}).Return(nil)

	cursor.On("Close", s.ctx).Return(nil)

	// Настраиваем мок для Find метода
	s.collection.On("Find", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		tags, ok := v["tags"].(bson.M)
		if !ok {
			return false
		}
		tagsAll, ok := tags["$all"].([]string)
		if !ok {
			return false
		}
		// Проверяем что в фильтре есть теги "primary" и "aerodynamic"
		return contains(tagsAll, "primary") && contains(tagsAll, "aerodynamic")
	})).Return(cursor, nil)

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
	require.True(s.T(), resultUuids[part3.Uuid])
	require.Len(s.T(), result, 2)
}

func (s *RepositorySuite) TestListPartsNoMatches() {
	// Создаем тестовую часть
	part1 := s.createTestPart()

	// Настраиваем мок для Find метода, который возвращает ошибку "нет документов"
	s.collection.On("Find", s.ctx, mock.MatchedBy(func(v bson.M) bool {
		uuids, ok := v["uuid"].(bson.M)
		if !ok {
			return false
		}
		uuidIn, ok := uuids["$in"].([]string)
		if !ok {
			return false
		}
		// Проверяем что в фильтре есть случайный UUID, которого нет в базе
		return len(uuidIn) == 1 && uuidIn[0] != part1.Uuid
	})).Return(nil, mongo.ErrNoDocuments)

	// Создаем фильтр, который не совпадает ни с одной частью
	filter := &model.PartsFilter{
		Uuids: []string{
			gofakeit.UUID(), // Случайный UUID, которого нет в репозитории
		},
	}

	// Вызываем метод
	result, err := s.repository.ListParts(s.ctx, filter)

	// Проверяем результат
	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, model.ErrPartNotFound)
	require.Empty(s.T(), result) // Не должно быть совпадений
}

func (s *RepositorySuite) TestListPartsEmptyRepository() {
	// Настраиваем мок для Find метода, который возвращает пустой курсор
	cursor := mocks.NewMongoCursor(s.T())
	cursor.On("All", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		parts := args.Get(1).(*[]*repoModel.Part)
		*parts = []*repoModel.Part{} // Пустой слайс
	}).Return(nil)

	cursor.On("Close", s.ctx).Return(nil)

	// Настраиваем мок для Find метода с пустым фильтром
	s.collection.On("Find", s.ctx, bson.M{}).Return(cursor, nil)

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

	// Создаем мок для курсора
	cursor := mocks.NewMongoCursor(s.T())
	cursor.On("All", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		parts := args.Get(1).(*[]*repoModel.Part)
		*parts = []*repoModel.Part{
			{
				Uuid:          part1.Uuid,
				Name:          part1.Name,
				Description:   part1.Description,
				Price:         part1.Price,
				StockQuantity: part1.StockQuantity,
				Category:      repoModel.Category(part1.Category),
				Tags:          part1.Tags,
				Metadata:      part1.Metadata,
				CreatedAt:     part1.CreatedAt,
				UpdatedAt:     part1.UpdatedAt,
			},
			{
				Uuid:          part2.Uuid,
				Name:          part2.Name,
				Description:   part2.Description,
				Price:         part2.Price,
				StockQuantity: part2.StockQuantity,
				Category:      repoModel.Category(part2.Category),
				Tags:          part2.Tags,
				Metadata:      part2.Metadata,
				CreatedAt:     part2.CreatedAt,
				UpdatedAt:     part2.UpdatedAt,
			},
		}
	}).Return(nil)

	cursor.On("Close", s.ctx).Return(nil)

	// Настраиваем мок для Find метода - можно вызывать многократно
	s.collection.On("Find", s.ctx, bson.M{}).Return(cursor, nil)

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
	for range 10 {
		<-done
	}
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
