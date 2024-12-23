package config

import (
  "log"
  "os"
  "github.com/joho/godotenv"
)

type AppConfig struct {
  SqlDb string
  ServerPort string
  ExchangeRateKey string
}

var Config AppConfig

func Load() {
  err := godotenv.Load()
  if err != nil {
    log.Println("No .env file found, relying on environment variables")
  }

  Config.ServerPort = getEnv("SERVER_PORT", "8080")
  Config.ExchangeRateKey = getEnv("EXCHANGE_RATE_KEY", "")
  Config.SqlDb = getEnv("SQL_URL", "")
}

func GetServerPort() string {
  return Config.ServerPort
}

func GetExchangeRateKey() string {
  return Config.ExchangeRateKey
}

func GetSqlDb() string {
  return Config.SqlDb
}

func getEnv(key, defaultVal string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  }
  return defaultVal
}
