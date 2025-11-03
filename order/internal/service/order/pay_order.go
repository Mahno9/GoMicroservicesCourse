package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) PayOrder(c context.Context, data model.PayOrderData) (string, error) {
	order, err := s.ordersRepo.Get(data.OrderUuid)
	if err != nil {
		return "", err
	}

	if order.Status != model.StatusPENDINGPAYMENT {
		return "", model.ErrOrderCancelConflict
	}

	transactionUuid, err := s.payment.PayOrder(c, model.PayOrderData{
		OrderUuid:     data.OrderUuid,
		UserUuid:      order.UserUuid,
		PaymentMethod: data.PaymentMethod,
	})
	if err != nil {
		return "", err
	}

	order.TransactionUuid = transactionUuid
	order.Status = model.StatusPAID
	order.PaymentMethod = data.PaymentMethod

	err = s.ordersRepo.Update(order)
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
