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
};

const NodeRender = ({ node: data, branchId, branchIndex }: NodeRenderProps) => {
  const render = () => {
    switch (data.type) {
      case WorkflowNodeType.Start:
      case WorkflowNodeType.Apply:
      case WorkflowNodeType.Deploy:
      case WorkflowNodeType.Notify:
        return <WorkflowElement node={data} />;
      case WorkflowNodeType.End:
        return <EndNode />;
      case WorkflowNodeType.Branch:
        return <BranchNode node={data} />;
      case WorkflowNodeType.Condition:
        return <ConditionNode node={data as WorkflowNode} branchId={branchId} branchIndex={branchIndex} />;
    }
  };

  return <>{render()}</>;
};

export default memo(NodeRender);
