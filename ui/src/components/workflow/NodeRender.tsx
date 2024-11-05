import { memo } from "react";
import { WorkflowBranchNode, WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import Node from "./Node";
import End from "./End";
import BranchNode from "./BranchNode";
import ConditionNode from "./ConditionNode";
import { NodeProps } from "./types";

const NodeRender = memo(({ data, branchId, branchIndex }: NodeProps) => {
  const render = () => {
    switch (data.type) {
      case WorkflowNodeType.Start:
      case WorkflowNodeType.Apply:
      case WorkflowNodeType.Deploy:
      case WorkflowNodeType.Notify:
        return <Node data={data} />;
      case WorkflowNodeType.End:
        return <End />;
      case WorkflowNodeType.Branch:
        return <BranchNode data={data as WorkflowBranchNode} />;
      case WorkflowNodeType.Condition:
        return <ConditionNode data={data as WorkflowNode} branchId={branchId} branchIndex={branchIndex} />;
    }
  };

  return <>{render()}</>;
});

export default NodeRender;
