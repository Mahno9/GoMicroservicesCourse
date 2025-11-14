package converter

import (
	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/model"
)

// ModelToRepositoryOrder конвертирует модель бизнес-логики в модель репозитория
func ModelToRepositoryOrder(order *model.Order) *repoModel.Order {
	if order == nil {
		return nil
	}

	var transactionUuidStr *string
	if order.TransactionUuid != nil {
		str := order.TransactionUuid.String()
		transactionUuidStr = &str
	}

	partUuidsStr := make([]string, len(order.PartUuids))
	for i, partUuid := range order.PartUuids {
		partUuidsStr[i] = partUuid.String()
	}

	return &repoModel.Order{
		OrderUuid:       order.OrderUuid.String(),
		UserUuid:        order.UserUuid.String(),
		PartUuids:       partUuidsStr,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: transactionUuidStr,
		PaymentMethod:   modelToRepoPaymentMethod(order.PaymentMethod),
		Status:          modelToRepoStatus(order.Status),
	}
}

// RepositoryToModelOrder конвертирует модель репозитория в модель бизнес-логики
func RepositoryToModelOrder(order *repoModel.Order) *model.Order {
	if order == nil {
		return nil
	}

	orderUuid, err := uuid.Parse(order.OrderUuid)
	if err != nil {
		return nil
	}
	userUuid, err := uuid.Parse(order.UserUuid)
	if err != nil {
		return nil
	}

	var transactionUuid *uuid.UUID
	if order.TransactionUuid != nil {
		parsed, err := uuid.Parse(*order.TransactionUuid)
		if err == nil {
			transactionUuid = &parsed
		}
	}

	partUuids := make([]uuid.UUID, len(order.PartUuids))
	for i, partUuidStr := range order.PartUuids {
		parsed, err := uuid.Parse(partUuidStr)
		if err != nil {
			continue
		}
		partUuids[i] = parsed
	}

	return &model.Order{
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		PartUuids:       partUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: transactionUuid,
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
