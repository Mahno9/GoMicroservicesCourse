package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	// Трансляция ошибок из модели в HTTP ошибки
	switch {
	case errors.Is(err, model.UnknownPaymentMethodErr):
		return &orderV1.GenericErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: orderV1.GenericError{
				Code:    orderV1.NewOptInt(http.StatusBadRequest),
				Message: orderV1.NewOptString("Unknown payment method"),
			},
		}
	case errors.Is(err, model.PartsNotAvailableErr):
		return &orderV1.GenericErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: orderV1.GenericError{
				Code:    orderV1.NewOptInt(http.StatusBadRequest),
				Message: orderV1.NewOptString("No required parts are available"),
			},
		}
	case errors.Is(err, model.OrderDoesNotExistErr):
		return &orderV1.GenericErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response: orderV1.GenericError{
				Code:    orderV1.NewOptInt(http.StatusNotFound),
				Message: orderV1.NewOptString("Order does not exist"),
			},
		}
	case errors.Is(err, model.OrderCancelConflictErr):
		return &orderV1.GenericErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response: orderV1.GenericError{
				Code:    orderV1.NewOptInt(http.StatusConflict),
				Message: orderV1.NewOptString("Order cancel conflict"),
			},
		}
	default:
		// Для всех остальных ошибок возвращаем внутреннюю ошибку сервера
		return &orderV1.GenericErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response: orderV1.GenericError{
				Code:    orderV1.NewOptInt(http.StatusInternalServerError),
				Message: orderV1.NewOptString(err.Error()),
			},
		}
	}
}
