package migrations

var Migrations = map[string]string{
	"001_add_user_id_to_accounts": `
			-- Step 1: Create a new table with the updated schema
			CREATE TABLE new_accounts (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					name TEXT UNIQUE NOT NULL,
					type TEXT NOT NULL,
					currency TEXT NOT NULL,
					balance REAL DEFAULT 0,
					user_id INTEGER,
					FOREIGN KEY (user_id) REFERENCES users(id)
			);

			-- Step 2: Copy data from the old table to the new table
			INSERT INTO new_accounts (id, created_at, name, type, currency, balance)
			SELECT id, created_at, name, type, currency, balance FROM accounts;

			-- Step 3: Drop the old table
			DROP TABLE accounts;

			-- Step 4: Rename the new table to the original table name
			ALTER TABLE new_accounts RENAME TO accounts;
	`,
	// Add similar migrations for categories and transactions
	"002_add_user_id_to_categories": `
			CREATE TABLE new_categories (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					name TEXT UNIQUE NOT NULL,
					main_category TEXT NOT NULL,
					user_id INTEGER,
					FOREIGN KEY (user_id) REFERENCES users(id)
			);
			INSERT INTO new_categories (id, created_at, name, main_category)
			SELECT id, created_at, name, main_category FROM categories;
			DROP TABLE categories;
			ALTER TABLE new_categories RENAME TO categories;
	`,
	"003_add_user_id_to_transactions": `
			CREATE TABLE new_transactions (
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
			);
			INSERT INTO new_transactions (id, created_at, description, amount, currency, amount_in_base_currency, exchange_rate, date, main_category, subcategory, category_id, account_id, related_account_id, transaction_type, fees)
			SELECT id, created_at, description, amount, currency, amount_in_base_currency, exchange_rate, date, main_category, subcategory, category_id, account_id, related_account_id, transaction_type, fees FROM transactions;
			DROP TABLE transactions;
			ALTER TABLE new_transactions RENAME TO transactions;
	`,
}
