package models

import (
	"context"
	"database/sql"
	"fmt"
	"guilliman/internal/utils"
	"guilliman/internal/utils/timeutils"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/guregu/null/v5"
)

const (
	TransactionTypeIncome   = "Income"
	TransactionTypeExpense  = "Expense"
	TransactionTypeSavings  = "Savings"
	TransactionTypeTransfer = "Transfer"
)

const (
	MainCategoryNeeds    = "Needs"
	MainCategoryWants    = "Wants"
	MainCategorySavings  = "Savings"
	MainCategoryTransfer = "Transfer"
)

type Transaction struct {
	ID                   string      `json:"id"`
	Description          string      `json:"description"`
	Amount               float64     `json:"amount"`
	Currency             string      `json:"currency"`
	AmountInBaseCurrency float64     `json:"amount_in_base_currency"`
	ExchangeRate         float64     `json:"exchange_rate"`
	Date                 int64       `json:"date"`
	MainCategory         string      `json:"main_category"`
	Subcategory          string      `json:"subcategory"`
	CategoryID           null.String `json:"category_id"`
	AccountID            null.String `json:"account_id"`
	RelatedAccountID     null.String `json:"related_account_id"`
	TransactionType      string      `json:"transaction_type"`
	Fees                 float64     `json:"fees"`
	UserID               string      `json:"user_id"`
}

func GetTransactionsByMainCategory(mainCategory string, startDay string, endDay string, uid string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var start, end int64

	startDate, endDate := timeutils.GetSalaryMonthRange(startDay, endDay)
	start = startDate.Unix()
	end = endDate.Unix()

	query := `SELECT 
	  id,	
	  description,	
	  amount,	
	  currency,	
	  amount_in_base_currency,	
	  exchange_rate,	
	  main_category,	
	  subcategory,	
	  date,	
	  category_id,
	  account_id,
	  related_account_id,
	  transaction_type
	FROM transactions`

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "user_id = $1")
	args = append(args, uid)

	if mainCategory != "" {
		conditions = append(conditions, "main_category = $2")
		args = append(args, mainCategory)
	}

	conditions = append(conditions, "date BETWEEN $3 AND $4")
	args = append(args, start, end)

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionByID(transactionID string, userID string) (Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction Transaction

	query := `SELECT 
		id, description, amount, currency, amount_in_base_currency, exchange_rate, 
		date, main_category, subcategory, category_id, account_id, 
		related_account_id, transaction_type
	FROM transactions 
	WHERE id = $1 AND user_id = $2`

	err := db.QueryRow(ctx, query, transactionID, userID).Scan(
		&transaction.ID,
		&transaction.Description,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.AmountInBaseCurrency,
		&transaction.ExchangeRate,
		&transaction.Date,
		&transaction.MainCategory,
		&transaction.Subcategory,
		&transaction.CategoryID,
		&transaction.AccountID,
		&transaction.RelatedAccountID,
		&transaction.TransactionType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Transaction{}, fmt.Errorf("transaction not found")
		}
		return Transaction{}, fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	return transaction, nil
}

func GetTransactions(transactionType string, accountId string, limitParam string, uid string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT 
	  id,	
	  description,	
	  amount,	
	  currency,	
	  amount_in_base_currency,	
	  exchange_rate,	
	  main_category,	
	  subcategory,	
	  date,	
	  category_id,
	  account_id,
	  related_account_id,
	  transaction_type
	FROM transactions`

	var conditions []string
	var args []interface{}

	// Always filter by user_id
	conditions = append(conditions, "user_id = $1")
	args = append(args, uid)
	argIndex := 2 // Track query argument numbers

	if transactionType != "" {
		conditions = append(conditions, fmt.Sprintf("transaction_type = $%d", argIndex))
		args = append(args, transactionType)
		argIndex++
	}

	if accountId != "" {
		conditions = append(conditions, fmt.Sprintf("account_id = $%d", argIndex))
		args = append(args, accountId)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY date DESC" // Ensure newest transactions are retrieved first

	if limitParam != "" {
		_, err := strconv.Atoi(limitParam)
		if err == nil { // Ensure limit is a valid integer
			query += fmt.Sprintf(" LIMIT $%d", argIndex)
			args = append(args, limitParam)
		}
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionsForPeriod(start int64, end int64, transactionType string, accountId string, uid string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT 
	  id,	
	  description,	
	  amount,	
	  currency,	
	  amount_in_base_currency,	
	  exchange_rate,	
	  main_category,	
	  subcategory,	
	  date,	
	  category_id,
	  account_id,
	  related_account_id,
	  transaction_type
	FROM transactions`

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "user_id = $1")
	args = append(args, uid)

	if transactionType != "" {
		conditions = append(conditions, "transaction_type = $2")
		args = append(args, transactionType)
	}

	if accountId != "" {
		conditions = append(conditions, "account_id = $3")
		args = append(args, accountId)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

/**
* Add a new transaction to the database
* Can add TransactionType = "expense", "income"
 */
func AddTransaction(transaction Transaction) (Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sourceAccount, err := GetAccountByID(transaction.AccountID, transaction.UserID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid account: %v", err)
	}

	if transaction.TransactionType == TransactionTypeExpense {
		if sourceAccount.Balance < transaction.Amount {
			return Transaction{}, fmt.Errorf("insufficient balance in account: %v", err)
		}
	}

	// Determine the main category based on the subcategory
	var categoryID string
	if transaction.CategoryID.Valid {
		categoryID = transaction.CategoryID.String
	} else {
		categoryID = "" // Handle empty case appropriately
	}

	mainCategory, err := GetMainCategory(categoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	subcategory, err := GetSubCategory(categoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	transaction.MainCategory = mainCategory
	transaction.Subcategory = subcategory

	if transaction.Date == 0 {
		transaction.Date = time.Now().Unix()
	}

	var exchangeRate float64
	var amountInBaseCurrency float64

	rate, err := utils.GetExchangeRate(transaction.Currency)
	if err != nil {
		// Log the error but proceed without exchange rate
		log.Printf("Warning: Exchange rate not found for currency '%s'. Transaction will be saved without conversion.", transaction.Currency)
		exchangeRate = 0
		amountInBaseCurrency = 0
	} else {
		exchangeRate = rate
		// Convert the transaction amount to the base currency
		amountInBaseCurrency = transaction.Amount * exchangeRate
	}

	transaction.ExchangeRate = exchangeRate
	transaction.AmountInBaseCurrency = amountInBaseCurrency

	// Start a database transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Insert the transaction into the database
	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (
		  description,
		  amount,
		  currency,
		  amount_in_base_currency,
		  exchange_rate,
		  date,
		  main_category,
		  subcategory,
		  category_id,
		  account_id,
		  related_account_id,
		  transaction_type,
      user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		transaction.Description,
		transaction.Amount,
		transaction.Currency,
		transaction.AmountInBaseCurrency,
		transaction.ExchangeRate,
		transaction.Date,
		transaction.MainCategory,
		transaction.Subcategory,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.RelatedAccountID,
		transaction.TransactionType,
		transaction.UserID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to insert transaction: %v", err)
	}

	// Update the account balance for the source account
	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
		transaction.Amount, transaction.AccountID,
	)

	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to update source account balance: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return transaction, nil
}

func UpdateTransaction(transactionID string, updatedTransaction Transaction) (Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve the existing transaction
	existingTransaction, err := GetTransactionByID(transactionID, updatedTransaction.UserID)
	if err != nil {
		return Transaction{}, fmt.Errorf("transaction not found: %v", err)
	}

	// Ensure the category exists
	var categoryID string
	if updatedTransaction.CategoryID.Valid {
		categoryID = updatedTransaction.CategoryID.String
	} else {
		categoryID = "" // Handle empty case appropriately
	}
	mainCategory, err := GetMainCategory(categoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	subcategory, err := GetSubCategory(categoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	updatedTransaction.MainCategory = mainCategory
	updatedTransaction.Subcategory = subcategory

	// Handle currency conversion if the currency changed
	if updatedTransaction.Currency != existingTransaction.Currency {
		rate, err := utils.GetExchangeRate(updatedTransaction.Currency)
		if err != nil {
			log.Printf("Warning: Exchange rate not found for currency '%s'. Transaction will be saved without conversion.", updatedTransaction.Currency)
			updatedTransaction.ExchangeRate = 0
			updatedTransaction.AmountInBaseCurrency = 0
		} else {
			updatedTransaction.ExchangeRate = rate
			updatedTransaction.AmountInBaseCurrency = updatedTransaction.Amount * rate
		}
	} else {
		// Retain the previous exchange rate if the currency hasn't changed
		updatedTransaction.ExchangeRate = existingTransaction.ExchangeRate
		updatedTransaction.AmountInBaseCurrency = updatedTransaction.Amount * existingTransaction.ExchangeRate
	}

	// Start a database transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Adjust the account balance: First revert the old transaction
	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance - $1 WHERE id = $2`,
		existingTransaction.Amount, existingTransaction.AccountID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to revert old transaction amount: %v", err)
	}

	// Then apply the updated transaction amount
	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
		updatedTransaction.Amount, updatedTransaction.AccountID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to update account balance: %v", err)
	}

	// Update the transaction in the database
	_, err = tx.Exec(ctx,
		`UPDATE transactions SET
		  description = $1, amount = $2, currency = $3, amount_in_base_currency = $4, exchange_rate = $5, 
		  date = $6, main_category = $7, subcategory = $8, category_id = $9, account_id = $10, 
		  related_account_id = $11, transaction_type = $12
		WHERE id = $13`,
		updatedTransaction.Description,
		updatedTransaction.Amount,
		updatedTransaction.Currency,
		updatedTransaction.AmountInBaseCurrency,
		updatedTransaction.ExchangeRate,
		updatedTransaction.Date,
		updatedTransaction.MainCategory,
		updatedTransaction.Subcategory,
		updatedTransaction.CategoryID,
		updatedTransaction.AccountID,
		updatedTransaction.RelatedAccountID,
		updatedTransaction.TransactionType,
		transactionID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to update transaction: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return updatedTransaction, nil
}

/**
* Add a new transfer to the database
* Can add TransactionType = "transfer" "savings"
 */
func AddTransfer(transaction Transaction) (Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var categoryID string
	if transaction.CategoryID.Valid {
		categoryID = transaction.CategoryID.String
	} else {
		categoryID = "" // Handle empty case appropriately
	}

	mainCategory, err := GetMainCategory(categoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	subcategory, err := GetSubCategory(categoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	transaction.MainCategory = mainCategory
	transaction.Subcategory = subcategory

	if transaction.Date == 0 {
		transaction.Date = time.Now().Unix()
	}

	var exchangeRate float64
	var amountInBaseCurrency float64

	rate, err := utils.GetExchangeRate(transaction.Currency)
	if err != nil {
		// Log the error but proceed without exchange rate
		log.Printf("Warning: Exchange rate not found for currency '%s'. Transaction will be saved without conversion.", transaction.Currency)
		exchangeRate = 0
		amountInBaseCurrency = 0
	} else {
		exchangeRate = rate
		// Convert the transaction amount to the base currency
		amountInBaseCurrency = transaction.Amount * exchangeRate
	}

	transaction.ExchangeRate = exchangeRate
	transaction.AmountInBaseCurrency = amountInBaseCurrency

	tx, err := db.Begin(ctx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Insert the transaction into the database
	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (
		  description,
		  amount,
		  currency,
		  amount_in_base_currency,
		  exchange_rate,
		  date,
		  main_category,
		  subcategory,
		  category_id,
		  account_id,
		  related_account_id,
		  transaction_type,
		  fees,
      user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		transaction.Description,
		transaction.Amount,
		transaction.Currency,
		transaction.AmountInBaseCurrency,
		transaction.ExchangeRate,
		transaction.Date,
		transaction.MainCategory,
		transaction.Subcategory,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.RelatedAccountID,
		transaction.TransactionType,
		transaction.Fees,
		transaction.UserID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to insert transaction: %v", err)
	}

	// Update the account balance for the source account
	_, err = tx.Exec(ctx,
    `UPDATE accounts SET balance = balance - ($1::NUMERIC + $2::NUMERIC) WHERE id = $3`,
		transaction.Amount, transaction.Fees, transaction.AccountID,
	)

	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to update source account balance: %v", err)
	}

	// Update the account balance for the destination account
	_, err = tx.Exec(ctx,
    `UPDATE accounts SET balance = balance + ($1::NUMERIC + $2::NUMERIC) WHERE id = $3`,
		transaction.Amount, transaction.Fees, transaction.RelatedAccountID,
	)

	if err != nil {
		tx.Rollback(ctx)
		return Transaction{}, fmt.Errorf("failed to update destination account balance: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return transaction, nil
}

func DeleteTransaction(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start a database transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Fetch the transaction details to retrieve its amount and account ID
	var transaction Transaction 

  log.Printf("Transaction ID: %s", id)

	err = tx.QueryRow(ctx,
		`SELECT amount, account_id, related_account_id, transaction_type, fees
		 FROM transactions 
		 WHERE id = $1`, id,
	).Scan(
		&transaction.Amount,
		&transaction.AccountID,
		&transaction.RelatedAccountID,
		&transaction.TransactionType,
		&transaction.Fees,
	)

  if err != nil {
    return fmt.Errorf("Error retrieving transaction: %v", err)
  }

	if err == sql.ErrNoRows {
		tx.Rollback(ctx)
		return fmt.Errorf("transaction with ID %d not found", id)
	} else if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to retrieve transaction: %v", err)
	}

  log.Printf("Transaction ID 2: %s", id)

	// Reverse the balance change for the source account
	_, err = tx.Exec(ctx,
    `UPDATE accounts SET balance = balance + ($1::NUMERIC + $2::NUMERIC) WHERE id = $3`,
		transaction.Amount, transaction.Fees, transaction.AccountID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to update source account balance: %v", err)
	}

	// If the transaction is a transfer, update the related account balance as well
	if transaction.TransactionType == TransactionTypeTransfer ||
		transaction.TransactionType == TransactionTypeSavings &&
		transaction.RelatedAccountID.Valid && transaction.RelatedAccountID.String != "" {
		_, err = tx.Exec(ctx,
      `UPDATE accounts SET balance = balance - ($1::NUMERIC + $2::NUMERIC) WHERE id = $3`,
			transaction.Amount, transaction.Fees, transaction.RelatedAccountID,
		)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to update destination account balance: %v", err)
		}
	}

	// Delete the transaction
	_, err = tx.Exec(ctx, `DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to delete transaction: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return nil
}

func GetTransactionsByAccount(accountID string, uid string) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Prepare the SQL query
	query := `
		SELECT 
			id, 
			description, 
			amount, 
			currency, 
			amount_in_base_currency, 
			exchange_rate, 
			date, 
			main_category, 
			subcategory, 
			category_id, 
			account_id, 
			related_account_id, 
			transaction_type, 
			fees 
		FROM transactions 
		WHERE account_id = $1 AND user_id = $2
	`

	// Initialize a slice to hold the transactions
	var transactions []Transaction

	// Execute the query
	rows, err := db.Query(ctx, query, accountID, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer rows.Close() // Ensure rows are closed

	// Iterate through the result set and scan into the slice
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.Date,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
			&transaction.Fees,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return transactions, nil
}
