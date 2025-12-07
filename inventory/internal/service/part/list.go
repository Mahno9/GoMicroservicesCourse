package part

import (
	"context"
	"log"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	timedCtx, cancel := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancel()

	log.Println("🟡 ListParts:", filter)

	return s.repository.ListParts(timedCtx, filter)
}
