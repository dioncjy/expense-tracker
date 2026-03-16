import { useEffect, useMemo, useState } from 'react';
import './App.css';
import type { MonthlySummary } from './types/api';
import { createExpense, getMonthlySummary, setMonthlyCeiling } from './api/expenses';

type Status = 'idle' | 'loading' | 'error' | 'success';

interface ExpenseFormState {
  amount: string;
  category: string;
  description: string;
}

function getCurrentYearMonth() {
  const now = new Date();
  return { year: now.getFullYear(), month: now.getMonth() + 1 };
}

function formatCurrency(amount: number): string {
  if (!Number.isFinite(amount)) return '-';
  return amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function App() {
  const [{ year, month }] = useState(getCurrentYearMonth);
  const [summary, setSummary] = useState<MonthlySummary | null>(null);
  const [summaryStatus, setSummaryStatus] = useState<Status>('idle');
  const [error, setError] = useState<string | null>(null);

  const [form, setForm] = useState<ExpenseFormState>({
    amount: '',
    category: '',
    description: '',
  });
  const [formStatus, setFormStatus] = useState<Status>('idle');

  const [ceilingInput, setCeilingInput] = useState('');
  const [ceilingStatus, setCeilingStatus] = useState<Status>('idle');

  const nearOrOverText = useMemo(() => {
    if (!summary) return '';
    if (summary.exceededCeiling) return 'You have exceeded your monthly ceiling.';
    if (summary.nearCeiling) return 'You are close to your monthly ceiling.';
    return '';
  }, [summary]);

  async function loadSummary() {
    try {
      setSummaryStatus('loading');
      setError(null);
      const data = await getMonthlySummary(year, month);
      setSummary(data);
      setSummaryStatus('success');
      if (data.budgetCeiling && !ceilingInput) {
        setCeilingInput(String(data.budgetCeiling));
      }
    } catch (e) {
      setSummaryStatus('error');
      setError(e instanceof Error ? e.message : 'Failed to load summary');
    }
  }

  useEffect(() => {
    void loadSummary();
  }, []);

  async function handleSubmitExpense(e: React.FormEvent) {
    e.preventDefault();
    const amount = Number(form.amount);
    if (!amount || amount <= 0 || !form.category.trim()) {
      setError('Please enter a positive amount and a category.');
      return;
    }
    try {
      setFormStatus('loading');
      setError(null);
      await createExpense({
        amount,
        category: form.category.trim(),
        description: form.description.trim() || undefined,
      });
      setFormStatus('success');
      setForm({ amount: '', category: '', description: '' });
      await loadSummary();
    } catch (e) {
      setFormStatus('error');
      setError(e instanceof Error ? e.message : 'Failed to add expense');
    }
  }

  async function handleSaveCeiling(e: React.FormEvent) {
    e.preventDefault();
    const value = Number(ceilingInput);
    if (!value || value <= 0) {
      setError('Please enter a positive ceiling.');
      return;
    }
    try {
      setCeilingStatus('loading');
      setError(null);
      await setMonthlyCeiling(year, month, value);
      setCeilingStatus('success');
      await loadSummary();
    } catch (e) {
      setCeilingStatus('error');
      setError(e instanceof Error ? e.message : 'Failed to save ceiling');
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <div>
          <h1>Expense Tracker</h1>
          <p>
            {year} – {month.toString().padStart(2, '0')}
          </p>
        </div>
      </header>

      {error && <div className="alert alert-error">{error}</div>}
      {nearOrOverText && !error && (
        <div className="alert alert-warning">
          <strong>Heads up:</strong> {nearOrOverText}
        </div>
      )}

      <main className="grid">
        <section className="card">
          <h2>Monthly summary</h2>
          {summaryStatus === 'loading' && <p>Loading summary…</p>}
          {summary && summaryStatus !== 'loading' && (
            <>
              <div className="summary-row">
                <span>Total spent</span>
                <strong>${formatCurrency(summary.total)}</strong>
              </div>
              <div className="summary-row">
                <span>Ceiling</span>
                <strong>
                  {summary.budgetCeiling
                    ? `$${formatCurrency(summary.budgetCeiling)}`
                    : 'Not set'}
                </strong>
              </div>
              {summary.budgetCeiling > 0 && (
                <div className="summary-row">
                  <span>Used</span>
                  <strong>{Math.round(summary.percentOfCeiling * 100)}%</strong>
                </div>
              )}

              {summary.categoryBreakdown.length > 0 && (
                <>
                  <h3>By category</h3>
                  <ul className="category-list">
                    {summary.categoryBreakdown.map((c) => (
                      <li key={c.category}>
                        <span>{c.category}</span>
                        <span>${formatCurrency(c.total)}</span>
                      </li>
                    ))}
                  </ul>
                </>
              )}
            </>
          )}
        </section>

        <section className="card">
          <h2>Add expense</h2>
          <form className="form" onSubmit={handleSubmitExpense}>
            <label>
              <span>Amount</span>
              <input
                type="number"
                step="0.01"
                min="0"
                value={form.amount}
                onChange={(e) => setForm((f) => ({ ...f, amount: e.target.value }))}
                required
              />
            </label>
            <label>
              <span>Category</span>
              <input
                type="text"
                value={form.category}
                onChange={(e) => setForm((f) => ({ ...f, category: e.target.value }))}
                placeholder="Food, Transport, etc."
                required
              />
            </label>
            <label>
              <span>Description (optional)</span>
              <input
                type="text"
                value={form.description}
                onChange={(e) =>
                  setForm((f) => ({ ...f, description: e.target.value }))
                }
                placeholder="Lunch, taxi, groceries…"
              />
            </label>
            <button type="submit" disabled={formStatus === 'loading'}>
              {formStatus === 'loading' ? 'Saving…' : 'Add expense'}
            </button>
          </form>
        </section>

        <section className="card">
          <h2>Monthly ceiling</h2>
          <form className="form-inline" onSubmit={handleSaveCeiling}>
            <input
              type="number"
              step="0.01"
              min="0"
              value={ceilingInput}
              onChange={(e) => setCeilingInput(e.target.value)}
              placeholder="Enter monthly ceiling"
            />
            <button type="submit" disabled={ceilingStatus === 'loading'}>
              {ceilingStatus === 'loading' ? 'Saving…' : 'Save'}
            </button>
          </form>
        </section>
      </main>
    </div>
  );
}

export default App;
