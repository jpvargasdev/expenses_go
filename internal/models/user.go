package models

import "fmt"

type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	DisplayName   string `json:"display_name"`
  PhoneNumber   string `json:"phone_number"`
  PhotoUrl      string `json:"photo_url`
}

func CreateUser(user User) error {
  // Check if the user already exists in the database
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", user.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("User already exists") 
	}

	// Insert the new user into the database
	query := `
		INSERT INTO users (id, email, display_name, phone_number, photo_url) 
		VALUES (?, ?)
	`
	_, err = db.Exec(query, user.ID, user.Email, user.DisplayName, user.PhoneNumber, user.PhotoUrl)
	if err != nil {
		return err
	}

  return nil
}

func DeleteUser(uid string) error {
  _, err := db.Exec("DELETE from users WHERE id = ?", uid)
  if err != nil {
    return err
  }
  return nil
}
