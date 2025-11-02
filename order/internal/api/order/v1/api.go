package v1

import (
	"time"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
)

const (
	createOrderTimeout   = 1 * time.Second
	commonRequestTimeout = 2 * time.Second
)

type apiHandler struct {
	orderService service.OrderService
}

func NewAPIHandler(service service.OrderService) *apiHandler {
	return &apiHandler{
		orderService: service,
	}
}
