package config

import (
  "os"
)

type AppConfig struct {
  ServerPort string
}

var Config AppConfig

func Load() {
  Config.ServerPort = getEnv("SERVER_PORT", "8080")
}

func GetServerPort() string {
  return Config.ServerPort
}

func getEnv(key, defaultVal string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  }
  return defaultVal
}
