package models

import (
  "database/sql"
  "log"
  _ "github.com/mattn/go-sqlite3"
)

type CategorySeed struct {
  Name         string
  MainCategory string
}

// Initial categories and subcategories
var initialCategories = []CategorySeed{
  {"Groceries", "Needs"},
  {"Rent", "Needs"},
  {"Utilities", "Needs"},
  {"Transportation", "Needs"},
  {"Restaurants", "Wants"},
  {"Entertainment", "Wants"},
  {"Shopping", "Wants"},
  {"Hobbies", "Wants"},
  {"Emergency Fund", "Savings"},
  {"Investments", "Savings"},
  {"Debt Repayment", "Savings"},
}

var db *sql.DB

func InitializeDatabase() {
  var err error
  db, err = sql.Open("sqlite3", "./guilliman.db")
  if err != nil {
    log.Fatalf("Failed to open SQLite database: %v", err)
  }

  createTables()
}

func SeedCategories() error {
  // Check if the table already has data
  var count int
  err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
  if err != nil && err != sql.ErrNoRows {
    return err
  }

  // If there are already entries in the categories table, skip seeding
  if count > 0 {
    log.Println("Categories table already populated, skipping seeding")
    return nil
  }

  // Insert initial categories
  for _, category := range initialCategories {
    _, err := db.Exec("INSERT INTO categories (name, main_category) VALUES (?, ?)", category.Name, category.MainCategory)
    if err != nil {
      return err
    }
  }

  log.Println("Categories table successfully seeded")
  return nil
}

func createTables() {
  categoryTable := `CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,         -- Name of the subcategory
    main_category TEXT NOT NULL        -- Main category (Needs, Wants, Savings)
  );`

  expenseTable := `CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    transaction_amount REAL NOT NULL,
    transaction_currency TEXT NOT NULL,
    amount_in_base_currency REAL,
    exchange_rate REAL,
    main_category TEXT NOT NULL,
    subcategory TEXT NOT NULL,
    date INTEGER NOT NULL,
    category_id INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories(id)
  );`

  _, err := db.Exec(categoryTable)
  if err != nil {
    log.Fatalf("Failed to create categories table: %v", err)
  }

  _, err = db.Exec(expenseTable)
  if err != nil {
    log.Fatalf("Failed to create expenses table: %v", err)
  }
}

func CloseDatabase() {
  if db != nil {
    db.Close()
  }
}

func ClearDatabase() error {
  _, err := db.Exec("DELETE FROM expenses")
  if err != nil {
    log.Printf("Error clearing expenses table: %v", err)
    return err
  }

  _, err = db.Exec("DELETE FROM categories")
  if err != nil {
    log.Printf("Error clearing categories table: %v", err)
    return err
  }

  log.Println("Database cleared successfully")
  return nil
}
