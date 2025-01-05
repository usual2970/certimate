import { memo, useMemo } from "react";

import { type WorkflowNode, WorkflowNodeType } from "@/domain/workflow";

import BranchNode from "./node/BranchNode";
import CommonNode from "./node/CommonNode";
import ConditionNode from "./node/ConditionNode";
import EndNode from "./node/EndNode";

export type WorkflowElementProps = {
  node: WorkflowNode;
  disabled?: boolean;
  branchId?: string;
  branchIndex?: number;
};

const WorkflowElement = ({ node, disabled, ...props }: WorkflowElementProps) => {
  const workflowNodeEl = useMemo(() => {
    switch (node.type) {
      case WorkflowNodeType.Start:
      case WorkflowNodeType.Apply:
      case WorkflowNodeType.Deploy:
      case WorkflowNodeType.Notify:
        return <CommonNode node={node} disabled={disabled} />;

      case WorkflowNodeType.Branch:
        return <BranchNode node={node} disabled={disabled} />;

      case WorkflowNodeType.Condition:
        return <ConditionNode node={node} disabled={disabled} branchId={props.branchId!} branchIndex={props.branchIndex!} />;

      case WorkflowNodeType.End:
        return <EndNode />;

      default:
        console.warn(`[certimate] unsupported workflow node type: ${node.type}`);
        return <></>;
    }
  }, [node, disabled, props]);

  return <>{workflowNodeEl}</>;
};

export default memo(WorkflowElement);
