package expenses

import "time"

// Expense represents a single spending record.
type Expense struct {
	ID          int64     `json:"id"`
	OccurredAt  time.Time `json:"occurredAt"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
}

// MonthlySummary aggregates expenses for a given month.
type MonthlySummary struct {
	Year              int     `json:"year"`
	Month             int     `json:"month"`
	Total             float64 `json:"total"`
	BudgetCeiling     float64 `json:"budgetCeiling"`
	PercentOfCeiling  float64 `json:"percentOfCeiling"`
	NearCeiling       bool    `json:"nearCeiling"`
	ExceededCeiling   bool    `json:"exceededCeiling"`
	CategoryBreakdown []CategoryTotal `json:"categoryBreakdown"`
}

// CategoryTotal represents subtotal for a category.
type CategoryTotal struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

