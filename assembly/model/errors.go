package model

import "errors"

var (
	ErrKafkaInvalidInputData = errors.New("kafka invalid input data")
	ErrKafkaSend             = errors.New("kafka send error")
)
