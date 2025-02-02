package models

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/null/v5"
	"github.com/jackc/pgx/v5"
)

// Account struct
type Account struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`     // Name of the account (e.g., "Checking Account", "Credit Card")
	Type     string  `json:"type"`     // Type of account (e.g., "Bank", "Credit Card", "Cash")
	Currency string  `json:"currency"` // Currency of the account (e.g., "USD", "EUR")
	Balance  float64 `json:"balance"`  // Balance of the account (optional)
	UserID   string  `json:"user_id"`
}

// GetAccounts retrieves all accounts for a user
func GetAccounts(id string, uid string) ([]Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, name, type, currency, balance FROM accounts WHERE user_id = $1"
	args := []interface{}{uid}

	if id != "" {
		query += " AND id = $2"
		args = append(args, id)
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.Type, &account.Currency, &account.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// GetAccountByID retrieves a single account by ID and user ID
func GetAccountByID(id null.String, uid string) (Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, name, type, currency, balance FROM accounts WHERE id = $1 AND user_id = $2"

	var account Account
	err := db.QueryRow(ctx, query, id, uid).Scan(
		&account.ID,
		&account.Name,
		&account.Type,
		&account.Currency,
		&account.Balance,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Account{}, fmt.Errorf("no account found with ID %d", id)
		}
		return Account{}, err
	}

	return account, nil
}

// AddAccount inserts a new account into the database
func AddAccount(account Account) (Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO accounts (name, type, currency, balance, user_id)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.QueryRow(ctx, query, account.Name, account.Type, account.Currency, account.Balance, account.UserID).Scan(&account.ID)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

// UpdateAccount updates an existing account
func UpdateAccount(account Account) (Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE accounts
		SET name = COALESCE(NULLIF($1, ''), name),
			type = COALESCE(NULLIF($2, ''), type),
			currency = COALESCE(NULLIF($3, ''), currency),
			balance = COALESCE(NULLIF($4, 0), balance)
		WHERE id = $5 AND user_id = $6`

	result, err := db.Exec(ctx, query, account.Name, account.Type, account.Currency, account.Balance, account.ID, account.UserID)
	if err != nil {
		return Account{}, fmt.Errorf("failed to update account: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return Account{}, fmt.Errorf("no account found with ID %d", account.ID)
	}

	return account, nil
}

// DeleteAccount removes an account from the database
func DeleteAccount(id string, uid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM accounts WHERE id = $1 AND user_id = $2"

	result, err := db.Exec(ctx, query, id, uid)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no account found with ID %d", id)
	}

	return nil
}
