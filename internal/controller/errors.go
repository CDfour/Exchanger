package controller

import "errors"

var (
	ErrNoRepository = errors.New("отсутсвует Repository")
	ErrNoExchange   = errors.New("отсутвует Exchange")
)
