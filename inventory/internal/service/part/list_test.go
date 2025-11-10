package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsMainFlow() {
	var (
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()
		partUuid3 = gofakeit.UUID()

		partName1 = gofakeit.Name()
		partName2 = gofakeit.Name()
		partName3 = gofakeit.Name()

		partPrice1 = gofakeit.Float64()
		partPrice2 = gofakeit.Float64()
		partPrice3 = gofakeit.Float64()

		stockQty1 = int64(gofakeit.Number(1, 100))
		stockQty2 = int64(gofakeit.Number(1, 100))
		stockQty3 = int64(gofakeit.Number(1, 100))

		category1 = model.CategoryEngine
		category2 = model.CategoryFuel
		category3 = model.CategoryWing

		country1 = gofakeit.Country()
		country2 = gofakeit.Country()

		tags1 = []string{gofakeit.Word(), gofakeit.Word()}
		tags2 = []string{gofakeit.Word(), gofakeit.Word()}
		tags3 = []string{gofakeit.Word(), gofakeit.Word()}

		filter = &model.PartsFilter{
			Uuids: []string{
				partUuid1,
				partUuid2,
				partUuid3,
			},
			Names: []string{
				partName1,
				partName2,
			},
			Categories: []model.Category{
				category1,
				category2,
			},
			ManufacturerCountries: []string{
				country1,
				country2,
			},
			Tags: []string{
				tags1[0],
				tags2[1],
			},
		}

		expectedParts = []*model.Part{
			{
				Uuid:          partUuid1,
				Name:          partName1,
				Description:   gofakeit.Sentence(5),
				Price:         partPrice1,
				StockQuantity: stockQty1,
				Category:      category1,
				Tags:          tags1,
			},
			{
				Uuid:          partUuid2,
				Name:          partName2,
				Description:   gofakeit.Sentence(5),
				Price:         partPrice2,
				StockQuantity: stockQty2,
				Category:      category2,
				Tags:          tags2,
			},
			{
				Uuid:          partUuid3,
				Name:          partName3,
				Description:   gofakeit.Sentence(5),
				Price:         partPrice3,
				StockQuantity: stockQty3,
				Category:      category3,
				Tags:          tags3,
			},
		}
	)

	s.repository.On("ListParts", s.ctx, mock.MatchedBy(func(f *model.PartsFilter) bool {
		return len(f.Uuids) == 3 &&
			len(f.Names) == 2 &&
			len(f.Categories) == 2 &&
			len(f.ManufacturerCountries) == 2 &&
			len(f.Tags) == 2
	})).Return(expectedParts, nil)

	result, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)
	s.Equal(expectedParts, result)
}

func (s *ServiceSuite) TestListPartsRepositoryError() {
	var (
		filter = &model.PartsFilter{
			Uuids: []string{
				gofakeit.UUID(),
			},
		}

		expectedError = model.ErrPartNotFound
	)

	s.repository.On("ListParts", s.ctx, mock.AnythingOfType("*model.PartsFilter")).Return(nil, expectedError)

	result, err := s.service.ListParts(s.ctx, filter)

	s.ErrorIs(err, expectedError)
	s.Nil(result)
}

func (s *ServiceSuite) TestListPartsEmptyFilter() {
	var (
		filter = &model.PartsFilter{}

		expectedParts = []*model.Part{}
	)

	s.repository.On("ListParts", s.ctx, mock.MatchedBy(func(f *model.PartsFilter) bool {
		return len(f.Uuids) == 0 &&
			len(f.Names) == 0 &&
			len(f.Categories) == 0 &&
			len(f.ManufacturerCountries) == 0 &&
			len(f.Tags) == 0
	})).Return(expectedParts, nil)

	result, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)
	s.Equal(expectedParts, result)
}

func (s *ServiceSuite) TestListPartsWithOnlyUuidsFilter() {
	var (
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()

		filter = &model.PartsFilter{
			Uuids: []string{
				partUuid1,
				partUuid2,
			},
		}

		expectedParts = []*model.Part{
			{
				Uuid: partUuid1,
				Name: gofakeit.Name(),
			},
			{
				Uuid: partUuid2,
				Name: gofakeit.Name(),
			},
		}
	)

	s.repository.On("ListParts", s.ctx, mock.MatchedBy(func(f *model.PartsFilter) bool {
		return len(f.Uuids) == 2 && len(f.Names) == 0 && len(f.Categories) == 0
	})).Return(expectedParts, nil)

	result, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)
	s.Equal(expectedParts, result)
}
