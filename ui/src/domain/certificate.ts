import { type WorkflowModel } from "./workflow";

export interface CertificateModel extends BaseModel {
  source: string;
  subjectAltNames: string;
  certificate: string;
  privateKey: string;
  effectAt: ISO8601String;
  expireAt: ISO8601String;
  workflowId: string;
  expand?: {
    workflowId?: WorkflowModel; // TODO: ugly, maybe to use an alias?
  };
}

export const CERTIFICATE_SOURCES = Object.freeze({
  WORKFLOW: "workflow",
  UPLOAD: "upload",
} as const);

export type CertificateSourceType = (typeof CERTIFICATE_SOURCES)[keyof typeof CERTIFICATE_SOURCES];
