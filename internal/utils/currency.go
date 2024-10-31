package utils

import (
  "encoding/json"
	"fmt"
	"guilliman/config"
  "io"
	"net/http"
	"time"
)

var (
  exchangeRates   map[string]float64 
  lastFetchTime   time.Time
  baseCurrency    = "SEK"                 // Default currency, should I set this in .env file???
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

func FetchExchangeRates(currency string) error {
  // check time
  if (time.Since(lastFetchTime) < time.Hour && exchangeRates != nil) {
    return nil
  }

  exchangeRateURL := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", config.GetExchangeRateKey(), currency)

  res, err := http.Get(exchangeRateURL)

  if err != nil {
    return fmt.Errorf("failed to fetch exchange rates: %v", err)
  }

  if res.StatusCode != http.StatusOK {
    bodyBytes, _ := io.ReadAll(res.Body)
    return fmt.Errorf("failed to fetch exchange rates: %s", string(bodyBytes))
  }

  defer res.Body.Close()

  var data ExchangeRateResponse
  if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
    return fmt.Errorf("failed to decode exchange rate response: %v", err)
  }

  if res.StatusCode == http.StatusOK {
    exchangeRates = data.ConversionRates
    baseCurrency = currency
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

  fmt.Println(rate)

  return rate, nil
}

