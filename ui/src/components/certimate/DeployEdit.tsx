import { createContext, useContext, type Context as ReactContext } from "react";

import { type DeployConfig } from "@/domain/domain";

export type DeployEditContext<T extends DeployConfig["config"] = DeployConfig["config"]> = {
  config: Omit<DeployConfig, "config"> & { config: T };
  setConfig: (config: Omit<DeployConfig, "config"> & { config: T }) => void;

  errors: { [K in keyof T]?: string };
  setErrors: (error: { [K in keyof T]?: string }) => void;
};

export const Context = createContext<DeployEditContext>({} as DeployEditContext);

export function useDeployEditContext<T extends DeployConfig["config"] = DeployConfig["config"]>() {
  return useContext<DeployEditContext<T>>(Context as unknown as ReactContext<DeployEditContext<T>>);
}
