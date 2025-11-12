package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderData model.PayOrderData) (string, error) {
	reqGetCtx, cancelGet := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancelGet()

	order, err := s.orderRepository.Get(reqGetCtx, orderData.OrderUuid)
	if err != nil {
		return "", err
	}

	if order.Status != model.StatusPENDINGPAYMENT {
		return "", model.ErrOrderCancelConflict
	}

	orderData.UserUuid = order.UserUuid

	reqPayCtx, cancelPay := context.WithTimeout(ctx, model.RequestTimeoutUpdate)
	defer cancelPay()

	transactionUuid, err := s.paymentClient.PayOrder(reqPayCtx, orderData)
	if err != nil {
		return "", err
	}

	order.TransactionUuid = &transactionUuid
	order.Status = model.StatusPAID
	order.PaymentMethod = orderData.PaymentMethod

	reqUpdateCtx, cancelUpdate := context.WithTimeout(ctx, model.RequestTimeoutUpdate)
	defer cancelUpdate()

	err = s.orderRepository.Update(reqUpdateCtx, order)
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
