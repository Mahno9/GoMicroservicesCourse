package model

import "errors"

var ErrPartNotFound = errors.New("part not found")
var ErrDbIndexInitFailed = errors.New("db index init failed")
