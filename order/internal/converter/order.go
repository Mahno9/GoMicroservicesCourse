package converter

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func ApiToModelOrderInfo(apiOrder *orderV1.CreateOrderReq) model.CreateOrderData {
	if apiOrder == nil {
		return model.CreateOrderData{}
	}

	partUuids := make([]string, len(apiOrder.PartUuids))
	copy(partUuids, apiOrder.PartUuids)

	return model.CreateOrderData{
		UserUuid:  string(apiOrder.UserUUID),
		PartUuids: partUuids,
	}
}

func ModelToApiGetOrder(modelOrder *model.Order) *orderV1.GetOrderOK {
	return &orderV1.GetOrderOK{
		OrderUUID:       orderV1.OrderUUID(modelOrder.OrderUuid),
		UserUUID:        orderV1.UserUUID(modelOrder.UserUuid),
		PartUuids:       modelOrder.PartUuids,
		TotalPrice:      orderV1.TotalPrice(modelOrder.TotalPrice),
		TransactionUUID: orderV1.TransactionUUID(modelOrder.TransactionUuid),
		PaymentMethod:   ModelToApiPaymentMethod(modelOrder.PaymentMethod),
		Status:          ModelToApiOrderStatus(int32(modelOrder.Status)),
	}
}

func ModelToApiOrderStatus(modelStatus int32) orderV1.Status {
	switch modelStatus {
	case 2:
		return orderV1.StatusPAID
	case 3:
		return orderV1.StatusCANCELLED
	default:
		return orderV1.StatusPENDINGPAYMENT
	}
}

func ModelToApiPaymentMethod(modelPaymentMethod int32) orderV1.PaymentMethod {
	switch modelPaymentMethod {
	case 1:
		return orderV1.PaymentMethodCARD
	case 2:
		return orderV1.PaymentMethodSBP
	case 3:
		return orderV1.PaymentMethodCREDITCARD
	case 4:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func ApiToModelPaymentMethod(apiPaymentMethod orderV1.PaymentMethod) int32 {
	switch apiPaymentMethod {
	case orderV1.PaymentMethodCARD:
		return 1
	case orderV1.PaymentMethodSBP:
		return 2
	case orderV1.PaymentMethodCREDITCARD:
		return 3
	case orderV1.PaymentMethodINVESTORMONEY:
		return 4
	default:
		return 0
	}
}
