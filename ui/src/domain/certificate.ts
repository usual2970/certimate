import { Workflow } from "./workflow";

export type Certificate = {
  id: string;
  san: string;
  certificate: string;
  privateKey: string;
  issuerCertificate: string;
  certUrl: string;
  certStableUrl: string;
  output: string;
  expireAt: string;
  workflow: string;
  nodeId: string;
  created: string;
  updated: string;

  expand: {
    workflow?: Workflow;
  };
};
