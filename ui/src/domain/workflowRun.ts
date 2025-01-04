export interface WorkflowRunModel extends BaseModel {
  workflowId: string;
  logs: WorkflowRunLog[];
  error: string;
  succeeded: boolean;
}

export type WorkflowRunLog = {
  nodeId: string;
  nodeName: string;
  outputs: WorkflowRunLogOutput[];
  error: string;
};

export type WorkflowRunLogOutput = {
  time: ISO8601String;
  title: string;
  content: string;
  error: string;
};
