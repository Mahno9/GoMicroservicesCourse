package part

import (
	"sync"

	def "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mut  sync.RWMutex
	data map[string]*repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: map[string]*repoModel.Part{},
		mut:  sync.RWMutex{},
	}
}

func (s *repository) Add(part *repoModel.Part) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.data[part.Uuid] = part
}

func (s *repository) Get(uuid string) (*repoModel.Part, bool) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	v, ok := s.data[uuid]
	return v, ok
}

func (s *repository) Contains(uuid string) bool {
	s.mut.RLock()
	defer s.mut.RUnlock()
	_, ok := s.data[uuid]
	return ok
}

func (s *repository) GetAll() []*repoModel.Part {
	s.mut.RLock()
	defer s.mut.RUnlock()
	parts := make([]*repoModel.Part, 0, len(s.data))
	for _, part := range s.data {
		parts = append(parts, part)
	}
	return parts
}

func (s *repository) Count() int {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return len(s.data)
}
