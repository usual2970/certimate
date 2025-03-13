import { useMemo } from "react";

import WorkflowElement from "@/components/workflow/WorkflowElement";
import { WorkflowNodeType, newNode } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

export type WorkflowElementsProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
};

const WorkflowElements = ({ className, style, disabled }: WorkflowElementsProps) => {
  const { workflow } = useWorkflowStore(useZustandShallowSelector(["workflow"]));

  const elements = useMemo(() => {
    const root = workflow.draft;
    const nodes: JSX.Element[] = [];

    let current = root as typeof root | undefined;
    while (current) {
      nodes.push(<WorkflowElement key={current.id} node={current} disabled={disabled} />);
      current = current.next;
    }

    nodes.push(<WorkflowElement key="end" node={newNode(WorkflowNodeType.End)} />);

    return nodes;
  }, [workflow, disabled]);

  return (
    <div className={className} style={style}>
      <div className="flex w-max min-w-full flex-col items-center">{elements}</div>
    </div>
  );
};

export default WorkflowElements;
