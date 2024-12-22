import { type WorkflowModel } from "./workflow";

export interface CertificateModel extends BaseModel {
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
