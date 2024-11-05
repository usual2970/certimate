import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import End from "@/components/workflow/End";
import NodeRender from "@/components/workflow/NodeRender";

import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import { WorkflowNode } from "@/domain/workflow";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useMemo } from "react";

import { useShallow } from "zustand/shallow";

const selectState = (state: WorkflowState) => ({
  root: state.root,
});

const Workflow = () => {
  // 3. 使用正确的选择器和 shallow 比较
  const { root } = useWorkflowStore(useShallow(selectState));

  const elements = useMemo(() => {
    let current = root;

    const elements: JSX.Element[] = [];

    while (current) {
      // 处理普通节点
      elements.push(<NodeRender data={current} key={current.id} />);
      current = current.next as WorkflowNode;
    }

    elements.push(<End key="workflow-end" />);

    return elements;
  }, [root]);

  return (
    <>
      <WorkflowProvider>
        <ScrollArea className="h-[100vh] w-full bg-slate-50 relative">
          <div className="h-16 sticky  top-0 left-0 z-20 shadow-md bg-white"></div>

          <div className=" flex flex-col items-center mt-8">{elements}</div>

          <ScrollBar orientation="vertical" />
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </WorkflowProvider>
    </>
  );
};

export default Workflow;
