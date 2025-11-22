package model

import "errors"

var (
	ErrUnknownPaymentMethod = errors.New("unknown payment method")
	ErrPartsNotAvailable    = errors.New("no required parts are available")
	ErrOrderDoesNotExist    = errors.New("order does not exist")
	ErrOrderCancelConflict  = errors.New("order cancel conflict")

	ErrQueryBuild            = errors.New("query build error")
	ErrQueryExecution        = errors.New("query execution error")
	ErrQueryResponseScanning = errors.New("query response scanning error")
)
