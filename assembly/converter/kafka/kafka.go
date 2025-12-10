package kafka

import "github.com/Mahno9/GoMicroservicesCourse/assembly/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}
