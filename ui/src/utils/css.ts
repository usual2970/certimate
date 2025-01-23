import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export const mergeCls = (...inputs: ClassValue[]) => {
  return twMerge(clsx(inputs));
};
