package part

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *service) ListParts(c context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	return s.repository.ListParts(c, filter)
}
