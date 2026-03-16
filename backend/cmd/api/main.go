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
		addCORSHeaders(w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/api/expenses", func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w, r)
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodPost:
			h.CreateExpense(w, r)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/api/summary/monthly", func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w, r)
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodGet:
			h.GetMonthlySummary(w, r)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/api/settings/ceiling", func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w, r)
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodPut:
			h.SetMonthlyCeiling(w, r)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	addr := ":8080"
	log.Printf("expense tracker API listening on %s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

