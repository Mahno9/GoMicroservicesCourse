package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	timedCtx, cancel := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancel()

	logger.Info(ctx, "🟡 GetPart:", zap.String("uuid", uuid))

	return s.repository.GetPart(timedCtx, uuid)
}
