import { memo } from "react";

import { type WorkflowNode, WorkflowNodeType } from "@/domain/workflow";

import WorkflowElement from "../WorkflowElement";
import BranchNode from "./BranchNode";
import ConditionNode from "./ConditionNode";
import EndNode from "./EndNode";

export type NodeRenderProps = {
  node: WorkflowNode;
  branchId?: string;
  branchIndex?: number;
  disabled?: boolean;
};

const NodeRender = ({ node: data, branchId, branchIndex, disabled }: NodeRenderProps) => {
  const render = () => {
    switch (data.type) {
      case WorkflowNodeType.Start:
      case WorkflowNodeType.Apply:
      case WorkflowNodeType.Deploy:
      case WorkflowNodeType.Notify:
        return <WorkflowElement node={data} disabled={disabled} />;
      case WorkflowNodeType.End:
        return <EndNode />;
      case WorkflowNodeType.Branch:
        return <BranchNode node={data} disabled={disabled} />;
      case WorkflowNodeType.Condition:
        return <ConditionNode node={data as WorkflowNode} branchId={branchId!} branchIndex={branchIndex!} disabled={disabled} />;
    }
  };

  return <>{render()}</>;
};

export default memo(NodeRender);
