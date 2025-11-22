package part

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *service) InitWithDummy(ctx context.Context) error {
	timedCtx, cancel := context.WithTimeout(ctx, model.RequestTimeoutUpdate)
	defer cancel()

	return s.repository.InitWithDummy(timedCtx)
}
