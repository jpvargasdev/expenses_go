package utils

import (
	"errors"
	"guilliman/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int) (string, error) {
  claims := jwt.MapClaims{
    "user_id": userID,
    "exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Should we add more time? maybe not expiring time??
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(config.GetSecretKey())
}

func ValidateToken(tokenString string) error {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return config.GetSecretKey(), nil
  })

  if err != nil {
    return err
  }

  if !token.Valid {
    return errors.New("invalid token")
  }

  return nil
}
