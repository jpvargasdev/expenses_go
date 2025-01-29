package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
)

var FirebaseApp *firebase.App

// InitFirebase initializes the Firebase app with the service account
func InitFirebase() error {
	// Initialize Firebase app
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
		return err
	}

	FirebaseApp = app
	log.Println("Firebase initialized successfully")
	return nil
}
