import { memo, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { CloseCircleOutlined as CloseCircleOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { useControllableValue } from "ahooks";
import { Avatar, Button, Card, Drawer, Dropdown, Modal, Popover, Space, Typography } from "antd";
import { produce } from "immer";
import { isEqual } from "radash";

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

export type CommonNodeProps = {
  node: WorkflowNode;
  disabled?: boolean;
};

const CommonNode = ({ node, disabled }: CommonNodeProps) => {
  const { t } = useTranslation();

  const { updateNode, removeNode } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeNode"]));

  const [drawerOpen, setDrawerOpen] = useState(false);

  const workflowNodeEl = useMemo(() => {
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

  const handleNodeClick = () => {
    setDrawerOpen(true);
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
              {workflowNodeEl}
            </div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} disabled={disabled} />

      <CommonNodeEditDrawer node={node} disabled={disabled} open={drawerOpen} onOpenChange={(open) => setDrawerOpen(open)} />
    </>
  );
};

type CommonNodeEditDrawerProps = CommonNodeProps & {
  defaultOpen?: boolean;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
};

const CommonNodeEditDrawer = ({ node, disabled, ...props }: CommonNodeEditDrawerProps) => {
  const { t } = useTranslation();

  const [modalApi, ModelContextHolder] = Modal.useModal();

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const { accesses } = useAccessesStore(useZustandShallowSelector("accesses"));
  const { addEmail } = useContactEmailsStore(useZustandShallowSelector(["addEmail"]));
  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const {
    form: formInst,
    formPending,
    formProps,
    submit: submitForm,
  } = useAntdForm({
    name: "workflowNodeForm",
    onSubmit: async (values) => {
      await sleep(5000);
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

  const formEl = useMemo(() => {
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

  const handleClose = () => {
    if (formPending) return;

    const oldValues = Object.fromEntries(Object.entries(node.config ?? {}).filter(([_, value]) => value !== null && value !== undefined));
    const newValues = Object.fromEntries(Object.entries(formInst.getFieldsValue(true)).filter(([_, value]) => value !== null && value !== undefined));
    const changed = !isEqual(oldValues, newValues);

    const { promise, resolve, reject } = Promise.withResolvers();
    if (changed) {
      modalApi.confirm({
        title: t("common.text.operation_confirm"),
        content: t("workflow_node.unsaved_changes.confirm"),
        onOk: () => resolve(void 0),
        onCancel: () => reject(),
      });
    } else {
      resolve(void 0);
    }

    promise.then(() => {
      setOpen(false);
    });
  };

  const handleCancelClick = () => {
    if (formPending) return;

    setOpen(false);
  };

  const handleOkClick = async () => {
    await submitForm();
    setOpen(false);
  };

  return (
    <>
      {ModelContextHolder}

      <Drawer
        destroyOnClose
        footer={
          <Space className="w-full justify-end">
            <Button onClick={handleCancelClick}>{t("common.button.cancel")}</Button>
            <Button loading={formPending} type="primary" onClick={handleOkClick}>
              {t("common.button.ok")}
            </Button>
          </Space>
        }
        open={open}
        title={node.name}
        width={640}
        onClose={handleClose}
      >
        {formEl}
      </Drawer>
    </>
  );
};

export default memo(CommonNode);
