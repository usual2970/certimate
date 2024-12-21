import { useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, message, notification, Switch } from "antd";
import { useShallow } from "zustand/shallow";
import { ArrowLeft as ArrowLeftIcon } from "lucide-react";

import Show from "@/components/Show";
import End from "@/components/workflow/End";
import NodeRender from "@/components/workflow/NodeRender";
import WorkflowBaseInfoEditDialog from "@/components/workflow/WorkflowBaseInfoEditDialog";
import WorkflowLog from "@/components/workflow/WorkflowLog";
import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import { cn } from "@/components/ui/utils";
import { allNodesValidated, WorkflowNode } from "@/domain/workflow";
import { useWorkflowStore, WorkflowState } from "@/stores/workflow";
import { run as runWorkflow } from "@/api/workflow";

const selectState = (state: WorkflowState) => ({
  workflow: state.workflow,
  init: state.init,
  switchEnable: state.switchEnable,
  save: state.save,
});

const WorkflowDetail = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [_, NotificationContextHolder] = notification.useNotification();

  // 3. 使用正确的选择器和 shallow 比较
  const { workflow, init, switchEnable, save } = useWorkflowStore(useShallow(selectState));

  // 从 url 中获取 workflowId
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
      messageApi.warning(t("workflow.detail.action.save.failed.uncompleted"));
      return;
    }
    switchEnable();
    if (!locId) {
      navigate(`/workflows/detail?id=${workflow.id}`);
    }
  };

  const handleWorkflowSaveClick = () => {
    if (!allNodesValidated(workflow.draft as WorkflowNode)) {
      messageApi.warning(t("workflow.detail.action.save.failed.uncompleted"));
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
      await runWorkflow(workflow.id as string);
      messageApi.success(t("workflow.detail.action.run.success"));
    } catch (e) {
      messageApi.warning(t("workflow.detail.action.run.failed"));
    }
    setRunning(false);
  };

  return (
    <div>
      {MessageContextHolder}
      {NotificationContextHolder}

      <WorkflowProvider>
        <div className="h-16 sticky  top-0 left-0 z-[1]` shadow-md bg-muted/40 flex justify-between items-center">
          <div className="px-5 text-stone-700 dark:text-stone-200 flex items-center space-x-2">
            <ArrowLeftIcon className="cursor-pointer" onClick={handleBackClick} />
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
                  <Button type="text" onClick={handleRunClick}>
                    {running ? t("workflow.detail.action.running") : t("workflow.detail.action.run")}
                  </Button>
                }
              >
                <Button type="primary" onClick={handleWorkflowSaveClick}>
                  {t("workflow.detail.action.save")}
                </Button>
              </Show>
            </Show>

            <Switch checked={workflow.enabled ?? false} onChange={handleEnableChange} />
          </div>
        </div>

        <Show when={tab === "workflow"}>
          <div className="p-4">
            <div className="flex flex-col items-center p-4 bg-background">{elements}</div>
          </div>
        </Show>

        <Show when={!!locId && tab === "history"}>
          <div className="p-4">
            <div className="flex flex-col items-center p-4 bg-background">
              <WorkflowLog />
            </div>
          </div>
        </Show>
      </WorkflowProvider>
    </div>
  );
};

export default WorkflowDetail;
