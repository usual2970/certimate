import { memo, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Flex, Typography } from "antd";
import { produce } from "immer";

import { type WorkflowNodeConfigForInspect, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import SharedNode, { type SharedNodeProps } from "./_SharedNode";
import InspectNodeConfigForm, { type InspectNodeConfigFormInstance } from "./InspectNodeConfigForm";

export type InspectNodeProps = SharedNodeProps;

const InspectNode = ({ node, disabled }: InspectNodeProps) => {
  if (node.type !== WorkflowNodeType.Inspect) {
    console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Inspect}`);
  }

  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const formRef = useRef<InspectNodeConfigFormInstance>(null);
  const [formPending, setFormPending] = useState(false);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const getFormValues = () => formRef.current!.getFieldsValue() as WorkflowNodeConfigForInspect;

  const wrappedEl = useMemo(() => {
    if (node.type !== WorkflowNodeType.Inspect) {
      console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Inspect}`);
    }

    if (!node.validated) {
      return <Typography.Link>{t("workflow_node.action.configure_node")}</Typography.Link>;
    }

    const config = (node.config as WorkflowNodeConfigForInspect) ?? {};
    return (
      <Flex className="size-full overflow-hidden" align="center" gap={8}>
        <Typography.Text className="truncate">{config.domain ?? ""}</Typography.Text>
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
        node={node}
        open={drawerOpen}
        pending={formPending}
        onConfirm={handleDrawerConfirm}
        onOpenChange={(open) => setDrawerOpen(open)}
        getFormValues={() => formRef.current!.getFieldsValue()}
      >
        <InspectNodeConfigForm ref={formRef} disabled={disabled} initialValues={node.config} />
      </SharedNode.ConfigDrawer>
    </>
  );
};

export default memo(InspectNode);
