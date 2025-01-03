import { useMemo } from "react";

import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import EndNode from "@/components/workflow/node/EndNode";
import NodeRender from "@/components/workflow/node/NodeRender";
import { type WorkflowNode } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

export type WorkflowElementsProps = {
  className?: string;
  style?: React.CSSProperties;
};

const WorkflowElements = ({ className, style }: WorkflowElementsProps) => {
  const { workflow } = useWorkflowStore(useZustandShallowSelector(["workflow"]));

  const elements = useMemo(() => {
    const nodes: JSX.Element[] = [];

    let current = workflow.draft as WorkflowNode;
    while (current) {
      nodes.push(<NodeRender node={current} key={current.id} />);
      current = current.next as WorkflowNode;
    }

    nodes.push(<EndNode key="workflow-end" />);

    return nodes;
  }, [workflow]);

  return (
    <div className={className} style={style}>
      <div className="flex flex-col items-center overflow-auto">
        <WorkflowProvider>{elements}</WorkflowProvider>
      </div>
    </div>
  );
};

export default WorkflowElements;
