package main

import (
	"log"
	"net/http"
	"os"

	"vibe-code/backend/internal/expenses"
	"vibe-code/backend/internal/storage"
)

func main() {
	dbPath := os.Getenv("EXPENSE_DB_PATH")
	if dbPath == "" {
		dbPath = "expense-tracker.db"
	}

	repo, err := storage.NewSQLiteRepository(dbPath)
	if err != nil {
		log.Fatalf("failed to init sqlite: %v", err)
	}

	svc := expenses.NewService(repo, 0.8)
	h := expenses.NewHandlers(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /api/expenses", h.CreateExpense)
	mux.HandleFunc("GET /api/summary/monthly", h.GetMonthlySummary)
	mux.HandleFunc("PUT /api/settings/ceiling", h.SetMonthlyCeiling)

	addr := ":8080"
	log.Printf("expense tracker API listening on %s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

