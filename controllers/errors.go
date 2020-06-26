package controllers

import (
	"errors"
)

var (
	ErrStringNotContains     = errors.New("string does not contain the required substring")
	ErrUnexpectedSliceLength = errors.New("unexpected slice length")
	ErrMapKeyNotFound        = errors.New("map does not contain the expected key")
)
