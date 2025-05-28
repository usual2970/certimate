import { memo, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Typography } from "antd";
import { produce } from "immer";

import { WORKFLOW_TRIGGERS, type WorkflowNodeConfigForStart, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import SharedNode, { type SharedNodeProps } from "./_SharedNode";
import StartNodeConfigForm, { type StartNodeConfigFormInstance } from "./StartNodeConfigForm";

export type StartNodeProps = SharedNodeProps;

const StartNode = ({ node, disabled }: StartNodeProps) => {
  if (node.type !== WorkflowNodeType.Start) {
    console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Start}`);
  }

  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const formRef = useRef<StartNodeConfigFormInstance>(null);
  const [formPending, setFormPending] = useState(false);
  const getFormValues = () => formRef.current!.getFieldsValue() as WorkflowNodeConfigForStart;

  const [drawerOpen, setDrawerOpen] = useState(false);

  const wrappedEl = useMemo(() => {
    if (node.type !== WorkflowNodeType.Start) {
      console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Start}`);
    }

    if (!node.validated) {
      return <Typography.Link>{t("workflow_node.action.configure_node")}</Typography.Link>;
    }

    const config = (node.config as WorkflowNodeConfigForStart) ?? {};
    return (
      <div className="flex items-center justify-between space-x-2">
        <Typography.Text className="truncate">
          {config.trigger === WORKFLOW_TRIGGERS.AUTO
            ? t("workflow.props.trigger.auto")
            : config.trigger === WORKFLOW_TRIGGERS.MANUAL
              ? t("workflow.props.trigger.manual")
              : "　"}
        </Typography.Text>
        <Typography.Text className="truncate" type="secondary">
          {config.trigger === WORKFLOW_TRIGGERS.AUTO ? config.triggerCron : ""}
        </Typography.Text>
      </div>
    );
  }, [node]);

  const handleDrawerConfirm = async () => {
    setFormPending(true);
    try {
      await formRef.current!.validateFields();
    } catch (err) {
      setFormPending(false);
      throw err;
    }

    try {
      const newValues = getFormValues();
      const newNode = produce(node, (draft) => {
        draft.config = {
          ...newValues,
        };
        draft.validated = true;
      });
      await updateNode(newNode);
    } finally {
      setFormPending(false);
    }
  };

  return (
    <>
      <SharedNode.Block node={node} disabled={disabled} onClick={() => setDrawerOpen(true)}>
        {wrappedEl}
      </SharedNode.Block>

      <SharedNode.ConfigDrawer
        getConfigNewValues={getFormValues}
        node={node}
        open={drawerOpen}
        pending={formPending}
        onConfirm={handleDrawerConfirm}
        onOpenChange={(open) => setDrawerOpen(open)}
      >
        <StartNodeConfigForm ref={formRef} disabled={disabled} initialValues={node.config} />
      </SharedNode.ConfigDrawer>
    </>
  );
};

export default memo(StartNode);
