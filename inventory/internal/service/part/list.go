package part

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	timedCtx, cancel := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancel()

	return s.repository.ListParts(timedCtx, filter)
}
