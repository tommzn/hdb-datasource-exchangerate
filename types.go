package exchangerate

import (
	events "github.com/tommzn/hdb-events-go"
)

// apiResponse covers reponse payload from exchange rate api.
type apiResponse struct {
	Amount       float64            `json:"amount"`
	BaseCurrency string             `json:"base"`
	Date         string             `json:"date"`
	Rates        map[string]float64 `json:"rates"`
}

// ExchangeRateApi is used as datasource to fetch exchange rates.
type ExchangeRateApi struct {

	// apiUrl speciffy the endpoint of an exchange service.
	apiUrl string

	// dateFormat defines the validity date format for an exchange rate.
	dateFormat string

	// exchangeRates is a list of currency pairs an exchange rate should be requested for.
	exchangeRates []*events.ExchangeRate
}
