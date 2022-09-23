package config

import "errors"

var (
	ErrNoApikey   = errors.New("отсутствует apikey")
	ErrNoDatabase = errors.New("отсутствует Database")
	ErrNoHost     = errors.New("отсутсвует Host")
	ErrNoPassword = errors.New("отсутствует Password")
	ErrNoPort     = errors.New("отсутствует Port")
	ErrNoSchedule = errors.New("отсутвует Schedule")
	ErrNoUsername = errors.New("отсутствует Username")
)
