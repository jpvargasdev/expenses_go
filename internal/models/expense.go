package models

import (
	"time"
  "fmt"
)

type Expense struct {
  ID           int     `json:"id"`
  Description  string  `json:"description"`
  Amount       float64 `json:"amount"`
  MainCategory string  `json:"main_category"`  // Needs, Wants, or Savings
  Subcategory  string  `json:"subcategory"`    // Specific subcategory name (e.g., Groceries, Rent)
  Date         int64   `json:"date"`           // Unix timestamp for the expense date
  CategoryID   int     `json:"category_id"`
}

type ExpenseWithCategory struct {
  ID           int     `json:"id"`
  Description  string  `json:"description"`
  Amount       float64 `json:"amount"`
  MainCategory string  `json:"main_category"`  // Needs, Wants, or Savings
  Subcategory  string  `json:"subcategory"`    // Specific subcategory name
  Date         int64   `json:"date"`           // Unix timestamp for the expense date
  CategoryID   int     `json:"category_id"`
  CategoryName string  `json:"category_name"`  // The name of the category from the categories table
}

func GetExpenses() ([]ExpenseWithCategory, error) {
  rows, err := db.Query(`
      SELECT expenses.id, expenses.description, expenses.amount, expenses.main_category, expenses.subcategory, expenses.date, expenses.category_id, categories.name 
      FROM expenses 
      LEFT JOIN categories ON expenses.category_id = categories.id
    `)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var expenses []ExpenseWithCategory
  for rows.Next() {
    var expense ExpenseWithCategory
    if err := rows.Scan(
      &expense.ID,
      &expense.Description,
      &expense.Amount,
      &expense.MainCategory,
      &expense.Subcategory,
      &expense.Date,
      &expense.CategoryID,
      &expense.CategoryName,
    ); err != nil {
      return nil, err
    }
    expenses = append(expenses, expense)
  }
  return expenses, nil 
}

func GetExpensesForPeriod(start, end int64) ([]ExpenseWithCategory, error) {
  rows, err := db.Query(`
      SELECT expenses.id, expenses.description, expenses.amount, expenses.main_category, expenses.subcategory, expenses.date, expenses.category_id, categories.name 
      FROM expenses 
      LEFT JOIN categories ON expenses.category_id = categories.id
      WHERE expenses.date >= ? AND expenses.date <= ?
    `, start, end)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

 var expenses []ExpenseWithCategory
  for rows.Next() {
      var expense ExpenseWithCategory
      if err := rows.Scan(
        &expense.ID, 
        &expense.Description, 
        &expense.Amount, 
        &expense.MainCategory,
        &expense.Subcategory,
        &expense.Date, 
        &expense.CategoryID, 
        &expense.CategoryName,
      ); err != nil {
        return nil, err
      }
      expenses = append(expenses, expense)
  }
  return expenses, nil
}

func AddExpense(expense Expense) error {
  mainCategory, err := GetMainCategory(expense.Subcategory)
  if err != nil {
      return fmt.Errorf("invalid subcategory: %v", err)
  }
  expense.MainCategory = mainCategory

  expense.Date = time.Now().Unix()
  _, err = db.Exec(
      "INSERT INTO expenses (description, amount, main_category, subcategory, date, category_id) VALUES (?, ?, ?, ?, ?, ?)",
      expense.Description, expense.Amount, expense.MainCategory, expense.Subcategory, expense.Date, expense.CategoryID,
  )
  return err
}

