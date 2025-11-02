package v1

import (
	"context"
	"net/http"

	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	// TODO: handle errors from order/internal/model/error.go

	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
