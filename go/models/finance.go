package models

import "encoding/json"

// Category classifies operating expenses and income.
type Category struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Expense is an operating expense linked to a category.
type Expense struct {
	ID         int             `json:"id"`
	Date       string          `json:"date"`
	Concept    string          `json:"concept"`
	CategoryID int             `json:"category_id"`
	Amount     float64         `json:"amount"`
	Quantity   float64         `json:"quantity"`
	Total      float64         `json:"total"`
	ReceiptID  int             `json:"receipt_id"`
	Meta       json.RawMessage `json:"meta,omitempty"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

// Income is a revenue record linked to a category.
type Income struct {
	ID         int             `json:"id"`
	Date       string          `json:"date"`
	Concept    string          `json:"concept"`
	CategoryID int             `json:"category_id"`
	Amount     float64         `json:"amount"`
	Quantity   float64         `json:"quantity"`
	Total      float64         `json:"total"`
	ReceiptID  int             `json:"receipt_id"`
	Meta       json.RawMessage `json:"meta,omitempty"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

// FuelExpense links a fuel expense to the expense record and the driver.
type FuelExpense struct {
	ID        int    `json:"id"`
	ExpenseID int    `json:"id_expense"`
	DriverID  int    `json:"id_driver"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Stats is an aggregated weekly financial summary.
type Stats struct {
	ID            int     `json:"id"`
	NumberWeek    int     `json:"number_week"`
	DateRangeInit string  `json:"date_range_init"`
	DateRangeEnd  string  `json:"date_range_end"`
	ExpenseTotal  float64 `json:"expense_total"`
	IncomeTotal   float64 `json:"income_total"`
	NetTotal      float64 `json:"net_total"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
