import { Domain } from "./domain";

export type Deployment = {
  id: string;
  domain: string;
  log: {
    apply?: Log[];
    check?: Log[];
    deploy?: Log[];
  };
  phase: Pahse;
  phaseSuccess: boolean;
  deployedAt: string;
  created: string;
  updated: string;
  expand: {
    domain?: Domain;
  };
};

export type Pahse = "apply" | "check" | "deploy";

export type Log = {
  time: string;
  message: string;
  error: string;
};

export type DeploymentListReq = {
  domain?: string;
  page?: number;
  perPage?: number;
};
