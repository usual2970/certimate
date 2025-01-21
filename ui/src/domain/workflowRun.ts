import type { WorkflowModel } from "./workflow";

export interface WorkflowRunModel extends BaseModel {
  workflowId: string;
  status: string;
  trigger: string;
  startedAt: ISO8601String;
  endedAt: ISO8601String;
  logs?: WorkflowRunLog[];
  error?: string;
  expand?: {
    workflowId?: WorkflowModel;
  };
}

export type WorkflowRunLog = {
  nodeId: string;
  nodeName: string;
  outputs?: WorkflowRunLogOutput[];
  error?: string;
};

export type WorkflowRunLogOutput = {
  time: ISO8601String;
  title: string;
  content: string;
  error?: string;
};

export const WORKFLOW_RUN_STATUSES = Object.freeze({
  PENDING: "pending",
  RUNNING: "running",
  SUCCEEDED: "succeeded",
  FAILED: "failed",
  CANCELED: "canceled",
} as const);

export type WorkflorRunStatusType = (typeof WORKFLOW_RUN_STATUSES)[keyof typeof WORKFLOW_RUN_STATUSES];
