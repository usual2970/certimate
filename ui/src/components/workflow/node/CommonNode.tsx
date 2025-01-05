import { memo, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { CloseCircleOutlined as CloseCircleOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { Avatar, Button, Card, Dropdown, Popover, Space, Typography } from "antd";
import { produce } from "immer";

import Show from "@/components/Show";
import { deployProvidersMap } from "@/domain/provider";
import { notifyChannelsMap } from "@/domain/settings";
import {
  WORKFLOW_TRIGGERS,
  type WorkflowNode,
  type WorkflowNodeConfigForApply,
  type WorkflowNodeConfigForDeploy,
  type WorkflowNodeConfigForNotify,
  type WorkflowNodeConfigForStart,
  WorkflowNodeType,
} from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { useContactEmailsStore } from "@/stores/contact";
import { useWorkflowStore } from "@/stores/workflow";

import AddNode from "./AddNode";
import ApplyNodeForm from "./ApplyNodeForm";
import DeployNodeForm from "./DeployNodeForm";
import NotifyNodeForm from "./NotifyNodeForm";
import StartNodeForm from "./StartNodeForm";
import { usePanelContext } from "../panel/PanelContext";

export type CommonNodeProps = {
  node: WorkflowNode;
  disabled?: boolean;
};

const CommonNode = ({ node, disabled }: CommonNodeProps) => {
  const { t } = useTranslation();

  const { accesses } = useAccessesStore(useZustandShallowSelector("accesses"));
  const { addEmail } = useContactEmailsStore(useZustandShallowSelector(["addEmail"]));
  const { updateNode, removeNode } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeNode"]));
  const { confirm: confirmPanel } = usePanelContext();

  const {
    form: formInst,
    formPending,
    formProps,
    submit: submitForm,
  } = useAntdForm({
    name: "workflowNodeForm",
    onSubmit: async (values) => {
      if (node.type === WorkflowNodeType.Apply) {
        await addEmail(values.contactEmail);
        await updateNode(
          produce(node, (draft) => {
            draft.config = {
              provider: accesses.find((e) => e.id === values.providerAccessId)?.provider,
              ...values,
            };
            draft.validated = true;
          })
        );
      } else {
        await updateNode(
          produce(node, (draft) => {
            draft.config = { ...values };
            draft.validated = true;
          })
        );
      }
    },
  });

  const nodeContentComponent = useMemo(() => {
    if (!node.validated) {
      return <Typography.Link>{t("workflow_node.action.configure_node")}</Typography.Link>;
    }

    switch (node.type) {
      case WorkflowNodeType.Start: {
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
      }

      case WorkflowNodeType.Apply: {
        const config = (node.config as WorkflowNodeConfigForApply) ?? {};
        return <Typography.Text className="truncate">{config.domains || "　"}</Typography.Text>;
      }

      case WorkflowNodeType.Deploy: {
        const config = (node.config as WorkflowNodeConfigForDeploy) ?? {};
        const provider = deployProvidersMap.get(config.provider);
        return (
          <Space>
            <Avatar src={provider?.icon} size="small" />
            <Typography.Text className="truncate">{t(provider?.name ?? "")}</Typography.Text>
          </Space>
        );
      }

      case WorkflowNodeType.Notify: {
        const config = (node.config as WorkflowNodeConfigForNotify) ?? {};
        const channel = notifyChannelsMap.get(config.channel as string);
        return (
          <div className="flex items-center justify-between space-x-2">
            <Typography.Text className="truncate">{t(channel?.name ?? "　")}</Typography.Text>
            <Typography.Text className="truncate" type="secondary">
              {config.subject ?? ""}
            </Typography.Text>
          </div>
        );
      }

      default: {
        console.warn(`[certimate] unsupported workflow node type: ${node.type}`);
        return <></>;
      }
    }
  }, [node]);

  const panelBodyComponent = useMemo(() => {
    const nodeFormProps = {
      form: formInst,
      formName: formProps.name,
      disabled: disabled || formPending,
      workflowNode: node,
    };

    switch (node.type) {
      case WorkflowNodeType.Start:
        return <StartNodeForm {...nodeFormProps} />;
      case WorkflowNodeType.Apply:
        return <ApplyNodeForm {...nodeFormProps} />;
      case WorkflowNodeType.Deploy:
        return <DeployNodeForm {...nodeFormProps} />;
      case WorkflowNodeType.Notify:
        return <NotifyNodeForm {...nodeFormProps} />;
      default:
        console.warn(`[certimate] unsupported workflow node type: ${node.type}`);
        return <> </>;
    }
  }, [node, disabled, formInst, formPending, formProps]);

  const handleNodeClick = () => {
    confirmPanel({
      title: node.name,
      children: panelBodyComponent,
      okText: t("common.button.save"),
      onOk: () => {
        return submitForm();
      },
    });
  };

  const handleNodeNameBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    const oldName = node.name;
    const newName = e.target.innerText.trim();
    if (oldName === newName) {
      return;
    }

    updateNode(
      produce(node, (draft) => {
        draft.name = newName;
      })
    );
  };

  return (
    <>
      <Popover
        arrow={false}
        content={
          <Show when={node.type !== WorkflowNodeType.Start}>
            <Dropdown
              menu={{
                items: [
                  {
                    key: "delete",
                    disabled: disabled,
                    label: t("workflow_node.action.delete_node"),
                    icon: <CloseCircleOutlinedIcon />,
                    danger: true,
                    onClick: () => {
                      if (disabled) return;

                      removeNode(node.id);
                    },
                  },
                ],
              }}
              trigger={["click"]}
            >
              <Button color="primary" icon={<EllipsisOutlinedIcon />} variant="text" />
            </Dropdown>
          </Show>
        }
        overlayClassName="shadow-md"
        overlayInnerStyle={{ padding: 0 }}
        placement="rightTop"
      >
        <Card className="relative w-[256px] overflow-hidden shadow-md" styles={{ body: { padding: 0 } }} hoverable>
          <div className="bg-primary flex h-[48px] flex-col items-center justify-center truncate px-4 py-2 text-white">
            <div
              className="focus:bg-background focus:text-foreground w-full overflow-hidden text-center outline-none focus:rounded-sm"
              contentEditable
              suppressContentEditableWarning
              onBlur={handleNodeNameBlur}
            >
              {node.name}
            </div>
          </div>

          <div className="flex flex-col justify-center px-4 py-2">
            <div className="cursor-pointer text-sm" onClick={handleNodeClick}>
              {nodeContentComponent}
            </div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} disabled={disabled} />
    </>
  );
};

export default memo(CommonNode);
