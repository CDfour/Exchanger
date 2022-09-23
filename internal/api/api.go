package api

import (
	"net/http"
	"strings"
	"time"

	"project/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Интерфейс для контроллера
type IController interface {
	GetCurrencies() ([]model.Currency, error)
	GetRates([]string, time.Time) ([]model.ExchangeRate, error)
}

// Структура API
type api struct {
	controller IController
}

// Конструктор структуры API
func NewAPI(controller IController) (*api, error) { // проверка на nil
	if controller == nil {
		return nil, ErrNoController
	}

	return &api{controller}, nil
}

// @Summary      currencies
// @Description  Get name, description and abbreviation for all currencies in database
// @Tags         Currencies
// @Produce      json
// @Success 200 {object} Currency
// @Failure 500 {object} ErrorStruct
// @Router      /currencies/ [get]
// Выводит названия, описания и аббревиатуры всех валют из базы данных
func (a *api) CurrenciesHandler(c *gin.Context) {
	logrus.Info("starting CurrenciesHandler")

	currencies, err := a.controller.GetCurrencies()
	if err != nil {
		logrus.Errorln("controller.GetCurrencies: ", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorStruct{"Error"})
		return
	}

	var apiCurrencies []Currency
	for _, currency := range currencies {
		apiCurrencies = append(apiCurrencies, Currency{ID: currency.ID, Name: currency.Name, Description: currency.Description, Abbreviation: currency.Abbreviation})
	}

	c.IndentedJSON(http.StatusOK, apiCurrencies)

	logrus.Info("ending CurrenciesHandler")
}

// @Summary rates
// @Description Get name, description, abbreviation for currency, name, description and abbreviation for base currency, time and rate
// @Tags Rates
// @Produce json
// @Success 201 {object} ExchangeRate
// @Failure 400 {object} ErrorStruct
// @Failure 500 {object} ErrorStruct
// @Router /rates/ [get]
// Выводит курсы указанных валют в указанное время или в текущий момент
func (a *api) RateHandler(c *gin.Context) {
	logrus.Info("starting RateHandler")

	str := c.Request.URL.Query().Get("currencies")
	abbreviations := strings.Split(str, ",")

	str = c.Request.URL.Query().Get("time")
	var t time.Time
	var err error
	if str == "" {
		t = time.Now()
	} else {
		t, err = time.Parse("2006-01-02 15:04:05", str)
		if err != nil {
			logrus.Errorln("time.Parse: ", err)
			c.IndentedJSON(http.StatusBadRequest, ErrorStruct{"Wrong time"})
			return
		}
	}

	logrus.Infof("abbreviation: %s, time %s", abbreviations, t)

	controllerExchanges, err := a.controller.GetRates(abbreviations, t)
	if err != nil {
		logrus.Errorln("controller.GetRates: ", err)
		c.IndentedJSON(http.StatusInternalServerError, ErrorStruct{"Error"})
		return
	}

	var apiExchangeRates []ExchangeRate
	for _, controllerExchange := range controllerExchanges {
		apiExchangeRates = append(apiExchangeRates, ExchangeRate{
			ID: controllerExchange.ID,
			BaseCurrency: Currency{
				ID:           controllerExchange.BaseCurrency.ID,
				Name:         controllerExchange.BaseCurrency.Name,
				Description:  controllerExchange.BaseCurrency.Description,
				Abbreviation: controllerExchange.BaseCurrency.Abbreviation,
			},
			Currency: Currency{
				ID:           controllerExchange.Currency.ID,
				Name:         controllerExchange.Currency.Name,
				Description:  controllerExchange.Currency.Description,
				Abbreviation: controllerExchange.Currency.Abbreviation,
			},
			TimeAdd: controllerExchange.TimeAdd,
			Rate:    controllerExchange.Rate,
		})

	}

	c.IndentedJSON(http.StatusOK, apiExchangeRates)

	logrus.Info("ending RatesHandler")
}
