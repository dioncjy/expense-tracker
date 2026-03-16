import { request } from './client';
import type { Expense, MonthlySummary } from '../types/api';

export async function createExpense(payload: Omit<Expense, 'id'>): Promise<void> {
  await request<{ status: string }>('/api/expenses', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function getMonthlySummary(
  year: number,
  month: number,
): Promise<MonthlySummary> {
  const params = new URLSearchParams({ year: String(year), month: String(month) });
  const data = await request<MonthlySummary>(
    `/api/summary/monthly?${params.toString()}`,
  );
  // Backend may return `null` for empty slice; normalize to empty array.
  return {
    ...data,
    categoryBreakdown: data.categoryBreakdown ?? [],
  };
}

export async function setMonthlyCeiling(
  year: number,
  month: number,
  ceiling: number,
): Promise<void> {
  await request<{ status: string }>('/api/settings/ceiling', {
    method: 'PUT',
    body: JSON.stringify({ year, month, ceiling }),
  });
}

