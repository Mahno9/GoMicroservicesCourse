package main

import (
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

type OrderInfo struct {
	orderUuid       string
	userUuid        string
	partUuids       []string
	totalPrice      float64
	transactionUuid string
	paymentMethod   paymentV1.PaymentMethod
	status          orderV1.Status
}

type OrdersStorage struct {
	orders map[string]OrderInfo
}
