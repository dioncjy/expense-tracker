export interface Expense {
  id?: number;
  occurredAt?: string;
  amount: number;
  category: string;
  description?: string;
}

export interface CategoryTotal {
  category: string;
  total: number;
}

export interface MonthlySummary {
  year: number;
  month: number;
  total: number;
  budgetCeiling: number;
  percentOfCeiling: number;
  nearCeiling: boolean;
  exceededCeiling: boolean;
  categoryBreakdown: CategoryTotal[];
}

