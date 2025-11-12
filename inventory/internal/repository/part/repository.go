package part

import (
	"context"
	"time"

	model "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	def "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databasePartsCollection = "parts"
	indexCreationTimeout    = 5 * time.Second
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection def.MongoCollection
}

func NewRepository(ctx context.Context, db def.MongoDatabase) (*repository, error) {
	collection := db.Collection(databasePartsCollection)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "uuid", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(ctx, indexCreationTimeout)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, model.ErrDbIndexInitFailed
	}

	return &repository{
		collection: collection,
	}, nil
}
