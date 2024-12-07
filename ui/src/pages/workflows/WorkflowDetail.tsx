import { run } from "@/api/workflow";
import Show from "@/components/Show";
import { Button } from "@/components/ui/button";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Switch } from "@/components/ui/switch";
import { Toaster } from "@/components/ui/toaster";
import { useToast } from "@/components/ui/use-toast";
import End from "@/components/workflow/End";
import NodeRender from "@/components/workflow/NodeRender";
import WorkflowBaseInfoEditDialog from "@/components/workflow/WorkflowBaseInfoEditDialog";
import WorkflowLog from "@/components/workflow/WorkflowLog";

import { cn } from "@/components/ui/utils";
import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import { allNodesValidated, WorkflowNode } from "@/domain/workflow";
import { getErrMessage } from "@/lib/error";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { ArrowLeft } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
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
  const [locId, setLocId] = useState<string>("");
  const id = searchParams.get("id");

  const [tab, setTab] = useState("workflow");

  const [running, setRunning] = useState(false);

  useEffect(() => {
    init(id ?? "");
    if (id) {
      setLocId(id);
    }
  }, [id]);

  const navigate = useNavigate();

  const { toast } = useToast();

  const { t } = useTranslation();

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
    // 返回上一步
    navigate(-1);
  };

  const handleEnableChange = () => {
    if (!workflow.enabled && !allNodesValidated(workflow.draft as WorkflowNode)) {
      toast({
        title: t("workflow.detail.action.save.failed"),
        description: t("workflow.detail.action.save.failed.uncompleted"),
        variant: "destructive",
      });
      return;
    }
    switchEnable();
    if (!locId) {
      navigate(`/workflows/detail?id=${workflow.id}`);
    }
  };

  const handleWorkflowSaveClick = () => {
    if (!allNodesValidated(workflow.draft as WorkflowNode)) {
      toast({
        title: t("workflow.detail.action.save.failed"),
        description: t("workflow.detail.action.save.failed.uncompleted"),
        variant: "destructive",
      });
      return;
    }
    save();
    if (!locId) {
      navigate(`/workflows/detail?id=${workflow.id}`);
    }
  };

  const getTabCls = (tabName: string) => {
    if (tab === tabName) {
      return "text-primary border-primary";
    }
    return "border-transparent hover:text-primary hover:border-b-primary";
  };

  const handleRunClick = async () => {
    if (running) {
      return;
    }
    setRunning(true);
    try {
      await run(workflow.id as string);
      toast({
        title: t("workflow.detail.action.run.success"),
        description: t("workflow.detail.action.run.success"),
        variant: "default",
      });
    } catch (e) {
      toast({
        title: t("workflow.detail.action.run.failed"),
        description: getErrMessage(e),
        variant: "destructive",
      });
    }
    setRunning(false);
  };

  return (
    <>
      <WorkflowProvider>
        <Toaster />
        <ScrollArea className="h-[100vh] w-full relative bg-background">
          <div className="h-16 sticky  top-0 left-0 z-20 shadow-md bg-muted/40 flex justify-between items-center">
            <div className="px-5 text-stone-700 dark:text-stone-200 flex items-center space-x-2">
              <ArrowLeft className="cursor-pointer" onClick={handleBackClick} />
              <WorkflowBaseInfoEditDialog
                trigger={
                  <div className="flex flex-col space-y-1 cursor-pointer items-start">
                    <div className="truncate  max-w-[200px]">{workflow.name ? workflow.name : t("workflow.props.name.default")}</div>
                    <div className="text-sm text-muted-foreground truncate  max-w-[200px]">
                      {workflow.description ? workflow.description : t("workflow.props.description.placeholder")}
                    </div>
                  </div>
                }
              />
            </div>

            <div className="flex justify-between space-x-5 text-muted-foreground text-lg h-full">
              <div
                className={cn("h-full flex items-center cursor-pointer border-b-2", getTabCls("workflow"))}
                onClick={() => {
                  setTab("workflow");
                }}
              >
                <div>{t("workflow.detail.title")}</div>
              </div>
              <div
                className={cn("h-full flex items-center cursor-pointer border-b-2", getTabCls("history"))}
                onClick={() => {
                  setTab("history");
                }}
              >
                <div>{t("workflow.detail.history")}</div>
              </div>
            </div>

            <div className="px-5 flex items-center space-x-3">
              <Show when={!!workflow.enabled}>
                <Show
                  when={!!workflow.hasDraft}
                  fallback={
                    <Button variant={"secondary"} onClick={handleRunClick}>
                      {running ? t("workflow.detail.action.running") : t("workflow.detail.action.run")}
                    </Button>
                  }
                >
                  <Button variant={"secondary"} onClick={handleWorkflowSaveClick}>
                    {t("workflow.detail.action.save")}
                  </Button>
                </Show>
              </Show>

              <Switch className="dark:data-[state=unchecked]:bg-stone-400" checked={workflow.enabled ?? false} onCheckedChange={handleEnableChange} />
            </div>
          </div>
          <Show when={tab === "workflow"}>
            <div className=" flex flex-col items-center mt-8">{elements}</div>
          </Show>

          <Show when={!!locId && tab === "history"}>
            <div className=" flex flex-col items-center mt-8">
              <WorkflowLog />
            </div>
          </Show>

          <ScrollBar orientation="vertical" />
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </WorkflowProvider>
    </>
  );
};

export default WorkflowDetail;
