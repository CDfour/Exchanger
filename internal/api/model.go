package api

import (
	"time"

	"github.com/google/uuid"
)

type ErrorStruct struct {
	Status string
}

type Currency struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Abbreviation string    `json:"abbreviation"`
}

type ExchangeRate struct {
	ID           uuid.UUID `json:"id"`
	BaseCurrency Currency  `json:"base_currency"`
	Currency     Currency  `json:"currency"`
	TimeAdd      time.Time `json:"time"`
	Rate         float64   `json:"rate"`
}
