package converter

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func ModelToAPIError(err error) error {
	if errors.Is(err, model.ErrPartNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
