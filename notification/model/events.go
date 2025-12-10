package model

type OrderPaidEvent struct {
	OrderID         string
	UserID          string
	PaymentMethod   string
	TransactionUUID string
}

type ShipAssembledEvent struct {
	OrderID    string
	TrackingID string
	UserID     string
}
