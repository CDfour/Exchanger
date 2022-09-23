package model

import (
	"time"

	"github.com/google/uuid"
)

type Currency struct {
	ID           uuid.UUID
	Name         string
	Description  string
	Abbreviation string
	IsBase       bool
}

type ExchangeRate struct {
	ID           uuid.UUID
	BaseCurrency Currency
	Currency     Currency
	TimeAdd      time.Time
	Rate         float64
}
