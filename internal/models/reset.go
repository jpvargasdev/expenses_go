package models

import (
	"log"
)

func ResetDatabase() error {
	err := ClearDatabase()
	if err == nil {
		log.Fatalf("Failed to clear database: %v", err)
	}

	return err
}
