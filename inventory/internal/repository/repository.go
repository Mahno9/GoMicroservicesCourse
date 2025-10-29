package repository

import (
	"context"

	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

type PartRepository interface {
	GetPart(c context.Context, uuid string) (*domainModel.Part, error)
	ListParts(c context.Context, filter *domainModel.PartsFilter) ([]*domainModel.Part, error)
	InitWithDummy() error
}
