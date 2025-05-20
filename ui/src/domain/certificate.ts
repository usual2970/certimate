import { type WorkflowModel } from "./workflow";

export interface CertificateModel extends BaseModel {
  source: string;
  subjectAltNames: string;
  serialNumber: string;
  certificate: string;
  privateKey: string;
  issuerOrg: string;
  keyAlgorithm: string;
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

export const CERTIFICATE_FORMATS = Object.freeze({
  PEM: "PEM",
  PFX: "PFX",
  JKS: "JKS",
} as const);

export type CertificateFormatType = (typeof CERTIFICATE_FORMATS)[keyof typeof CERTIFICATE_FORMATS];
