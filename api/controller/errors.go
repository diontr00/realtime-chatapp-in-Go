package controller

import "errors"

var (
	EvenNotSupportedError = errors.New("This event type is not supported")
	InvalidMessageError   = errors.New("Invalid message")
)
