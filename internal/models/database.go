package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategorySeed struct {
	Name         string
	MainCategory string
}

// Initial categories and subcategories
var initialCategories = []CategorySeed{
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
	{"Family Support", "Needs"},

	{"Streaming Services", "Wants"},
	{"Health", "Wants"},
	{"Leisure ", "Wants"},
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

var db *pgxpool.Pool

func InitializeDatabase() {
	dsn := os.Getenv("DATABASE_URL") // Example: "postgres://user:password@localhost:5432/dbname?sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	db, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL")
}

// CreateTables creates the necessary tables if they don't exist
func CreateTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	categoryTable := `CREATE TABLE IF NOT EXISTS categories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name TEXT UNIQUE NOT NULL,        
		main_category TEXT NOT NULL,
		user_id TEXT REFERENCES users(id) ON DELETE CASCADE
	);`

	transactionsTable := `CREATE TABLE IF NOT EXISTS transactions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		description TEXT NOT NULL,
		amount REAL NOT NULL,                       
		currency TEXT NOT NULL,
		amount_in_base_currency REAL,
		exchange_rate REAL,
		date INTEGER NOT NULL,
		main_category TEXT NOT NULL,
		subcategory TEXT NOT NULL,
		category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
		account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
		related_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
		transaction_type TEXT NOT NULL,             
		fees REAL DEFAULT 0,
		user_id TEXT REFERENCES users(id) ON DELETE CASCADE
	);`

	accountsTable := `CREATE TABLE IF NOT EXISTS accounts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name TEXT UNIQUE NOT NULL,   
		type TEXT NOT NULL,          
		currency TEXT NOT NULL,      
		balance REAL DEFAULT 0,      
		user_id TEXT REFERENCES users(id) ON DELETE CASCADE
	);`

	userTable := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		email TEXT NOT NULL UNIQUE,
		photo_url TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		display_name TEXT NOT NULL
	);`

	migrationTable := `CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	tableStatements := []string{userTable, accountsTable, categoryTable, transactionsTable, migrationTable}

	for _, stmt := range tableStatements {
		_, err := db.Exec(ctx, stmt)
		if err != nil {
			log.Printf("Failed to create table: %v", err)
			return err
		}
	}

	log.Println("Tables created successfully")
	return nil
}

// **SeedCategories: Adds default categories to the database if not present**
func SeedCategories() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if categories already exist
	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check categories count: %w", err)
	}

	// Skip seeding if categories exist
	if count > 0 {
		log.Println("✅ Categories already exist, skipping seed.")
		return nil
	}

	log.Println("🌱 Seeding categories into database...")

	// Use batch processing for efficiency
	batch := &pgx.Batch{}
	for _, category := range initialCategories {
		batch.Queue("INSERT INTO categories (name, main_category) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			category.Name, category.MainCategory)
	}

	// Execute batch insert
	br := db.SendBatch(ctx, batch)
	defer br.Close()

	for range initialCategories {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert category: %w", err)
		}
	}

	log.Println("✅ Categories successfully seeded.")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if db != nil {
		db.Close()
	}
}

// ClearDatabase removes all data from tables
func ClearDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, "TRUNCATE accounts, categories, transactions RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Printf("Error clearing database: %v", err)
		return err
	}

	log.Println("Database cleared successfully")
	return nil
}
