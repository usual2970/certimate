import { parseExpression } from "cron-parser";

export const validCronExpression = (expr: string): boolean => {
  try {
    parseExpression(expr);

    if (expr.trim().split(" ").length !== 5) return false; // pocketbase 后端仅支持五段式的表达式
    return true;
  } catch {
    return false;
  }
};

export const getNextCronExecutions = (expr: string, times = 1): Date[] => {
  if (!validCronExpression(expr)) return [];

  const now = new Date();
  const cron = parseExpression(expr, { currentDate: now, iterator: true });

  const result: Date[] = [];
  for (let i = 0; i < times; i++) {
    const next = cron.next();
    result.push(next.value.toDate());
  }
  return result;
};
