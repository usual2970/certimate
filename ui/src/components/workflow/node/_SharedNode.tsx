import { memo } from "react";
import { useTranslation } from "react-i18next";
import { CloseCircleOutlined as CloseCircleOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { useControllableValue } from "ahooks";
import { Button, Card, Drawer, Dropdown, Modal, Popover, Space } from "antd";
import { produce } from "immer";
import { isEqual } from "radash";

import Show from "@/components/Show";
import { type WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import AddNode from "./AddNode";

export type SharedNodeProps = {
  node: WorkflowNode;
  disabled?: boolean;
};

type SharedNodeWrapperProps = SharedNodeProps & {
  children: React.ReactNode;
  onClick?: (e: React.MouseEvent) => void;
};

const SharedNodeWrapper = ({ children, node, disabled, onClick }: SharedNodeWrapperProps) => {
  const { t } = useTranslation();

  const { updateNode, removeNode } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeNode"]));

  const handleNodeClick = (e: React.MouseEvent) => {
    onClick?.(e);
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

          <div className="flex cursor-pointer flex-col justify-center px-4 py-2" onClick={handleNodeClick}>
            <div className="text-sm">{children}</div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} disabled={disabled} />
    </>
  );
};

type SharedNodeEditDrawerProps = SharedNodeProps & {
  children: React.ReactNode;
  footer?: boolean;
  loading?: boolean;
  open?: boolean;
  pending?: boolean;
  onOpenChange?: (open: boolean) => void;
  onConfirm: () => void | Promise<unknown>;
  getFormValues: () => NonNullable<unknown>;
};

const SharedNodeConfigDrawer = ({
  children,
  node,
  disabled,
  footer = true,
  loading,
  pending,
  onConfirm,
  getFormValues,
  ...props
}: SharedNodeEditDrawerProps) => {
  const { t } = useTranslation();

  const [modalApi, ModelContextHolder] = Modal.useModal();

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const handleConfirmClick = async () => {
    await onConfirm();
    setOpen(false);
  };

  const handleCancelClick = () => {
    if (pending) return;

    setOpen(false);
  };

  const handleClose = () => {
    if (pending) return;

    const oldValues = Object.fromEntries(Object.entries(node.config ?? {}).filter(([_, value]) => value !== null && value !== undefined));
    const newValues = Object.fromEntries(Object.entries(getFormValues()).filter(([_, value]) => value !== null && value !== undefined));
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

    promise.then(() => setOpen(false));
  };

  return (
    <>
      {ModelContextHolder}

      <Drawer
        afterOpenChange={(open) => setOpen(open)}
        destroyOnClose
        loading={loading}
        footer={
          !!footer && (
            <Space className="w-full justify-end">
              <Button onClick={handleCancelClick}>{t("common.button.cancel")}</Button>
              <Button disabled={disabled} loading={pending} type="primary" onClick={handleConfirmClick}>
                {t("common.button.save")}
              </Button>
            </Space>
          )
        }
        open={open}
        width={640}
        onClose={handleClose}
      >
        {children}
      </Drawer>
    </>
  );
};

export default {
  Wrapper: memo(SharedNodeWrapper),
  ConfigDrawer: memo(SharedNodeConfigDrawer),
};
