package part

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/converter"
)

func (r *repository) ListParts(ctx context.Context, filters *model.PartsFilter) ([]*model.Part, error) {
	parts := make([]*model.Part, 0)
	repoFilter := converter.ModelToRepositoryFilter(filters)

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

	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
