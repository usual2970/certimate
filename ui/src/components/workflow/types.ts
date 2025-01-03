import { type WorkflowBranchNode, type WorkflowNode } from "@/domain/workflow";

/**
 * @deprecated
 */
export type NodeProps = {
  node: WorkflowNode | WorkflowBranchNode;
  branchId?: string;
  branchIndex?: number;
};

/**
 * @deprecated
 */
export type BrandNodeProps = {
  node: WorkflowBranchNode;
};
