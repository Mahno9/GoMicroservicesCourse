package v1

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
)

type apiHandler struct {
	orderService service.OrderService
}

func NewAPIHandler(service service.OrderService) *apiHandler {
	return &apiHandler{
		orderService: service,
	}
}
