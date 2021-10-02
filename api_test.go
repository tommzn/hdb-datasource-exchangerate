package exchangerate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	events "github.com/tommzn/hdb-events-go"
)

type ExchangeRateApiTestSuite struct {
	suite.Suite
	conf       config.Config
	dateFormat string
}

func TestExchangeRateApiTestSuite(t *testing.T) {
	suite.Run(t, new(ExchangeRateApiTestSuite))
}

func (suite *ExchangeRateApiTestSuite) SetupTest() {
	suite.conf = loadConfigForTest(nil)
	suite.dateFormat = "2006-01-02"
}

func (suite *ExchangeRateApiTestSuite) TestConvertTime() {

	now := time.Now()
	timestamp := asTimeStamp(now)
	suite.Equal(now.UTC().Format(time.RFC3339), timestamp.AsTime().Format(time.RFC3339))
}

func (suite *ExchangeRateApiTestSuite) TestParseDate() {

	dateString := "2021-02-01"
	date := parseDate(dateString, suite.dateFormat)
	suite.Equal(dateString, date.Format(suite.dateFormat))

	date2 := parseDate("xxx", suite.dateFormat)
	suite.Equal(time.Now().Format(suite.dateFormat), date2.Format(suite.dateFormat))
}

func (suite *ExchangeRateApiTestSuite) TestLoadExchangeRatesFromConfig() {

	exchangeRates := exchangeRatesFromConfig(suite.conf)
	suite.Len(exchangeRates, 2)
}

func (suite *ExchangeRateApiTestSuite) TestGetExchangeRates() {

	datasource, err := New(suite.conf)
	suite.Nil(err)
	suite.NotNil(datasource)

	message, err := datasource.Fetch()
	suite.Nil(err)
	suite.NotNil(message)

	exchangerates, ok := message.(*events.ExchangeRates)
	suite.True(ok)
	suite.Len(exchangerates.Rates, 2)
	for _, rate := range exchangerates.Rates {
		suite.True(rate.Rate != 0)
		suite.NotNil(rate.Timestamp)
	}

}

func (suite *ExchangeRateApiTestSuite) TestGetRateForInvalidCurrency() {

	configFile := "fixtures/invalid_testconfig.yml"

	datasource, err := New(loadConfigForTest(&configFile))
	suite.Nil(err)
	suite.NotNil(datasource)

	message, err := datasource.Fetch()
	suite.NotNil(err)
	suite.Nil(message)
}

func (suite *ExchangeRateApiTestSuite) TestCreateWithoutExchangeRateApiUrl() {

	configFile := "fixtures/invalid_testconfig_02.yml"

	datasource, err := New(loadConfigForTest(&configFile))
	suite.NotNil(err)
	suite.Nil(datasource)
}
