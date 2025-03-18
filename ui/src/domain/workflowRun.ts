import { type WorkflowModel } from "./workflow";

export interface WorkflowRunModel extends BaseModel {
  workflowId: string;
  status: string;
  trigger: string;
  startedAt: ISO8601String;
  endedAt: ISO8601String;
  error?: string;
  expand?: {
    workflowId?: WorkflowModel; // TODO: ugly, maybe to use an alias?
  };
}

export const WORKFLOW_RUN_STATUSES = Object.freeze({
  PENDING: "pending",
  RUNNING: "running",
  SUCCEEDED: "succeeded",
  FAILED: "failed",
  CANCELED: "canceled",
} as const);

export type WorkflorRunStatusType = (typeof WORKFLOW_RUN_STATUSES)[keyof typeof WORKFLOW_RUN_STATUSES];
