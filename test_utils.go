package exchangerate

import (
	config "github.com/tommzn/go-config"
)

// loadConfigForTest loads test config.
func loadConfigForTest(configFile *string) config.Config {

	if configFile == nil {
		configFile = config.AsStringPtr("fixtures/testconfig.yml")
	}
	configLoader := config.NewFileConfigSource(configFile)
	config, _ := configLoader.Load()
	return config
}
