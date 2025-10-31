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

func convertPaymentMethod(orderPaymentMethod *orderV1.PaymentMethod) (int32, error) {
	paymentMethodBytes, err := orderPaymentMethod.MarshalText()
	if err != nil {
		return 0, model.unknownPaymentMethodErr

		log.Printf("❗ Unknown payment method: %v\n", orderPaymentMethod)
		paymentMethodBytes = []byte(paymentV1.PaymentMethod_UNKNOWN.String())
	}

	return paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(paymentMethodBytes)])
}
