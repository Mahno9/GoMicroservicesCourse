package model

type ShipAssembled struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type OrderPaid struct {
	Uuid            string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   string
	TransactionUuid string
}
