package models

import (
	"database/sql"
	"guilliman/migrations"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitializeDatabase() {
	// Check for the development mode using the DEV environment variable
	devMode := os.Getenv("DEV") == "true"

	// Set database path based on the mode
	databasePath := "/data/database.db"
	if devMode {
		databasePath = "./dev_database.db"
		log.Println("Running in development mode. Using local database:", databasePath)
	} else {
		// Ensure the /data directory exists in production mode
		dataDir := "/data"
		err := os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
		log.Println("Running in production mode.")
	}

	// Open the SQLite database
	var err error
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	log.Printf("Connected to SQLite database at %s", databasePath)
}

func CreateTables() error {
	// User table
	userTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	// Accounts table
	accountsTable := `CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name TEXT UNIQUE NOT NULL,
		type TEXT NOT NULL,
		currency TEXT NOT NULL,
		balance REAL DEFAULT 0,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Categories table
	categoryTable := `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name TEXT UNIQUE NOT NULL,
		main_category TEXT NOT NULL,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Transactions table
	transactionsTable := `CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		description TEXT NOT NULL,
		amount REAL NOT NULL,
		currency TEXT NOT NULL,
		amount_in_base_currency REAL,
		exchange_rate REAL,
		date INTEGER NOT NULL,
		main_category TEXT NOT NULL,
		subcategory TEXT NOT NULL,
		category_id INTEGER,
		user_id INTEGER,
		account_id INTEGER,
		related_account_id INTEGER,
		transaction_type TEXT NOT NULL,
		fees INTEGER DEFAULT 0,
		FOREIGN KEY (category_id) REFERENCES categories(id),
		FOREIGN KEY (account_id) REFERENCES accounts(id),
		FOREIGN KEY (related_account_id) REFERENCES accounts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Migrations table
	migrationsTable := `CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Execute table creation scripts
	tableStatements := []string{userTable, accountsTable, categoryTable, transactionsTable, migrationsTable}
	for _, stmt := range tableStatements {
		if _, err := db.Exec(stmt); err != nil {
			log.Printf("Failed to create table: %v", err)
			return err
		}
	}

	log.Println("Tables created successfully")
	return nil
}

func ApplyMigrations() error {

	for name, stmt := range migrations.Migrations {
		// Check if this migration has already been applied
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = ?)", name).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			log.Printf("Migration %s already applied", name)
			continue
		}

		// Apply the migration
		log.Printf("Applying migration %s...", name)
		if _, err := db.Exec(stmt); err != nil {
			log.Printf("Failed to apply migration %s: %v", name, err)
			return err
		}

		// Record the migration
		_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", name)
		if err != nil {
			return err
		}
	}

	log.Println("All migrations applied successfully")
	return nil
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
	initialCategories := []struct {
		Name         string
		MainCategory string
	}{
		{"Emergency Fund", "Savings"},
		{"Investments", "Savings"},
		{"Debt Repayment", "Savings"},
		{"Short Term", "Savings"},
		{"Travels", "Savings"},
		{"Savings", "Savings"},
		{"Interests Earned", "Savings"},
		{"Pets", "Needs"},
		{"House Services", "Needs"},
		{"Bank Fees", "Needs"},
		{"Groceries", "Needs"},
		{"Rent", "Needs"},
		{"Utilities", "Needs"},
		{"Transportation", "Needs"},
		{"Work Lunchs", "Needs"},
		{"Streaming Services", "Wants"},
		{"Health", "Wants"},
		{"Leisure", "Wants"},
		{"Self Care", "Wants"},
		{"Entertainment", "Wants"},
		{"Shopping", "Wants"},
		{"Hobbies", "Wants"},
		{"Taxi", "Wants"},
		{"Restaurants", "Wants"},
		{"Transfer", "Transfer"},
		{"Salary", "Income"},
		{"Interests", "Income"},
		{"Payments", "Income"},
	}

	for _, category := range initialCategories {
		_, err := db.Exec("INSERT INTO categories (name, main_category) VALUES (?, ?)", category.Name, category.MainCategory)
		if err != nil {
			return err
		}
	}

	log.Println("Categories table successfully seeded")
	return nil
}

func CloseDatabase() {
	if db != nil {
		db.Close()
	}
}

func ClearDatabase() error {
	_, err := db.Exec("DELETE FROM accounts")
	if err != nil {
		log.Printf("Error clearing accounts table: %v", err)
		return err
	}

	_, err = db.Exec("DELETE FROM categories")
	if err != nil {
		log.Printf("Error clearing categories table: %v", err)
		return err
	}

	_, err = db.Exec("DELETE FROM transactions")
	if err != nil {
		log.Printf("Error clearing transactions table: %v", err)
		return err
	}

	log.Println("Database cleared successfully")
	return nil
}
