export interface WorkflowRunModel extends BaseModel {
  workflow: string;
  log: WorkflowRunLog[];
  error: string;
  succeed: boolean;
}

export type WorkflowRunLog = {
  nodeName: string;
  error: string;
  outputs: WorkflowRunLogOutput[];
};

export type WorkflowRunLogOutput = {
  time: ISO8601String;
  title: string;
  content: string;
  error: string;
};
