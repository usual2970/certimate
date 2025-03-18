export interface WorkflowLogModel extends Omit<BaseModel, "updated"> {
  nodeId: string;
  nodeName: string;
  level: "DEBUG" | "INFO" | "WARN" | "ERROR";
  message: string;
  data: Record<string, any>;
}
