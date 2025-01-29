package models

import (
	"fmt"
	"log"
)

type Account struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`     // Name of the account (e.g., "Checking Account", "Credit Card")
	Type     string  `json:"type"`     // Type of account (e.g., "Bank", "Credit Card", "Cash")
	Currency string  `json:"currency"` // Currency of the account (e.g., "USD", "EUR")
	Balance  float64 `json:"balance"`  // Balance of the account (optional)
  UserID   string  `json:"user_id"`
}

func GetAccounts(id string, uid string) ([]Account, error) {
	query := "SELECT id, name, type, currency, balance FROM accounts"

  query += fmt.Sprintf(" WHERE user_id = '%s'", uid)

	if id != "" {
		query += " AND WHERE id = ?"
	}

	rows, err := db.Query(query, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.Type,
			&account.Currency,
			&account.Balance,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func GetAccountByID(id int, uid string) (Account, error) {
	query := "SELECT id, name, type, currency, balance FROM accounts WHERE id = ? AND user_id = ?"

	var account Account
	if err := db.QueryRow(query, id, uid).Scan(
		&account.ID,
		&account.Name,
		&account.Type,
		&account.Currency,
		&account.Balance,
	); err != nil {
		return Account{}, err
	}

	return account, nil
}

func AddAccount(acccount Account) (Account, error) {
	result, err := db.Exec(`INSERT INTO accounts (
		name,
		type,
		currency,
		balance,
    user_id
	) VALUES (?, ?, ?, ?, ?)`,
		acccount.Name,
		acccount.Type,
		acccount.Currency,
		acccount.Balance,
    acccount.UserID,
	)
	if err != nil {
		return Account{}, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Println("Warning: Could not retrieve last insert ID for account")
	} else {
		acccount.ID = int(lastID)
	}

	return acccount, nil
}

func UpdateAccount(account Account) (Account, error) {
	query := `
		UPDATE accounts
		SET
			name = COALESCE(?, name),
			type = COALESCE(?, type),
			currency = COALESCE(?, currency),
			balance = COALESCE(?, balance)
		WHERE id = ?`

	// Execute the query
	result, err := db.Exec(query, account.Name, account.Type, account.Currency, account.Balance, account.ID)
	if err != nil {
		return Account{}, fmt.Errorf("failed to update account: %v", err)
	}

	// Check if the account was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Account{}, fmt.Errorf("failed to fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return Account{}, fmt.Errorf("no account found with ID %d", account.ID)
	}

	return account, nil
}

func DeleteAccount(id int, uid string) (Account, error) {
	query := "DELETE FROM accounts WHERE id = ? AND user_id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		return Account{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Account{}, err
	}

	if rowsAffected == 0 {
		return Account{}, fmt.Errorf("no account found with ID %d", id)
	}

	return Account{ID: id}, nil
}
