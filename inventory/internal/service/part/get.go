package part

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	timedCtx, cancel := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancel()

	return s.repository.GetPart(timedCtx, uuid)
}
