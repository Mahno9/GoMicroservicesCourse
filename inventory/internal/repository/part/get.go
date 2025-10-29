package part

import (
	"context"

	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, partUuid string) (*domainModel.Part, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	part, ok := r.parts[partUuid]
	if !ok {
		return nil, domainModel.ErrPartNotFound
	}

	return converter.RepoToDomainPart(part), nil
}
