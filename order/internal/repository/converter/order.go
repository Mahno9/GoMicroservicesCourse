package converter

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/model"
)

// ModelToRepositoryOrder конвертирует модель бизнес-логики в модель репозитория
func ModelToRepositoryOrder(order *model.Order) *repoModel.Order {
	if order == nil {
		return nil
	}

	return &repoModel.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   modelToRepoPaymentMethod(order.PaymentMethod),
		Status:          modelToRepoStatus(order.Status),
	}
}

// RepositoryToModelOrder конвертирует модель репозитория в модель бизнес-логики
func RepositoryToModelOrder(order *repoModel.Order) *model.Order {
	if order == nil {
		return nil
	}

	return &model.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   repoToModelPaymentMethod(order.PaymentMethod),
		Status:          repoToModelStatus(order.Status),
	}
}

func modelToRepoStatus(status model.OrderStatus) string {
	switch status {
	case model.StatusPAID:
		return repoModel.StatusPaid
	case model.StatusCANCELLED:
		return repoModel.StatusCancelled
	default:
		return repoModel.StatusPendingPayment
	}
}

func repoToModelStatus(status string) model.OrderStatus {
	switch status {
	case repoModel.StatusPaid:
		return model.StatusPAID
	case repoModel.StatusCancelled:
		return model.StatusCANCELLED
	default:
		return model.StatusPENDINGPAYMENT
	}
}

func modelToRepoPaymentMethod(method int32) string {
	switch method {
	case 1:
		return repoModel.PaymentMethodCard
	case 2:
		return repoModel.PaymentMethodSBP
	case 3:
		return repoModel.PaymentMethodCreditCard
	case 4:
		return repoModel.PaymentMethodInvestor
	default:
		return repoModel.PaymentMethodUnknown
	}
}

func repoToModelPaymentMethod(method string) int32 {
	switch method {
	case repoModel.PaymentMethodCard:
		return 1
	case repoModel.PaymentMethodSBP:
		return 2
	case repoModel.PaymentMethodCreditCard:
		return 3
	case repoModel.PaymentMethodInvestor:
		return 4
	default:
		return 0
	}
}
