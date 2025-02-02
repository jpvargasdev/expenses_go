package models

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	PhoneNumber string `json:"phone_number"`
	PhotoUrl    string `json:"photo_url"`
}

// CreateUser inserts a new user if they do not exist
func CreateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the user already exists
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", user.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return nil // User already exists, nothing to do
	}

	// Insert new user
	query := `
		INSERT INTO users (id, email, display_name, phone_number, photo_url) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = db.Exec(ctx, query, user.ID, user.Email, user.DisplayName, user.PhoneNumber, user.PhotoUrl)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// DeleteUser removes a user by their UID
func DeleteUser(uid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, "DELETE FROM users WHERE id = $1", uid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
