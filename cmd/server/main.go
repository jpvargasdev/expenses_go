package main

import (
  "guilliman/config"
  "guilliman/internal/models"
  "guilliman/internal/routes"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"
)

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

  // Seed the database with initial categories
  if err := models.SeedCategories(); err != nil {
    log.Fatalf("Failed to seed categories: %v", err)
  }

  router := routes.SetupRouter()

  port := config.GetServerPort()

  fmt.Printf("Guilliman server is running on port %s...\n", port)

  if err := router.Run(":" + port); err != nil {
    log.Fatalf("Error starting Guilliman server: %v", err)
  }
}
