package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) PayOrder(c context.Context, orderData model.PayOrderData) (string, error) {
	order, err := s.ordersRepo.Get(orderData.OrderUuid)
	if err != nil {
		return "", err
	}

	if order.Status != model.StatusPENDINGPAYMENT {
		return "", model.ErrOrderCancelConflict
	}

	orderData.UserUuid = order.UserUuid

	transactionUuid, err := s.payment.PayOrder(c, orderData)
	if err != nil {
		return "", err
	}

	order.TransactionUuid = transactionUuid
	order.Status = model.StatusPAID
	order.PaymentMethod = orderData.PaymentMethod

	err = s.ordersRepo.Update(order)
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
