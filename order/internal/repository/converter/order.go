package converter

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/model"
)

// ModelToRepositoryOrder конвертирует модель бизнес-логики в модель репозитория
func ModelToRepositoryOrder(order *model.Order) *repoModel.Order {
	return &repoModel.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}

// RepositoryOrderToModel конвертирует модель репозитория в модель бизнес-логики
func RepositoryOrderToModel(order *repoModel.Order) *model.Order {
	return &model.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}
