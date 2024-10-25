import { createContext, useContext } from "react";

import { DeployConfig } from "@/domain/domain";

type DeployEditContext = {
  deploy: DeployConfig;
  error: Record<string, string>;
  setDeploy: (deploy: DeployConfig) => void;
  setError: (error: Record<string, string | undefined>) => void;
};

export const Context = createContext<DeployEditContext>({} as DeployEditContext);

export const useDeployEditContext = () => {
  return useContext(Context);
};
