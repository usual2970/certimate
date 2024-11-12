import Show from "@/components/Show";
import { Button } from "@/components/ui/button";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Switch } from "@/components/ui/switch";
import End from "@/components/workflow/End";
import NodeRender from "@/components/workflow/NodeRender";

import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import { WorkflowNode } from "@/domain/workflow";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { ArrowLeft } from "lucide-react";
import { useEffect, useMemo } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

import { useShallow } from "zustand/shallow";

const selectState = (state: WorkflowState) => ({
  workflow: state.workflow,
  init: state.init,
  switchEnable: state.switchEnable,
  save: state.save,
});

const WorkflowDetail = () => {
  // 3. 使用正确的选择器和 shallow 比较
  const { workflow, init, switchEnable, save } = useWorkflowStore(useShallow(selectState));

  // 从 url 中获取 workflowId
  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");

  useEffect(() => {
    console.log(id);
    init(id ?? "");
  }, [id]);

  const navigate = useNavigate();

  const elements = useMemo(() => {
    let current = workflow.draft as WorkflowNode;

    const elements: JSX.Element[] = [];

    while (current) {
      // 处理普通节点
      elements.push(<NodeRender data={current} key={current.id} />);
      current = current.next as WorkflowNode;
    }

    elements.push(<End key="workflow-end" />);

    return elements;
  }, [workflow]);

  const handleBackClick = () => {
    navigate("/workflow");
  };

  const handleEnableChange = () => {
    switchEnable();
  };

  const handleWorkflowSaveClick = () => {
    save();
  };

  return (
    <>
      <WorkflowProvider>
        <ScrollArea className="h-[100vh] w-full relative bg-background">
          <div className="h-16 sticky  top-0 left-0 z-20 shadow-md bg-muted/40 flex justify-between items-center">
            <div className="px-5 text-stone-700 dark:text-stone-200 flex items-center space-x-2">
              <ArrowLeft className="cursor-pointer" onClick={handleBackClick} />
              <div className="flex flex-col space-y-2">
                <div className="">工作流</div>
                <div className="text-sm text-muted-foreground">工作流详情</div>
              </div>
            </div>
            <div className="px-5 flex items-center space-x-3">
              <Show when={!!workflow.enabled}>
                <Show when={!!workflow.hasDraft} fallback={<Button variant={"secondary"}>立即执行</Button>}>
                  <Button variant={"secondary"} onClick={handleWorkflowSaveClick}>
                    保存变更
                  </Button>
                </Show>
              </Show>

              <Switch className="dark:data-[state=unchecked]:bg-stone-400" checked={workflow.enabled ?? false} onCheckedChange={handleEnableChange} />
            </div>
          </div>

          <div className=" flex flex-col items-center mt-8">{elements}</div>

          <ScrollBar orientation="vertical" />
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </WorkflowProvider>
    </>
  );
};

export default WorkflowDetail;
