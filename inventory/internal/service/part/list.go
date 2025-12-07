package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	timedCtx, cancel := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancel()

	logger.Info(ctx, "🟡 ListParts:", zap.Any("filter", filter))

	return s.repository.ListParts(timedCtx, filter)
}
