# Expense Tracker – Go + React

A personal **expense tracker** built to help you log daily/weekly spending, see **monthly and yearly trends**, and stay within a **monthly budget ceiling**.

## Tech stack

- **Backend**: Go (HTTP server + SQLite)
  - Stores expenses and monthly budget ceilings in a local SQLite database file.
  - Exposes JSON APIs for:
    - Creating expenses.
    - Getting monthly summaries (totals, category breakdowns, near/over ceiling).
    - Setting the monthly ceiling.
- **Frontend**: React + TypeScript (Vite)
  - Dashboard UI to:
    - Add new expenses.
    - View current month totals and category breakdown.
    - See warnings when you're close to or over your monthly ceiling.

## Running the backend

From the `backend` folder:

```bash
cd backend
go run ./cmd/api
```

- Server listens on `http://localhost:8080`.
- Uses `expense-tracker.db` in the project root by default.
- You can override the DB path with:

```bash
EXPENSE_DB_PATH=/path/to/your.db go run ./cmd/api
```

## Running the frontend

From the `frontend` folder:

```bash
cd frontend
npm install
npm run dev
```

This will start the React app (by default on `http://localhost:5173`).

## Project layout

- `backend/` – Go API + SQLite storage.
- `frontend/` – React + TypeScript UI (Vite).
