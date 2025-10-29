package payment

import (
	def "github.com/Mahno9/GoMicroservicesCourse/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
