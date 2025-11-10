package model

type CreateOrderData struct {
	UserUuid  string
	PartUuids []string
}

type PayOrderData struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod int32
}

type Order struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid string
	PaymentMethod   int32
	Status          OrderStatus
}

type OrderStatus int32

const (
	StatusPENDINGPAYMENT OrderStatus = 0
	StatusPAID           OrderStatus = 1
	StatusCANCELLED      OrderStatus = 2
)
