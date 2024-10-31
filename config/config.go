package config

import (
  "log"
  "os"
  "github.com/joho/godotenv"
)

type AppConfig struct {
  ServerPort string
  ExchangeRate string
}

var Config AppConfig

func Load() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  Config.ServerPort = getEnv("SERVER_PORT", "8080")
  Config.ExchangeRate = getEnv("EXCHANGE_RATE_KEY", "")
}

func GetServerPort() string {
  return Config.ServerPort
}

func GetExchangeRateKey() string {
  return Config.ExchangeRate
}

func getEnv(key, defaultVal string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  }
  return defaultVal
}
