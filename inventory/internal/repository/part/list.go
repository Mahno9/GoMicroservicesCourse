package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/converter"
	remoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
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
			logger.Error(ctx, "❗ failed to close cursor: %v", zap.Error(cerr))
		}
	}()

	repoParts := make([]*remoModel.Part, 0)
	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, err
	}

	return converter.RepositoryToModelParts(repoParts), nil
}
