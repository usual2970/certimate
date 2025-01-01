import { type WorkflowBranchNode, type WorkflowNode } from "@/domain/workflow";

export type NodeProps = {
  data: WorkflowNode | WorkflowBranchNode;
  branchId?: string;
  branchIndex?: number;
};

export type BrandNodeProps = {
  data: WorkflowBranchNode;
};
