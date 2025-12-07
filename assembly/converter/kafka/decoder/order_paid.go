package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/model"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
)

type OrderPaidDecoder struct{}

func (d *OrderPaidDecoder) Decode(data []byte) (model.OrderPaid, error) {
	var order eventsV1.OrderPaid
	err := proto.Unmarshal(data, &order)
	if err != nil {
		return model.OrderPaid{}, fmt.Errorf("filaed to unmarshal event: %w", err)
	}

	return model.OrderPaid{
		Uuid:            order.Uuid,
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PaymentMethod:   order.PaymentMethod,
		TransactionUuid: order.TransactionUuid,
	}, nil
}
