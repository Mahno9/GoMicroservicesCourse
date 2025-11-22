package part

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/converter"
	remoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filters *model.PartsFilter) ([]*model.Part, error) {
	repoFilter := converter.ModelToMongoFilter(filters)

	cursor, err := r.collection.Find(ctx, repoFilter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("❗ failed to close cursor: %v\n", cerr)
		}
	}()

	repoParts := make([]*remoModel.Part, 0)
	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, err
	}

	return converter.RepositoryToModelParts(repoParts), nil
}
