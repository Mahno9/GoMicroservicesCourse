package order

import (
	"log"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Update(order *model.Order) error {
	log.Printf("🔃 Updating order\n%+v\nvvv\n%+v\n", r.orders[order.OrderUuid], order)
	r.orders[order.OrderUuid] = converter.ModelToRepositoryOrder(order)
	return nil
}
