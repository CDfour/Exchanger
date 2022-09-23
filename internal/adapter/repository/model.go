package repository

import (
	"time"

	"github.com/google/uuid"
)

// Структура строки таблицы currency
type repositoryCurrency struct {
	ID           uuid.UUID
	Name         string
	Description  string
	Abbreviation string
	IsBase       bool
}

// Структура строки таблицы exchange
type repositoryExchange struct {
	ID         uuid.UUID
	CurrencyID uuid.UUID
	BaseID     uuid.UUID
	TimeAdd    time.Time
	Rate       float64
}
