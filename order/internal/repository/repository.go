package repository

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	Get(ctx context.Context, orderUuid string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
