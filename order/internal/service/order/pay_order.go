package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderData model.PayOrderData) (uuid.UUID, error) {
	reqGetCtx, cancelGet := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancelGet()

	order, err := s.getOrderFromRepo(reqGetCtx, orderData)
	if err != nil {
		return uuid.Nil, err
	}

	// There are no data about user in pay request. Enrich with it
	orderData.UserUuid = order.UserUuid

	transactionUuid, err := s.payOrder(ctx, orderData)
	if err != nil {
		return uuid.Nil, err
	}

	err = s.updateRepoWithPaymentData(ctx, order, &orderData, transactionUuid)
	if err != nil {
		return uuid.Nil, err
	}

	// TODO: add test for this
	err = s.produceOrderPaidEvent(ctx, order)
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUuid, nil
}

func (s *service) getOrderFromRepo(reqGetCtx context.Context, orderData model.PayOrderData) (*model.Order, error) {
	order, err := s.orderRepository.Get(reqGetCtx, orderData.OrderUuid)
	if err != nil {
		return nil, err
	}
	if order.Status != model.StatusPENDINGPAYMENT {
		return nil, model.ErrOrderCancelConflict
	}

	return order, nil
}

func (s *service) payOrder(ctx context.Context, orderData model.PayOrderData) (uuid.UUID, error) {
	reqPayCtx, cancelPay := context.WithTimeout(ctx, model.RequestTimeoutUpdate)
	defer cancelPay()

	return s.paymentClient.PayOrder(reqPayCtx, orderData)
}

func (s *service) updateRepoWithPaymentData(ctx context.Context, order *model.Order, orderData *model.PayOrderData, transactionUuid uuid.UUID) error {
	order.TransactionUuid = &transactionUuid
	order.Status = model.StatusPAID
	order.PaymentMethod = orderData.PaymentMethod

	reqUpdateCtx, cancelUpdate := context.WithTimeout(ctx, model.RequestTimeoutUpdate)
	defer cancelUpdate()

	return s.orderRepository.Update(reqUpdateCtx, order)
}

func (s *service) produceOrderPaidEvent(ctx context.Context, order *model.Order) error {
	err := s.producerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		Uuid:            uuid.New().String(),
		OrderUuid:       order.TransactionUuid.String(),
		UserUuid:        order.UserUuid.String(),
		PaymentMethod:   string(converter.ModelToApiPaymentMethod(order.PaymentMethod)),
		TransactionUuid: order.TransactionUuid.String(),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", model.ErrKafkaSend, err)
	}

	return nil
}
