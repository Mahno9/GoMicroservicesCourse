package model

type Order struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid string
	PaymentMethod   int32
	Status          int32
}
