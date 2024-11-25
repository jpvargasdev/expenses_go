package models

import (
	"time"
  "fmt"
  "log"
  "database/sql"

  "guilliman/internal/utils"
)

type Income struct {
  ID                    int     `json:"id"`
  Description           string  `json:"description"`
  Amount                float64 `json:"amount"`        // Amount in transaction currency
  Currency              string  `json:"currency"`      // Currency code of the transaction
  AmountInBaseCurrency  float64 `json:"amount_in_base_currency"`   // Amount converted to base currency
  ExchangeRate          float64 `json:"exchange_rate"`             // Exchange rate used for conversion
  Date                  int64   `json:"date"`                      // Unix timestamp
}


func GetIncomes() ([]Income, error) {
  rows, err := db.Query(`
    SELECT 
      id, 
      description, 
      amount, 
      currency, 
      amount_in_base_currency, 
      exchange_rate, 
      date, 
    FROM incomes
  `)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var incomes []Income
  for rows.Next() {
    var income Income
    if err := rows.Scan(
      &income.ID,
      &income.Description,
      &income.Amount,
      &income.Currency,
      &income.AmountInBaseCurrency,
      &income.ExchangeRate,
      &income.Date,
    ); err != nil {
      return nil, err
    }
    incomes = append(incomes, income)
  }
  return incomes, nil 
}

func AddIncome(income Income) error {
  if income.Date == 0 {
    income.Date = time.Now().Unix()
  }

  var exchangeRate float64
  var amountInBaseCurrency float64

  rate, err := utils.GetExchangeRate(income.Currency)
  if err != nil {
    // Log the error but proceed without exchange rate
    log.Printf("Warning: Exchange rate not found for currency '%s'. Income will be saved without conversion.", income.Currency)
    exchangeRate = 0
    amountInBaseCurrency = 0
  } else {
    exchangeRate = rate
    // Convert the transaction amount to the base currency
    amountInBaseCurrency = income.Amount * exchangeRate
  }

  income.ExchangeRate = exchangeRate
  income.AmountInBaseCurrency = amountInBaseCurrency

  // Insert the expense into the database
  _, err = db.Exec(
    `INSERT INTO expenses (
      description,
      amount,
      currency,
      amount_in_base_currency,
      exchange_rate,
      date,
    ) VALUES (?, ?, ?, ?, ?, ?)`,
    income.Description,
    income.Amount,
    income.Currency,
    income.AmountInBaseCurrency, // May be zero or null
    income.ExchangeRate,         // May be zero or null
    income.Date,
  )
  if err != nil {
    return err
  }

  return nil
}

func DeleteIncome(id int) error {
  result, err := db.Exec("DELETE FROM incomes WHERE id = ?", id)
  if err != nil {
    return fmt.Errorf("could not delete income: %v", err)
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return fmt.Errorf("could not retrieve affected rows: %v", err)
  }

  if rowsAffected == 0 {
    return sql.ErrNoRows
  }

  return nil
}

