package model

type CreateOrderData struct {
	UserUuid  string
	PartUuids map[string]any
}

type PayOrderData struct {
	PaymentMethod int32
	OrderUuid     string
}

type OrderInfo struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid string
	PaymentMethod   int32
	Status          int32
}

type OrdersStorage struct {
	Orders map[string]OrderInfo
}
