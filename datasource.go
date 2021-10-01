package exchangerate

import (
	config "github.com/tommzn/go-config"
	core "github.com/tommzn/hdb-datasource-core"
)

// New returns a new weather datasource
func New(config config.Config) (core.DataSource, error) {
	return newExchangeRateApi(config)

}
