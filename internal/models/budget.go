package models

import (
	"fmt"
	"guilliman/internal/utils/timeutils"
)

type BudgetSummary struct {
	TotalIncome       float64 `json:"total_income"`
	TotalExpenses     float64 `json:"total_expenses"`
	NetBalance        float64 `json:"net_balance"`
	NeedsAmount       float64 `json:"needs_amount"`
	WantsAmount       float64 `json:"wants_amount"`
	SavingsAmount     float64 `json:"savings_amount"`
	NeedsPercentage   float64 `json:"needs_percentage"`
	WantsPercentage   float64 `json:"wants_percentage"`
	SavingsPercentage float64 `json:"savings_percentage"`
	NeedsBudget       float64 `json:"needs_budget"`
	WantsBudget       float64 `json:"wants_budget"`
	SavingsBudget     float64 `json:"savings_budget"`
	NetWorth          float64 `json:"net_worth"`
}

func GetBudgetSummary(startDay string, endDay string) (BudgetSummary, error) {
	var summary BudgetSummary
	var start, end int64

	startDate, endDate := timeutils.GetSalaryMonthRange(startDay, endDay)
	start = startDate.Unix()
	end = endDate.Unix()

	// Fetch total income
	err := db.QueryRow(`
        SELECT COALESCE(SUM(amount), 0) 
        FROM transactions
        WHERE transaction_type = 'Income' 
          AND date BETWEEN ? AND ?`,
		start, end).Scan(&summary.TotalIncome)
	if err != nil {
		return summary, fmt.Errorf("failed to retrieve total income: %v", err)
	}

	// Fetch total expenses
	err = db.QueryRow(`
        SELECT COALESCE(SUM(amount), 0) 
        FROM transactions
        WHERE transaction_type = 'Expense' 
          AND date BETWEEN ? AND ?`,
		start, end).Scan(&summary.TotalExpenses)
	if err != nil {
		return summary, fmt.Errorf("failed to retrieve total expenses: %v", err)
	}
	// Convert expenses to a positive number (since stored as negative)
	summary.TotalExpenses = -summary.TotalExpenses

	// Calculate net balance
	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses

	// Calculate net worth: sum of all account balances
	err = db.QueryRow(`
        SELECT COALESCE(SUM(balance), 0) 
        FROM accounts`,
	).Scan(&summary.NetWorth)
	if err != nil {
		return summary, fmt.Errorf("failed to retrieve net worth: %v", err)
	}

	// Fetch total expenses grouped by main_category
	rows, err := db.Query(`
        SELECT main_category, COALESCE(SUM(amount), 0) 
        FROM transactions
        WHERE transaction_type = 'Expense' 
          AND date BETWEEN ? AND ?
        GROUP BY main_category`,
		start, end)
	if err != nil {
		return summary, fmt.Errorf("failed to retrieve expenses: %v", err)
	}
	defer rows.Close()

	var needsAmount, wantsAmount, savingsAmount float64

	for rows.Next() {
		var mainCategory string
		var amount float64
		if err := rows.Scan(&mainCategory, &amount); err != nil {
			return summary, fmt.Errorf("failed to scan expense row: %v", err)
		}

		amount = -amount // convert negative to positive expense

		switch mainCategory {
		case "Needs":
			needsAmount += amount
		case "Wants":
			wantsAmount += amount
		case "Savings":
			savingsAmount += amount
		}
	}

	if err := rows.Err(); err != nil {
		return summary, fmt.Errorf("error iterating expense rows: %v", err)
	}

	summary.NeedsAmount = needsAmount
	summary.WantsAmount = wantsAmount
	summary.SavingsAmount = savingsAmount

	// 50/30/20 budget allocations
	summary.NeedsBudget = summary.TotalIncome * 0.50
	summary.WantsBudget = summary.TotalIncome * 0.30
	summary.SavingsBudget = summary.TotalIncome * 0.20

	// Calculate actual percentages spent
	if summary.TotalIncome > 0 {
		summary.NeedsPercentage = (needsAmount / summary.TotalIncome) * 100
		summary.WantsPercentage = (wantsAmount / summary.TotalIncome) * 100
		summary.SavingsPercentage = (savingsAmount / summary.TotalIncome) * 100
	} else {
		summary.NeedsPercentage = 0
		summary.WantsPercentage = 0
		summary.SavingsPercentage = 0
	}

	return summary, nil
}
