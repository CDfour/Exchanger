package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"project/internal/model"

	"github.com/google/uuid"
)

type client struct {
	apikey string
}

func NewExchangeClient(apikey string) (*client, error) {
	if apikey == "" {
		return nil, ErrNoApikey
	}
	return &client{apikey}, nil
}

func (c *client) GetRates(currencies []model.Currency, baseCurrency model.Currency) (rates []model.ExchangeRate, err error) {

	var exchangeCurrencies Currencies

	symbols := getSymbols(currencies)
	url := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/latest?symbols=%s&base=%s", symbols, baseCurrency.Abbreviation)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", c.apikey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrStatusBad
	}

	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&exchangeCurrencies); err != nil {
		return nil, err
	}

	for abb, rate := range exchangeCurrencies.Rates {
		for _, currency := range currencies {
			if abb == currency.Abbreviation {
				rates = append(rates, model.ExchangeRate{
					ID:           uuid.New(),
					BaseCurrency: baseCurrency,
					Currency:     currency,
					TimeAdd:      time.Unix(int64(exchangeCurrencies.Timestamp), 0),
					Rate:         rate,
				})
				break
			}
		}
	}
	return rates, nil
}

func getSymbols(currencies []model.Currency) string {

	var abbreviations []string

	for _, currency := range currencies {
		abbreviations = append(abbreviations, currency.Abbreviation)
	}

	return strings.Join(abbreviations, ",")
}
