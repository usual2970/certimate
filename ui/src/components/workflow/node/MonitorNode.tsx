import { memo, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Flex, Typography } from "antd";
import { produce } from "immer";

import { type WorkflowNodeConfigForMonitor, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import SharedNode, { type SharedNodeProps } from "./_SharedNode";
import MonitorNodeConfigForm, { type MonitorNodeConfigFormInstance } from "./MonitorNodeConfigForm";

export type MonitorNodeProps = SharedNodeProps;

const MonitorNode = ({ node, disabled }: MonitorNodeProps) => {
  if (node.type !== WorkflowNodeType.Monitor) {
    console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Monitor}`);
  }

  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const formRef = useRef<MonitorNodeConfigFormInstance>(null);
  const [formPending, setFormPending] = useState(false);
  const getFormValues = () => formRef.current!.getFieldsValue() as WorkflowNodeConfigForMonitor;

  const [drawerOpen, setDrawerOpen] = useState(false);

  const wrappedEl = useMemo(() => {
    if (node.type !== WorkflowNodeType.Monitor) {
      console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Monitor}`);
    }

    if (!node.validated) {
      return <Typography.Link>{t("workflow_node.action.configure_node")}</Typography.Link>;
    }

    const config = (node.config as WorkflowNodeConfigForMonitor) ?? {};
    return (
      <Flex className="size-full overflow-hidden" align="center" gap={8}>
        <Typography.Text className="truncate">{config.domain || config.host || ""}</Typography.Text>
      </Flex>
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
        <MonitorNodeConfigForm ref={formRef} disabled={disabled} initialValues={node.config} />
      </SharedNode.ConfigDrawer>
    </>
  );
};

export default memo(MonitorNode);
