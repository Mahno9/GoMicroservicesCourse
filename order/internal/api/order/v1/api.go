package v1

import (
	"context"
	"log"
	"time"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	model "github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/service"

	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

type apiHandler struct {
	orderService service.OrderService
}

func NewAPIHandler() {

}

const (
	createOrderTimeout = 1 * time.Second
)

func (h *apiHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderReq) (orderV1.CreateOrderRes, error) {
	log.Printf("Creating order with details: %v\n", req)

	timedCtx, cancel := context.WithTimeout(ctx, createOrderTimeout)
	defer cancel()

	modelReqData := converter.ApiToModelOrderInfo(req)
	response, err := h.orderService.CreateOrder(timedCtx, modelReqData)

	// response, err := h.inventory.ListParts(timedCtx, &inventoryV1.ListPartsRequest{
	// 	Filter: &inventoryV1.PartsFilter{
	// 		Uuids: req.PartUuids,
	// 	},
	// })
	if err != nil {
		log.Printf("❗ Failed to get parts: %v\nNo order is created.", err)
		return nil, err
	}

	if len(response.Parts) != len(req.PartUuids) {
		log.Printf("❗ Some of requested parts are absent")
		return nil, model.PartsNotAvailableErr
	}

	totalPrice := float64(0.0)
	for _, part := range response.Parts {
		totalPrice += part.GetPrice()
	}

	orderUUID := uuid.New().String()
	h.store.orders[orderUUID] = OrderInfo{
		orderUuid:  orderUUID,
		userUuid:   string(req.UserUUID),
		partUuids:  req.PartUuids,
		totalPrice: totalPrice,
		status:     orderV1.StatusPENDINGPAYMENT,
	}

	log.Printf("❕ New order created:\n%+v\n", h.store.orders[orderUUID])

	return &orderV1.CreateOrderCreated{
		OrderUUID:  orderV1.OrderUUID(orderUUID),
		TotalPrice: orderV1.TotalPrice(totalPrice),
	}, nil
}

func (h *apiHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	for _, order := range h.store.orders {
		if order.orderUuid == string(params.OrderUUID) {
			var paymentMethod orderV1.PaymentMethod
			err := paymentMethod.UnmarshalText([]byte(order.paymentMethod.String()))
			if err != nil {
				paymentMethod = orderV1.PaymentMethodUNKNOWN
			}

			return &orderV1.GetOrderOK{
				OrderUUID:       orderV1.OrderUUID(order.orderUuid),
				UserUUID:        orderV1.UserUUID(order.userUuid),
				PartUuids:       order.partUuids,
				TotalPrice:      orderV1.TotalPrice(order.totalPrice),
				TransactionUUID: orderV1.TransactionUUID(order.transactionUuid),
				PaymentMethod:   paymentMethod,
				Status:          order.status,
			}, nil
		}
	}

	return &orderV1.NotFoundError{}, nil
}

func (h *apiHandler) OrderCancel(ctx context.Context, params orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	order, ok := h.store.orders[params.OrderUUID]
	if !ok {
		return &orderV1.OrderCancelNotFound{}, nil
	}

	if order.status != orderV1.StatusPENDINGPAYMENT {
		return &orderV1.OrderCancelConflict{}, nil
	}

	order.status = orderV1.StatusCANCELLED
	h.store.orders[params.OrderUUID] = order

	return &orderV1.OrderCancelNoContent{}, nil
}

func (h *apiHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderReq, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, ok := h.store.orders[params.OrderUUID]
	if !ok {
		return nil, nil
	}

	if order.status != orderV1.StatusPENDINGPAYMENT {
		log.Printf("❗ Invalid order status (%s). Unable to make payment.\n", order.status)
		return &orderV1.PayOrderConflict{}, nil
	}

	paymentMethod := convertPaymentMethod(&req.PaymentMethod)

	payResp, err := h.payment.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID,
		UserUuid:      order.userUuid,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		log.Printf("❗ Failed to pay order: %v\n", err)
		return nil, err
	}

	order.paymentMethod = paymentMethod
	order.transactionUuid = payResp.TransactionUuid
	order.status = orderV1.StatusPAID

	h.store.orders[order.orderUuid] = order

	return &orderV1.PayOrderOK{
		TransactionUUID: orderV1.TransactionUUID(payResp.TransactionUuid),
	}, nil
}

func convertPaymentMethod(orderPaymentMethod *orderV1.PaymentMethod) paymentV1.PaymentMethod {
	paymentMethodBytes, err := orderPaymentMethod.MarshalText()
	if err != nil {
		log.Printf("❗ Unknown payment method: %v\n", orderPaymentMethod)
		paymentMethodBytes = []byte(paymentV1.PaymentMethod_UNKNOWN.String())
	}

	return paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(paymentMethodBytes)])
}

func (h *apiHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
