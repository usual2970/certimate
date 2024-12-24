import { type WorkflowOutput } from "./workflow";

export interface WorkflowRunModel extends BaseModel {
  workflow: string;
  log: WorkflowRunLog[];
  error: string;
  succeed: boolean;
}

export type WorkflowRunLog = {
  nodeName: string;
  error: string;
  outputs: WorkflowOutput[];
};
