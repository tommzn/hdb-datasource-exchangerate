[![Actions Status](https://github.com/tommzn/hdb-datasource-exchangerate/actions/workflows/go.image.build.yml/badge.svg)](https://github.com/tommzn/hdb-datasource-exchangerate/actions)
[![CircleCI](https://circleci.com/gh/tommzn/hdb-datasource-exchangerate/tree/main.svg?style=svg)](https://circleci.com/gh/tommzn/hdb-datasource-exchangerate/tree/main)

# Exchange Rate Collector
This package composes a [data collector](https://github.com/tommzn/hdb-datasource-core/collector.go) and [exchange rate data source](https://github.com/tommzn/hdb-datasource-exchangerate) to fetch excjange rates and publish it to a SQS queue.

## Config
This collector requires a config to get settings for an exchange rate api, AWS SQS and maybe some logging settings.

### Example 
```yaml
log:
  loglevel: error
  shipper: logzio  

hdb:
  queue: MyWExchangeRateQueue
  archive: MyEventArchiveQueue

exchangerate:
  url: "https://api.frankfurter.app/latest"
  date_format: "2006-01-02"
  conversions:
    - from: "EUR"
      to: "USD"
    - from: "USD"
      to: "EUR"

aws:
  sqs:
    region: eu-west-1
```

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
