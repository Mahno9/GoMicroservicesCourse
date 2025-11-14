package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=OrderRepository --output=mocks

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	Get(ctx context.Context, orderUuid uuid.UUID) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
