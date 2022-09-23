package exchange

import "errors"

var (
	ErrNoApikey  = errors.New("отсутствует apikey")
	ErrStatusBad = errors.New("ошибка статуса ответа")
)
