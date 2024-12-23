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

var db *sql.DB

func InitializeDatabase() {
	const remoteDBURL = config.GetSqlDb()
	
	if remoteDBURL != "" {
		// Connect to the hosted SQLite database
		db, err = sql.Open("sqlite3", remoteDBURL)
		if err != nil {
			log.Fatalf("Failed to open remote SQLite database: %v", err)
		}
		log.Println("Connected to remote SQLite database:", remoteDBURL)
	} else {
		// Fallback to local SQLite file
		db, err = sql.Open("sqlite3", "./guilliman.db")
		if err != nil {
			log.Fatalf("Failed to open local SQLite database: %v", err)
		}
		log.Println("Connected to local SQLite database: guilliman.db")
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
		main_category TEXT NOT NULL        -- Main category (Needs, Wants, Savings, Transfer)
	);`

	transactionsTable := `CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		amount REAL NOT NULL,                        -- Positive for income, negative for expenses
		currency TEXT NOT NULL,
		amount_in_base_currency REAL,
		exchange_rate REAL,
		date INTEGER NOT NULL,
		main_category TEXT NOT NULL,                 -- Needs, Wants, Savings
		subcategory TEXT NOT NULL,                   -- Name of the subcategory
		category_id INTEGER,
		account_id INTEGER,                          -- Account from which the transaction is made
		related_account_id INTEGER,                  -- Account to which the transaction is made (for transfers)
		transaction_type TEXT NOT NULL,              -- 'Expense', 'Income', 'Savings', 'Transfer'
		fees INTEGER DEFAULT 0,                      -- Fees associated with the transaction
		FOREIGN KEY (category_id) REFERENCES categories(id),
		FOREIGN KEY (account_id) REFERENCES accounts(id),
		FOREIGN KEY (related_account_id) REFERENCES accounts(id)
	);`

	accountsTable := `CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,   -- Name of the account (e.g., "Checking Account", "Credit Card")
		type TEXT NOT NULL,          -- Type of account (e.g., "Bank", "Credit Card", "Cash")
		currency TEXT NOT NULL,      -- Currency of the account (e.g., "USD", "EUR")
		balance REAL DEFAULT 0       -- Current balance of the account (optional)
	);`

	_, err := db.Exec(categoryTable)
	if err != nil {
		log.Fatalf("Failed to create categories table: %v", err)
	}

	_, err = db.Exec(transactionsTable)
	if err != nil {
		log.Fatalf("Failed to create transactions table: %v", err)
	}

	_, err = db.Exec(accountsTable)
	if err != nil {
		log.Fatalf("Failed to create accounts table: %v", err)
	}

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
