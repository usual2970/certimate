import { memo, useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Space, Typography } from "antd";
import { produce } from "immer";

import { deployProvidersMap } from "@/domain/provider";
import { type WorkflowNodeConfigForDeploy, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import DeployNodeConfigForm, { type DeployNodeConfigFormInstance } from "./DeployNodeConfigForm";
import SharedNode, { type SharedNodeProps } from "./_SharedNode";

export type DeployNodeProps = SharedNodeProps;

const DeployNode = ({ node, disabled }: DeployNodeProps) => {
  if (node.type !== WorkflowNodeType.Deploy) {
    console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Deploy}`);
  }

  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const formRef = useRef<DeployNodeConfigFormInstance>(null);
  const [formPending, setFormPending] = useState(false);

  const [fieldProvider, setFieldProvider] = useState<string | undefined>((node.config as WorkflowNodeConfigForDeploy)?.provider);
  useEffect(() => {
    setFieldProvider((node.config as WorkflowNodeConfigForDeploy)?.provider);
  }, [node.config?.provider]);

  const [drawerOpen, setDrawerOpen] = useState(false);
  const getFormValues = () => formRef.current!.getFieldsValue() as WorkflowNodeConfigForDeploy;

  const wrappedEl = useMemo(() => {
    if (node.type !== WorkflowNodeType.Deploy) {
      console.warn(`[certimate] current workflow node type is not: ${WorkflowNodeType.Deploy}`);
    }

    if (!node.validated) {
      return <Typography.Link>{t("workflow_node.action.configure_node")}</Typography.Link>;
    }

    const config = (node.config as WorkflowNodeConfigForDeploy) ?? {};
    const provider = deployProvidersMap.get(config.provider);
    return (
      <Space>
        <Avatar src={provider?.icon} size="small" />
        <Typography.Text className="truncate">{t(provider?.name ?? "")}</Typography.Text>
      </Space>
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

  const handleFormValuesChange = (values: Partial<WorkflowNodeConfigForDeploy>) => {
    setFieldProvider(values.provider!);
  };

  return (
    <>
      <SharedNode.Block node={node} disabled={disabled} onClick={() => setDrawerOpen(true)}>
        {wrappedEl}
      </SharedNode.Block>

      <SharedNode.ConfigDrawer
        node={node}
        footer={!!fieldProvider}
        open={drawerOpen}
        pending={formPending}
        onConfirm={handleDrawerConfirm}
        onOpenChange={(open) => setDrawerOpen(open)}
        getFormValues={() => formRef.current!.getFieldsValue()}
      >
        <DeployNodeConfigForm ref={formRef} disabled={disabled} initialValues={node.config} nodeId={node.id} onValuesChange={handleFormValuesChange} />
      </SharedNode.ConfigDrawer>
    </>
  );
};

export default memo(DeployNode);
