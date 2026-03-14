package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"vibe-code/backend/internal/expenses"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		occurred_at TEXT NOT NULL,
		amount REAL NOT NULL,
		category TEXT NOT NULL,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS monthly_ceilings (
		year INTEGER NOT NULL,
		month INTEGER NOT NULL,
		ceiling REAL NOT NULL,
		PRIMARY KEY (year, month)
	);
	`
	_, err := db.Exec(schema)
	return err
}

func (r *SQLiteRepository) InsertExpense(e expenses.Expense) error {
	const q = `
		INSERT INTO expenses (occurred_at, amount, category, description)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(q, e.OccurredAt.Format(time.RFC3339), e.Amount, e.Category, e.Description)
	return err
}

func (r *SQLiteRepository) ListExpensesByMonth(year int, month int) ([]expenses.Expense, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	const q = `
		SELECT id, occurred_at, amount, category, description
		FROM expenses
		WHERE occurred_at >= ? AND occurred_at < ?
		ORDER BY occurred_at ASC
	`

	rows, err := r.db.Query(q, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []expenses.Expense
	for rows.Next() {
		var (
			id          int64
			occurredRaw string
			amount      float64
			category    string
			description sql.NullString
		)

		if err := rows.Scan(&id, &occurredRaw, &amount, &category, &description); err != nil {
			return nil, err
		}

		t, err := time.Parse(time.RFC3339, occurredRaw)
		if err != nil {
			return nil, err
		}

		out = append(out, expenses.Expense{
			ID:          id,
			OccurredAt:  t,
			Amount:      amount,
			Category:    category,
			Description: description.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *SQLiteRepository) GetMonthlyCeiling(year int, month int) (float64, error) {
	const q = `
		SELECT ceiling FROM monthly_ceilings
		WHERE year = ? AND month = ?
	`
	var ceiling float64
	err := r.db.QueryRow(q, year, month).Scan(&ceiling)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return ceiling, err
}

func (r *SQLiteRepository) SetMonthlyCeiling(year int, month int, ceiling float64) error {
	const q = `
		INSERT INTO monthly_ceilings (year, month, ceiling)
		VALUES (?, ?, ?)
		ON CONFLICT(year, month) DO UPDATE SET ceiling = excluded.ceiling
	`
	_, err := r.db.Exec(q, year, month, ceiling)
	return err
}

