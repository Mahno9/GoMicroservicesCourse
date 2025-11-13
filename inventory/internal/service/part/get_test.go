package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartMainFlow() {
	var (
		partUuid  = gofakeit.UUID()
		partName  = gofakeit.Name()
		partDesc  = gofakeit.Sentence(5)
		partPrice = gofakeit.Float64Range(10.0, 1000.0)
		stockQty  = int64(gofakeit.Number(1, 100))
		category  = model.CategoryEngine
		country   = gofakeit.Country()
		website   = gofakeit.URL()
		tag1      = gofakeit.Word()
		tag2      = gofakeit.Word()

		now = time.Now()

		expectedPart = &model.Part{
			Uuid:          partUuid,
			Name:          partName,
			Description:   partDesc,
			Price:         partPrice,
			StockQuantity: stockQty,
			Category:      category,
			Dimensions: &model.Dimensions{
				Length: gofakeit.Float64Range(1.0, 100.0),
				Width:  gofakeit.Float64Range(1.0, 100.0),
				Height: gofakeit.Float64Range(1.0, 100.0),
				Weight: gofakeit.Float64Range(1.0, 100.0),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: country,
				Website: website,
			},
			Tags: []string{tag1, tag2},
			Metadata: func() map[string]any {
				value := any(new(string))
				return map[string]any{"key1": value}
			}(),
			CreatedAt: now,
			UpdatedAt: &now,
		}
	)

	s.repository.On("GetPart", mock.Anything, partUuid).Return(expectedPart, nil)

	result, err := s.service.GetPart(s.ctx, partUuid)

	s.NoError(err)
	s.Equal(expectedPart, result)
}

func (s *ServiceSuite) TestGetPartRepositoryError() {
	var (
		partUuid      = gofakeit.UUID()
		expectedError = model.ErrPartNotFound
	)

	s.repository.On("GetPart", mock.Anything, partUuid).Return(nil, expectedError)

	result, err := s.service.GetPart(s.ctx, partUuid)

	s.ErrorIs(err, expectedError)
	s.Nil(result)
}

func (s *ServiceSuite) TestGetPartEmptyUuid() {
	var (
		partUuid      = ""
		expectedError = model.ErrPartNotFound
	)

	s.repository.On("GetPart", mock.Anything, partUuid).Return(nil, expectedError)

	result, err := s.service.GetPart(s.ctx, partUuid)

	s.ErrorIs(err, expectedError)
	s.Nil(result)
}

func (s *ServiceSuite) TestGetPartWithMinimalData() {
	var (
		partUuid = gofakeit.UUID()
		partName = gofakeit.Name()

		expectedPart = &model.Part{
			Uuid:          partUuid,
			Name:          partName,
			Description:   "",
			Price:         0.0,
			StockQuantity: 0,
			Category:      model.CategoryUnknown,
			// Все остальные поля nil или пустые
		}
	)

	s.repository.On("GetPart", mock.Anything, partUuid).Return(expectedPart, nil)

	result, err := s.service.GetPart(s.ctx, partUuid)

	s.NoError(err)
	s.Equal(expectedPart, result)
	s.Equal(partUuid, result.Uuid)
	s.Equal(partName, result.Name)
	s.Equal("", result.Description)
	s.Equal(0.0, result.Price)
	s.Equal(int64(0), result.StockQuantity)
	s.Equal(model.CategoryUnknown, result.Category)
	s.Equal(time.Time{}, result.CreatedAt)
	s.Nil(result.Dimensions)
	s.Nil(result.Manufacturer)
	s.Nil(result.Tags)
	s.Nil(result.Metadata)
	s.Nil(result.UpdatedAt)
}
