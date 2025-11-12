package part

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

func (r *repository) InitWithDummy(ctx context.Context) error {
	parts := generateParts()

	// for _, part := range parts {
	// 	r.Add(part)
	// }

	interfaceParts := make([]interface{}, len(parts))
	for i, part := range parts {
		interfaceParts[i] = part
	}

	result, err := r.collection.InsertMany(ctx, interfaceParts)
	if err != nil {
		return err
	}

	log.Printf("📝 Inserted %d parts\n", len(result.InsertedIDs))

	return nil
}

func generateParts() []*repoModel.Part {
	names := []string{
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

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	categories := []repoModel.Category{
		repoModel.CategoryEngine,
		repoModel.CategoryFuel,
		repoModel.CategoryPorthole,
		repoModel.CategoryWing,
	}

	var parts []*repoModel.Part
	for i := 0; i < gofakeit.Number(20, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &repoModel.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      categories[gofakeit.Number(0, len(categories)-1)],
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     time.Now(),
		})
	}

	return parts
}

func generateDimensions() *repoModel.Dimensions {
	return &repoModel.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *repoModel.Manufacturer {
	return &repoModel.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]any {
	metadata := make(map[string]any)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		value := generateMetadataValue()
		metadata[gofakeit.Word()] = &value
	}

	return metadata
}

func generateMetadataValue() any {
	switch gofakeit.Number(0, 3) {
	case 0:
		return gofakeit.Word()

	case 1:
		return int64(gofakeit.Number(1, 100))

	case 2:
		return roundTo(gofakeit.Float64Range(1, 100))

	case 3:
		return gofakeit.Bool()

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
