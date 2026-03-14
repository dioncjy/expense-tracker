package expenses

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Handlers struct {
	svc *Service
}

func NewHandlers(svc *Service) *Handlers {
	return &Handlers{svc: svc}
}

type createExpenseRequest struct {
	OccurredAt  *time.Time `json:"occurredAt,omitempty"`
	Amount      float64    `json:"amount"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
}

func (h *Handlers) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var req createExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 || req.Category == "" {
		http.Error(w, "amount must be > 0 and category required", http.StatusBadRequest)
		return
	}

	exp := Expense{
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
	}
	if req.OccurredAt != nil {
		exp.OccurredAt = *req.OccurredAt
	}

	if err := h.svc.AddExpense(exp); err != nil {
		http.Error(w, "failed to save expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handlers) GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	year := parseIntQuery(r, "year", now.Year())
	month := parseIntQuery(r, "month", int(now.Month()))

	summary, err := h.svc.GetMonthlySummary(year, month)
	if err != nil {
		http.Error(w, "failed to load summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(summary)
}

type setCeilingRequest struct {
	Year    int     `json:"year"`
	Month   int     `json:"month"`
	Ceiling float64 `json:"ceiling"`
}

func (h *Handlers) SetMonthlyCeiling(w http.ResponseWriter, r *http.Request) {
	var req setCeilingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Year == 0 || req.Month == 0 || req.Ceiling <= 0 {
		http.Error(w, "year, month and positive ceiling are required", http.StatusBadRequest)
		return
	}

	if err := h.svc.repo.SetMonthlyCeiling(req.Year, req.Month, req.Ceiling); err != nil {
		http.Error(w, "failed to save ceiling", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func parseIntQuery(r *http.Request, key string, fallback int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}

