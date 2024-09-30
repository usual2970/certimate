import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function wrapperEnv(envConf: Recordable): ViteEnv {
  const ret: any = {};

  for (const key of Object.keys(envConf)) {
    let value = envConf[key].replace(/\\n/g, "\n");
    value = value === "true" ? true : value === "false" ? false : value;

    if (key === "VITE_PORT") {
      value = Number(value);
    }
    if (key === "VITE_PROXY") {
      try {
        value = JSON.parse(value);
      } catch (err) {
        console.log(err);
      }
    }
    ret[key] = value;
    process.env[key] = value;
  }
  return ret;
}
