package controller

import (
	"time"

	"project/internal/model"

	"github.com/sirupsen/logrus"
)

// Интерфейс репозитория
type IRepository interface {
	AddRates(exchangeRates []model.ExchangeRate) error
	GetCurrencies() ([]model.Currency, error)
	GetRates(currencies []model.Currency, baseCurrency model.Currency, t time.Time) ([]model.ExchangeRate, error)
	GetCurrency(abb string) (*model.Currency, error)
}

// Интерфейс биржи
type IExchange interface {
	GetRates(currencies []model.Currency, baseCurrency model.Currency) (rates []model.ExchangeRate, err error)
}

// Структура контроллера
type controller struct { // приватный
	repository IRepository
	exchange   IExchange
}

// Конструктор котроллера
func NewController(repository IRepository, exchange IExchange) (*controller, error) {
	if repository == nil {
		return nil, ErrNoRepository
	}
	if exchange == nil {
		return nil, ErrNoExchange
	}
	return &controller{repository, exchange}, nil
}

// Добавление курсов с биржи в базу данных
func (c *controller) AddRates() {
	logrus.Info("starting AddRates")

	currencies, err := c.repository.GetCurrencies()
	if err != nil {
		logrus.Errorln("repository.GetCurrencies: ", err)
		return
	}

	var (
		baseCurrency          model.Currency
		CurrenciesWithoutBase []model.Currency
	)
	for _, currency := range currencies {
		if currency.IsBase {
			baseCurrency = currency
		} else {
			CurrenciesWithoutBase = append(CurrenciesWithoutBase, currency)
		}
	}

	exchangeRates, err := c.exchange.GetRates(CurrenciesWithoutBase, baseCurrency)
	if err != nil {
		logrus.Errorln("exchange.GetRates: ", err)
		return
	}

	err = c.repository.AddRates(exchangeRates)
	if err != nil {
		logrus.Errorln("repository.AddRates: ", err)
		return
	}

	logrus.Info("ending AddRates")
}

// Вызывает у репозитория метод возвращения всех доступных валют и затем возвращает этот результат в API
func (c *controller) GetCurrencies() ([]model.Currency, error) {
	currencies, err := c.repository.GetCurrencies()
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

// Вызывает у репозитория метод возвращения курса для указанной валюты в определенное время и возвращает этот результат в API
func (c *controller) GetRates(abbreviations []string, t time.Time) ([]model.ExchangeRate, error) {
	var currencies []model.Currency
	for _, abb := range abbreviations {
		currency, err := c.repository.GetCurrency(abb)
		if err != nil {
			return nil, err
		}

		currencies = append(currencies, *currency)
	}

	var baseCurrency model.Currency

	allCurrencies, err := c.repository.GetCurrencies()
	if err != nil {
		return nil, err
	}

	for _, currency := range allCurrencies {
		if currency.IsBase {
			baseCurrency = currency
		}
	}

	exchangeRates, err := c.repository.GetRates(currencies, baseCurrency, t)
	if err != nil {
		return nil, err
	}

	return exchangeRates, nil
}
