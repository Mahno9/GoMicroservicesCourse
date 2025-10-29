package part

import (
	repository "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository"
	def "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	repository repository.PartRepository
}

func NewService(repository repository.PartRepository) *service {
	return &service{
		repository: repository,
	}
}
