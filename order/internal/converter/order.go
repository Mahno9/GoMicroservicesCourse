package converter

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func ApiToModelOrderInfo(apiOrder *genOrderV1.CreateOrderReq) model.CreateOrderData {
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

func ModelToApiGetOrder(modelOrder *model.Order) *genOrderV1.GetOrderOK {
	return &genOrderV1.GetOrderOK{
		OrderUUID:       genOrderV1.OrderUUID(modelOrder.OrderUuid),
		UserUUID:        genOrderV1.UserUUID(modelOrder.UserUuid),
		PartUuids:       modelOrder.PartUuids,
		TotalPrice:      genOrderV1.TotalPrice(modelOrder.TotalPrice),
		TransactionUUID: genOrderV1.TransactionUUID(modelOrder.TransactionUuid),
		PaymentMethod:   ModelToApiPaymentMethod(modelOrder.PaymentMethod),
		Status:          ModelToApiOrderStatus(int32(modelOrder.Status)),
	}
}

func ModelToApiOrderStatus(modelStatus int32) genOrderV1.Status {
	switch modelStatus {
	case 1:
		return genOrderV1.StatusPAID
	case 2:
		return genOrderV1.StatusCANCELLED
	default:
		return genOrderV1.StatusPENDINGPAYMENT
	}
}

func ModelToApiPaymentMethod(modelPaymentMethod int32) genOrderV1.PaymentMethod {
	switch modelPaymentMethod {
	case 1:
		return genOrderV1.PaymentMethodCARD
	case 2:
		return genOrderV1.PaymentMethodSBP
	case 3:
		return genOrderV1.PaymentMethodCREDITCARD
	case 4:
		return genOrderV1.PaymentMethodINVESTORMONEY
	default:
		return genOrderV1.PaymentMethodUNKNOWN
	}
}

func ApiToModelPaymentMethod(apiPaymentMethod genOrderV1.PaymentMethod) int32 {
	switch apiPaymentMethod {
	case genOrderV1.PaymentMethodCARD:
		return 1
	case genOrderV1.PaymentMethodSBP:
		return 2
	case genOrderV1.PaymentMethodCREDITCARD:
		return 3
	case genOrderV1.PaymentMethodINVESTORMONEY:
		return 4
	default:
		return 0
	}
}
