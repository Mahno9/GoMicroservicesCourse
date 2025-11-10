package repository

import (
	"context"

	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, uuid string) (*domainModel.Part, error)
	ListParts(ctx context.Context, filter *domainModel.PartsFilter) ([]*domainModel.Part, error)
	InitWithDummy() error
}
