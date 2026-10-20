package dtos

import "encoding/json"

// CategoryDTO is the write DTO for creating/updating a category.
type CategoryDTO struct {
	Type string `json:"type"`
}

// ExpenseDTO is the write DTO for creating/updating an expense.
type ExpenseDTO struct {
	Date       string          `json:"date"`
	Concept    string          `json:"concept"`
	CategoryID int             `json:"category_id"`
	Amount     float64         `json:"amount"`
	Quantity   float64         `json:"quantity"`
	Total      float64         `json:"total"`
	ReceiptID  int             `json:"receipt_id"`
	Meta       json.RawMessage `json:"meta,omitempty"`
}

// IncomeDTO is the write DTO for creating/updating an income record.
type IncomeDTO struct {
	Date       string          `json:"date"`
	Concept    string          `json:"concept"`
	CategoryID int             `json:"category_id"`
	Amount     float64         `json:"amount"`
	Quantity   float64         `json:"quantity"`
	Total      float64         `json:"total"`
	ReceiptID  int             `json:"receipt_id"`
	Meta       json.RawMessage `json:"meta,omitempty"`
}

// FuelExpenseDTO is the write DTO for creating/updating a fuel expense.
type FuelExpenseDTO struct {
	ExpenseID int `json:"id_expense"`
	DriverID  int `json:"id_driver"`
}

// StatsDTO is the write DTO for creating/updating a weekly stats summary.
type StatsDTO struct {
	NumberWeek    int     `json:"number_week"`
	DateRangeInit string  `json:"date_range_init"`
	DateRangeEnd  string  `json:"date_range_end"`
	ExpenseTotal  float64 `json:"expense_total"`
	IncomeTotal   float64 `json:"income_total"`
	NetTotal      float64 `json:"net_total"`
}
