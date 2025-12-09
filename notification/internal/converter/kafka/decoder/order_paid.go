package decoder

import (
	"github.com/Mahno9/GoMicroservicesCourse/notification/model"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type OrderPaidDecoder struct{}

func NewOrderPaidDecoder() *OrderPaidDecoder {
	return &OrderPaidDecoder{}
}

func (d *OrderPaidDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var protoEvent eventsV1.OrderPaid
	err := proto.Unmarshal(data, &protoEvent)
	if err != nil {
		return model.OrderPaidEvent{}, err
	}

	return model.OrderPaidEvent{
		OrderID:         protoEvent.GetOrderUuid(),
		UserID:          protoEvent.GetUserUuid(),
		PaymentMethod:   protoEvent.GetPaymentMethod(),
		TransactionUUID: protoEvent.GetTransactionUuid(),
	}, nil
}
