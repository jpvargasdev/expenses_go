package main

import (
	"fmt"
	"guilliman/config"
	"guilliman/internal/models"
	"guilliman/internal/routes"
  "guilliman/cmd/firebase"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title           Guilliman API
// @version         1.0
// @description     This is the Guilliman API

// @contact.name   Juan Vargas
// @contact.url    https://github.com/jpvargasdev/guilliman

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
func main() {
	config.Load()

	models.InitializeDatabase()
	defer models.CloseDatabase()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("\nShutting down server...")
		models.CloseDatabase()
		os.Exit(0)
	}()

	if err := models.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Seed the database with initial categories
	if err := models.SeedCategories(); err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}

  // Init Firebase
  err := firebase.InitFirebase()
  if err != nil {
    fmt.Printf("Failed to initialize Firebase: %v", err)
  }

	router := routes.SetupRouter()

	port := config.GetServerPort()

	fmt.Printf("Guilliman server is running on port %s...\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting Guilliman server: %v", err)
	}
}
