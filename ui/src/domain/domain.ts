import { Deployment, Pahse } from "./deployment";

export type Domain = {
  id: string;
  domain: string;
  crontab: string;
  access: string;
  targetAccess: string;
  targetType: string;
  expiredAt?: string;
  phase?: Pahse;
  phaseSuccess?: boolean;
  lastDeployedAt?: string;
  enabled?: boolean;
  created?: string;
  updated?: string;
  deleted?: string;
  rightnow?: boolean;
  expand?: {
    lastDeployment?: Deployment;
  };
};

export const getLastDeployment = (domain: Domain): Deployment | undefined => {
  return domain.expand?.lastDeployment;
};
