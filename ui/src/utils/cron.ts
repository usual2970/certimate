import { CronExpressionParser } from "cron-parser";

export const validCronExpression = (expr: string): boolean => {
  try {
    CronExpressionParser.parse(expr);

    if (expr.trim().split(" ").length !== 5) return false; // pocketbase 后端仅支持五段式的表达式
    return true;
  } catch {
    return false;
  }
};

export const getNextCronExecutions = (expr: string, times = 1): Date[] => {
  if (!validCronExpression(expr)) return [];

  const now = new Date();
  const cron = CronExpressionParser.parse(expr, { currentDate: now });

  return cron.take(times).map((date) => date.toDate());
};
