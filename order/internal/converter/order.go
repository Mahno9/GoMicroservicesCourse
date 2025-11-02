package converter

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func ApiToModelOrderInfo(apiOrder *orderV1.CreateOrderReq) model.CreateOrderData {
	newOrderData := model.CreateOrderData{
		UserUuid: string(apiOrder.UserUUID),
	}

	newOrderData.PartUuids = make(map[string]any)
	for _, partUuid := range apiOrder.PartUuids {
		newOrderData.PartUuids[string(partUuid)] = nil
	}

	return newOrderData
}

func ConvertPaymentMethod(orderPaymentMethod *orderV1.PaymentMethod) (int32, error) {
	paymentMethodBytes, err := orderPaymentMethod.MarshalText()
	if err != nil {
		return 0, model.UnknownPaymentMethodErr
	}

	// Convert string representation to payment method enum value
	// This is a simplified conversion - in a real implementation you might need
	// to map between different enum types more carefully
	paymentMethodStr := string(paymentMethodBytes)
	
	// Default to unknown payment method
	if paymentMethodStr == "" {
		return 0, model.UnknownPaymentMethodErr
	}

	// Simple conversion - in real implementation you'd need proper mapping
	// between orderV1.PaymentMethod and int32 representation
	switch orderV1.PaymentMethod(paymentMethodStr) {
	case orderV1.PaymentMethodCREDIT_CARD:
		return 1, nil
	case orderV1.PaymentMethodCASH:
		return 2, nil
	case orderV1.PaymentMethodBANK_TRANSFER:
		return 3, nil
	default:
		return 0, model.UnknownPaymentMethodErr
	}
}
