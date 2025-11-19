package model

const (
	// Order status constants
	StatusPaid           = "PAID"
	StatusCancelled      = "CANCELLED"
	StatusPendingPayment = "PENDING_PAYMENT"

	// Payment method constants
	PaymentMethodCard       = "CARD"
	PaymentMethodSBP        = "SBP"
	PaymentMethodCreditCard = "CREDIT_CARD"
	PaymentMethodInvestor   = "INVESTOR_MONEY"
	PaymentMethodUnknown    = "UNKNOWN"
)
