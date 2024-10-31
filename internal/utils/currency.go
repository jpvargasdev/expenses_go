package utils

import (
	"fmt"
	"guilliman/config"
	"net/http"
	"time"
)

var (
  exchangeRates   map[string]float64
  lastFetchTime   time.Time
  baseCurrency    = "SEK"
  cacheCurrency   = "SEK"
  apiKey          string
  exchangeRateURL string
)

type ExchangeRateResponse struct {
  Result             string             `json:"result"`
  BaseCode           string             `json:"base_code"`
  ConversionRates    map[string]float64 `json:"conversion_rates"`
  TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
}

/**
* TODO: Cache improvements:
* - store multiple base currencies
* - store multiple exchangeRates per currency
*/
func FetchExchangeRates(currency string) error {

  // check time
  if (time.Since(lastFetchTime) < time.Hour && exchangeRates != nil)  || (currency != cacheCurrency) {
    cacheCurrency = currency
    return nil
  }

  exchangeRateURL := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", config.GetExchangeRateKey(), currency)
  
  res, err := http.Get(exchangeRateURL)
  if err != nil {
    cacheCurrency = currency
    return fmt.Errorf("failed to fetch exchange rates: %v", err)
  }

  defer res.Body.Close()

  var data ExchangeRateResponse

  if res.StatusCode == http.StatusOK {
     cacheCurrency = currency
    exchangeRates = data.ConversionRates
    lastFetchTime = time.Now()
    return nil
  }

  return nil
}


func GetExchangeRate(currency string) (float64, error) {
  if (currency == baseCurrency) {
    return 1.0, nil
  }

  if exchangeRates == nil {
    if err := FetchExchangeRates(currency); err != nil {
      return 0, err
    }
  }

  rate, exists := exchangeRates[currency]
  if !exists {
    return 0, fmt.Errorf("exchange rate not found for currency: %s", currency)
  }

  return rate, nil
}

