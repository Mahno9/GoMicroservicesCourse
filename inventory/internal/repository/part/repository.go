package part

import (
	"sync"

	def "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mut   sync.RWMutex
	parts repoModel.RepoStorage
}

func NewRepository() *repository {
	return &repository{
		parts: make(repoModel.RepoStorage),
	}
}
