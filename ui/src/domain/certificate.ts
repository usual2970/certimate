import { type BaseModel } from "pocketbase";

import { WorkflowModel } from "./workflow";

export interface CertificateModel extends Omit<BaseModel, "created" | "updated"> {
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
  expand: {
    workflow?: WorkflowModel;
  };
}
