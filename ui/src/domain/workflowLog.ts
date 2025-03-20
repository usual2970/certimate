export interface WorkflowLogModel extends Omit<BaseModel, "updated"> {
  nodeId: string;
  nodeName: string;
  timestamp: ReturnType<typeof Date.prototype.getTime>;
  level: "DEBUG" | "INFO" | "WARN" | "ERROR";
  message: string;
  data: Record<string, any>;
}
