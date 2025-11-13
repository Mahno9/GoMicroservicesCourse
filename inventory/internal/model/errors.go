package model

import "errors"

var (
	ErrPartNotFound      = errors.New("part not found")
	ErrDbIndexInitFailed = errors.New("db index init failed")
)
