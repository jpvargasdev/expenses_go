package models

import (
	"log"
)

func ResetDatabase() {
	if err := ClearDatabase(); err != nil {
		log.Fatalf("Failed to clear database: %v", err)
	}
}
