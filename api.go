package exchangerate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	config "github.com/tommzn/go-config"
	core "github.com/tommzn/hdb-datasource-core"
	events "github.com/tommzn/hdb-events-go"
)

func newExchangeRateApi(conf config.Config) (core.DataSource, error) {

	apiUrl := conf.Get("exchangerate.url", nil)
	if apiUrl == nil {
		return nil, errors.New("Exchange rate api not spefified.")
	}
	dateFormat := conf.Get("exchangerate.date_formar", config.AsStringPtr("2006-01-02"))
	return &ExchangeRateApi{
		apiUrl:        *apiUrl,
		dateFormat:    *dateFormat,
		exchangeRates: exchangeRatesFromConfig(conf),
	}, nil
}

// exchangeRatesFromConfig extracts a list of requested currency conversion from given config.
func exchangeRatesFromConfig(conf config.Config) []*events.ExchangeRate {

	rates := []*events.ExchangeRate{}
	for _, rate := range conf.GetAsSliceOfMaps("exchangerate.conversions") {
		fromCurrency, ok1 := rate["from"]
		toCurrency, ok2 := rate["to"]
		if ok1 && ok2 {
			rates = append(rates, &events.ExchangeRate{FromCurrency: fromCurrency, ToCurrency: toCurrency})
		}
	}
	return rates
}

// Fetch calls the specified exchange rate api to get a conversion rate for all requested currency pairs.
func (client *ExchangeRateApi) Fetch() (proto.Message, error) {

	rates := &events.ExchangeRates{Rates: []*events.ExchangeRate{}}
	for _, exchangeRate := range client.exchangeRates {
		exchangeRate, err := client.requestRate(exchangeRate)
		if err != nil {
			return nil, err
		}
		rates.Rates = append(rates.Rates, exchangeRate)
	}
	return rates, nil
}

func (client *ExchangeRateApi) requestRate(exchangeRate *events.ExchangeRate) (*events.ExchangeRate, error) {

	resp, err := http.Get(fmt.Sprintf("%s?from=%s&to=%s", client.apiUrl, exchangeRate.FromCurrency, exchangeRate.ToCurrency))
	if err != nil {
		return nil, fmt.Errorf("Unable to request exchange rate for %s/%s, reason: %s", exchangeRate.FromCurrency, exchangeRate.ToCurrency, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected HTTP status code %d\n", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to response body, reason: %s\n", err)
	}

	apiResponse := apiResponse{}
	decodeErr := json.Unmarshal(body, &apiResponse)
	if decodeErr != nil {
		return nil, fmt.Errorf("Unable to process api response, reason: %s", decodeErr)
	}

	if len(apiResponse.Rates) == 0 {
		return nil, fmt.Errorf("No rates available for %s/%s", exchangeRate.FromCurrency, exchangeRate.ToCurrency)
	}

	rate, ok := apiResponse.Rates[exchangeRate.ToCurrency]
	if !ok {
		return nil, fmt.Errorf("Target rate not available for %s/%s", exchangeRate.FromCurrency, exchangeRate.ToCurrency)
	}

	exchangeRate.Rate = rate
	exchangeRate.Timestamp = asTimeStamp(parseDate(apiResponse.Date, client.dateFormat))
	return exchangeRate, nil

}

// parseDate tries to pass passed date with given foramt. If parsing fails, it will return current Time.
func parseDate(date, format string) time.Time {

	t, err := time.Parse(format, date)
	if err == nil {
		return t
	}
	return time.Now()
}

// asTimeStamp converts a unix epoch timestamp to a Protobuf timestamp.
func asTimeStamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
