import { parseExpression } from "cron-parser";

export const validCronExpression = (expr: string): boolean => {
  try {
    parseExpression(expr);
    return true;
  } catch {
    return false;
  }
};

export const getNextCronExecutions = (expr: string, times = 1): Date[] => {
  if (!expr) return [];

  try {
    const now = new Date();
    const cron = parseExpression(expr, { currentDate: now, iterator: true });

    const result: Date[] = [];
    for (let i = 0; i < times; i++) {
      const next = cron.next();
      result.push(next.value.toDate());
    }
    return result;
  } catch {
    return [];
  }
};
