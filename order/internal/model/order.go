package model

import (
	"github.com/google/uuid"
)

type CreateOrderData struct {
	UserUuid  uuid.UUID
	PartUuids []uuid.UUID
}

type PayOrderData struct {
	OrderUuid     uuid.UUID
	UserUuid      uuid.UUID
	PaymentMethod int32
}

type Order struct {
	OrderUuid       uuid.UUID
	UserUuid        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float64
	TransactionUuid *uuid.UUID
	PaymentMethod   int32
	Status          OrderStatus
}

type OrderStatus int32

const (
	StatusPENDINGPAYMENT OrderStatus = 0
	StatusPAID           OrderStatus = 1
	StatusCANCELLED      OrderStatus = 2
)
