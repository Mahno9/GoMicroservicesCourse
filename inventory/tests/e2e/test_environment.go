package e2e

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"

	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

const (
	mongoDatabaseNameKey = "MONGO_INITDB_DATABASE"
	mongoCollectionName  = "parts"
)

func grpcToMongoCategory(grpcCategory inventoryV1.Category) string {
	switch grpcCategory {
	case inventoryV1.Category_ENGINE:
		return "ENGINE"
	case inventoryV1.Category_FUEL:
		return "FUEL"
	case inventoryV1.Category_PORTHOLE:
		return "PORTHOLE"
	case inventoryV1.Category_WING:
		return "WING"
	default:
		return "UNKNOWN"
	}
}

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           gofakeit.Name(),
		"description":    gofakeit.Sentence(5),
		"price":          gofakeit.Float64Range(10.0, 1000.0),
		"stock_quantity": int64(gofakeit.Number(1, 100)),
		"category":       grpcToMongoCategory(inventoryV1.Category_ENGINE),
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(1.0, 100.0),
			"width":  gofakeit.Float64Range(1.0, 100.0),
			"height": gofakeit.Float64Range(1.0, 100.0),
			"weight": gofakeit.Float64Range(1.0, 100.0),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags": []string{gofakeit.Word(), gofakeit.Word()},
		"metadata": bson.M{
			"key1": gofakeit.Word(),
		},
		"created_at": now,
		"updated_at": &now,
	}

	databaseName := getEnvWithLogging(ctx, mongoDatabaseNameKey)

	_, err := env.Mongo.Client().Database(databaseName).Collection(mongoCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, info *inventoryV1.Part) (string, error) {
	partUUID := gofakeit.UUID()

	createdAt := time.Now()
	createdAtTimestampbp := info.GetCreatedAt()
	if createdAtTimestampbp != nil {
		createdAt = createdAtTimestampbp.AsTime()
	}

	var updatedAt time.Time
	updatedAtTimestampbp := info.GetUpdatedAt()
	if updatedAtTimestampbp != nil {
		updatedAt = updatedAtTimestampbp.AsTime()
	}

	metadata := info.GetMetadata()
	metadataBson := bson.M{}
	if len(metadata) > 0 {
		for k, v := range metadata {
			metadataBson[k] = v
		}
	}

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           info.GetName(),
		"description":    info.GetDescription(),
		"price":          info.GetPrice(),
		"stock_quantity": info.GetStockQuantity(),
		"category":       grpcToMongoCategory(info.GetCategory()),
		"dimensions": bson.M{
			"length": info.GetDimensions().GetLength(),
			"width":  info.GetDimensions().GetWidth(),
			"height": info.GetDimensions().GetHeight(),
			"weight": info.GetDimensions().GetWeight(),
		},
		"manufacturer": bson.M{
			"name":    info.GetManufacturer().GetName(),
			"country": info.GetManufacturer().GetCountry(),
			"website": info.GetManufacturer().GetWebsite(),
		},
		"tags":       info.GetTags(),
		"metadata":   metadataBson,
		"created_at": createdAt,
		"updated_at": updatedAt,
	}

	_, err := env.Mongo.Client().Database(getEnvWithLogging(ctx, mongoDatabaseNameKey)).Collection(mongoCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

func (env *TestEnvironment) GetTestPartInfo() *inventoryV1.Part {
	metadata := map[string]*inventoryV1.Value{
		"key1": {
			Kind: &inventoryV1.Value_StringValue{
				StringValue: "value1",
			},
		},
		"key2": {
			Kind: &inventoryV1.Value_StringValue{
				StringValue: "value2",
			},
		},
	}

	return &inventoryV1.Part{
		Name:          "Boeing wing",
		Description:   "Most expensive detail in the world",
		Price:         1000,
		StockQuantity: 10,
		Category:      inventoryV1.Category_WING,
		Dimensions: &inventoryV1.Dimensions{
			Length: 10,
			Width:  10,
			Height: 10,
			Weight: 10,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Spirit AeroSystems",
			Country: "USA",
			Website: "https://www.spiritaero.com",
		},
		Tags:     []string{"wing", "detail"},
		Metadata: metadata,
	}
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := getEnvWithLogging(ctx, mongoDatabaseNameKey)
	if databaseName == "" {
		databaseName = "part-service" // fallback value
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(mongoCollectionName).DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}

	return nil
}
