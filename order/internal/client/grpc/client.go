package grpc

import (
	"context"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, in *model.PartsFilter) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, paymentData model.PayOrderData) (string, error)
}
