package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) NewError(_ context.Context, err error) *genOrderV1.GenericErrorStatusCode {
	// Трансляция ошибок из модели в HTTP ошибки
	switch {
	case errors.Is(err, model.ErrUnknownPaymentMethod):
		return &genOrderV1.GenericErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: genOrderV1.GenericError{
				Code:    genOrderV1.NewOptInt(http.StatusBadRequest),
				Message: genOrderV1.NewOptString("Unknown payment method"),
			},
		}
	case errors.Is(err, model.ErrPartsNotAvailable):
		return &genOrderV1.GenericErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response: genOrderV1.GenericError{
				Code:    genOrderV1.NewOptInt(http.StatusBadRequest),
				Message: genOrderV1.NewOptString("No required parts are available"),
			},
		}
	case errors.Is(err, model.ErrOrderDoesNotExist):
		return &genOrderV1.GenericErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response: genOrderV1.GenericError{
				Code:    genOrderV1.NewOptInt(http.StatusNotFound),
				Message: genOrderV1.NewOptString("Order does not exist"),
			},
		}
	case errors.Is(err, model.ErrOrderCancelConflict):
		return &genOrderV1.GenericErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response: genOrderV1.GenericError{
				Code:    genOrderV1.NewOptInt(http.StatusConflict),
				Message: genOrderV1.NewOptString("Order cancel conflict"),
			},
		}
	default:
		// Для всех остальных ошибок возвращаем внутреннюю ошибку сервера
		return &genOrderV1.GenericErrorStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response: genOrderV1.GenericError{
				Code:    genOrderV1.NewOptInt(http.StatusInternalServerError),
				Message: genOrderV1.NewOptString(err.Error()),
			},
		}
	}
}
