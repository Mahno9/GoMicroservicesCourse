package model

import "errors"

var UnknownPaymentMethodErr = errors.New("unknown payment method")
var PartsNotAvailableErr = errors.New("no required parts are available")
var OrderDoesNotExistErr = errors.New("order does not exist")
