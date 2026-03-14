# Backend (Go) – Expense Tracker

This folder contains the Go backend for the expense tracker.

## Layout

- `cmd/api/main.go`: HTTP server entrypoint (currently exposes `/healthz`).
- `internal/expenses`: Core domain for expenses and summaries.
- `go.mod`: Go module definition.

## Next steps

- Add HTTP handlers for:
  - Creating an expense.
  - Listing expenses by month/week.
  - Getting monthly/yearly summaries.
  - Getting/setting the monthly ceiling.
- Implement a storage layer (e.g. SQLite) that satisfies the `Repository` interface in `internal/expenses/service.go`.
- Wire the service into the HTTP handlers and expose JSON APIs for the React frontend.

