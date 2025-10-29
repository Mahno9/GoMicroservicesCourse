package service

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

type PartService interface {
	GetPart(c context.Context, uuid string) (*model.Part, error)
	ListParts(c context.Context, filter *model.PartsFilter) ([]*model.Part, error)
	InitWithDummy() error
}
