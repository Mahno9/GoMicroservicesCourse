package part

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
)

func (s *service) GetPart(c context.Context, uuid string) (*model.Part, error) {
	return s.repository.GetPart(c, uuid)
}
