package expenses

import "time"

// Repository defines the storage interface for expenses and settings.
type Repository interface {
	InsertExpense(e Expense) error
	ListExpensesByMonth(year int, month int) ([]Expense, error)
	GetMonthlyCeiling(year int, month int) (float64, error)
	SetMonthlyCeiling(year int, month int, ceiling float64) error
}

type Service struct {
	repo          Repository
	nearThreshold float64
}

// NewService creates a new Service with the given repository.
// nearThreshold is the fraction of the ceiling at which we consider the user "near" their limit (e.g. 0.8 for 80%).
func NewService(repo Repository, nearThreshold float64) *Service {
	if nearThreshold <= 0 || nearThreshold > 1 {
		nearThreshold = 0.8
	}
	return &Service{repo: repo, nearThreshold: nearThreshold}
}

func (s *Service) AddExpense(e Expense) error {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	return s.repo.InsertExpense(e)
}

// GetMonthlySummary computes totals and ceiling status for the given month.
func (s *Service) GetMonthlySummary(year int, month int) (MonthlySummary, error) {
	expenses, err := s.repo.ListExpensesByMonth(year, month)
	if err != nil {
		return MonthlySummary{}, err
	}

	var total float64
	categoryTotals := make(map[string]float64)

	for _, e := range expenses {
		total += e.Amount
		categoryTotals[e.Category] += e.Amount
	}

	var breakdown []CategoryTotal
	for cat, t := range categoryTotals {
		breakdown = append(breakdown, CategoryTotal{
			Category: cat,
			Total:    t,
		})
	}

	ceiling, err := s.repo.GetMonthlyCeiling(year, month)
	if err != nil {
		// Treat missing ceiling as zero; caller can decide how to handle.
		ceiling = 0
	}

	summary := MonthlySummary{
		Year:              year,
		Month:             month,
		Total:             total,
		BudgetCeiling:     ceiling,
		CategoryBreakdown: breakdown,
	}

	if ceiling > 0 {
		summary.PercentOfCeiling = total / ceiling
		summary.NearCeiling = summary.PercentOfCeiling >= s.nearThreshold && total <= ceiling
		summary.ExceededCeiling = total > ceiling
	}

	return summary, nil
}

