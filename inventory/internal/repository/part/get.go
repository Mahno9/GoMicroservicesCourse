package part

import (
	"context"
	"errors"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/converter"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetPart(ctx context.Context, partUuid string) (*model.Part, error) {
	var part repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": partUuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}

	return converter.RepositoryToModelPart(&part), nil
}
