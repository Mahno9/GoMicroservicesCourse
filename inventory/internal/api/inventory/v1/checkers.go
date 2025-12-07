package v1

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkUUID(uuidProbe string) error {
	_, err := uuid.Parse(uuidProbe)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid UUID")
	}
	return nil
}
