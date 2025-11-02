package order

import (
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"

	repoModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	orders map[string]*repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]*repoModel.Order),
	}
}
