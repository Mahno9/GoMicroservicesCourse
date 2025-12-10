package model

import "errors"

var (
	ErrFailedToSendNotification = errors.New("failed to send notification")
	ErrInvalidEvent             = errors.New("invalid event")
	ErrTelegramClient           = errors.New("telegram client error")
)
