package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"project/internal/model"
)

// Структура репозитория
type repository struct {
	dbConn *pgx.Conn
}

// Конструктор репозитория
func NewRepository(dbConn *pgx.Conn) (*repository, error) {
	if dbConn == nil {
		return nil, ErrNoConn
	}

	return &repository{dbConn}, nil
}

// Добавление строки в таблицу
func (r *repository) AddRates(exchangeRates []model.ExchangeRate) error {
	query := `
		INSERT INTO public.exchange (id, currency_id, base_currency_id, time_add, rate)
		VALUES
		($1,$2,$3,$4,$5);
	`
	for _, exchangeRate := range exchangeRates {
		_, err := r.dbConn.Exec(
			context.Background(),
			query,
			exchangeRate.ID,
			exchangeRate.Currency.ID,
			exchangeRate.BaseCurrency.ID,
			exchangeRate.TimeAdd,
			exchangeRate.Rate,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Возвращает список доступных валют
func (r *repository) GetCurrencies() ([]model.Currency, error) {
	query := `	SELECT id, 
		  			   name, 
					   description, 
					   abbreviation, 
					   is_base
				FROM public.currency;`
	rows, err := r.dbConn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var (
		currency   repositoryCurrency
		currencies []model.Currency
	)

	for rows.Next() {
		err = rows.Scan(
			&currency.ID,
			&currency.Name,
			&currency.Description,
			&currency.Abbreviation,
			&currency.IsBase,
		)
		if err != nil {
			return nil, err
		}

		currencies = append(currencies, model.Currency{
			ID:           currency.ID,
			Name:         currency.Name,
			Description:  currency.Description,
			Abbreviation: currency.Abbreviation,
			IsBase:       currency.IsBase,
		})
	}
	return currencies, nil
}

func (r *repository) GetCurrency(abb string) (*model.Currency, error) {
	currency := repositoryCurrency{}

	query := `SELECT id,
			  		 name,
			 		 description,
			 		 abbreviation,
			 		 is_base
		FROM public.currency
		WHERE abbreviation = $1;`

	if err := r.dbConn.QueryRow(context.Background(), query, abb).Scan(
		&currency.ID,
		&currency.Name,
		&currency.Description,
		&currency.Abbreviation,
		&currency.IsBase,
	); err != nil {
		return nil, err
	}

	return &model.Currency{
		ID:           currency.ID,
		Name:         currency.Name,
		Description:  currency.Description,
		Abbreviation: currency.Abbreviation,
		IsBase:       currency.IsBase,
	}, nil
}

// Возвращает курс для указанной валюты в определенное время
func (r *repository) GetRates(currencies []model.Currency, baseCurrency model.Currency, t time.Time) ([]model.ExchangeRate, error) {

	var (
		repositoryExchange repositoryExchange
		exchangeRates      []model.ExchangeRate
	)

	query := `	SELECT id,
					   currency_id,
					   base_currency_id,
					   time_add,
					   rate
				FROM public.exchange
				WHERE currency_id = $1 and time_add <= $2
				ORDER BY time_add desc
				LIMIT 1;`

	for _, currency := range currencies {
		if err := r.dbConn.QueryRow(context.Background(), query, currency.ID, t).Scan(
			&repositoryExchange.ID,
			&repositoryExchange.CurrencyID,
			&repositoryExchange.BaseID,
			&repositoryExchange.TimeAdd,
			&repositoryExchange.Rate,
		); err != nil {
			return nil, err
		}
		exchangeRates = append(exchangeRates, model.ExchangeRate{
			ID:           repositoryExchange.ID,
			BaseCurrency: baseCurrency,
			Currency:     currency,
			TimeAdd:      repositoryExchange.TimeAdd,
			Rate:         repositoryExchange.Rate,
		})
	}

	return exchangeRates, nil
}
