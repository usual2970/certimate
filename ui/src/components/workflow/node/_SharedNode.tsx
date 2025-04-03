import { memo, useRef } from "react";
import { useTranslation } from "react-i18next";
import {
  CloseCircleOutlined as CloseCircleOutlinedIcon,
  EllipsisOutlined as EllipsisOutlinedIcon,
  FormOutlined as FormOutlinedIcon,
  MoreOutlined as MoreOutlinedIcon,
} from "@ant-design/icons";
import { useControllableValue } from "ahooks";
import { Button, Card, Drawer, Dropdown, Input, type InputRef, Modal, Popover, Space } from "antd";
import { produce } from "immer";
import { isEqual } from "radash";

import { type WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import AddNode from "./AddNode";

export type SharedNodeProps = {
  node: WorkflowNode;
  disabled?: boolean;
};

// #region Title
type SharedNodeTitleProps = SharedNodeProps & {
  className?: string;
  style?: React.CSSProperties;
};

const SharedNodeTitle = ({ className, style, node, disabled }: SharedNodeTitleProps) => {
  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const handleBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    const oldName = node.name;
    const newName = e.target.innerText.trim().substring(0, 64) || oldName;
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
    <div className="w-full cursor-text overflow-hidden text-center">
      <div className={className} style={style} contentEditable={!disabled} suppressContentEditableWarning onBlur={handleBlur}>
        {node.name}
      </div>
    </div>
  );
};
// #endregion

// #region Menu
type SharedNodeMenuProps = SharedNodeProps & {
  branchId?: string;
  branchIndex?: number;
  trigger: React.ReactNode;
  afterUpdate?: () => void;
  afterDelete?: () => void;
};

const isBranchingNode = (node: WorkflowNode) => {
  return (
    node.type === WorkflowNodeType.Branch ||
    node.type === WorkflowNodeType.Condition ||
    node.type === WorkflowNodeType.ExecuteResultBranch ||
    node.type === WorkflowNodeType.ExecuteSuccess ||
    node.type === WorkflowNodeType.ExecuteFailure
  );
};

const SharedNodeMenu = ({ trigger, node, disabled, branchId, branchIndex, afterUpdate, afterDelete }: SharedNodeMenuProps) => {
  const { t } = useTranslation();

  const { updateNode, removeNode, removeBranch } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeNode", "removeBranch"]));

  const [modalApi, ModelContextHolder] = Modal.useModal();

  const nameInputRef = useRef<InputRef>(null);
  const nameRef = useRef<string>();

  const handleRenameConfirm = async () => {
    const oldName = node.name;
    const newName = nameRef.current?.trim()?.substring(0, 64) || oldName;
    if (oldName === newName) {
      return;
    }

    await updateNode(
      produce(node, (draft) => {
        draft.name = newName;
      })
    );

    afterUpdate?.();
  };

  const handleDeleteClick = async () => {
    if (isBranchingNode(node)) {
      await removeBranch(branchId!, branchIndex!);
    } else {
      await removeNode(node.id);
    }

    afterDelete?.();
  };

  return (
    <>
      {ModelContextHolder}

      <Dropdown
        menu={{
          items: [
            {
              key: "rename",
              disabled: disabled,
              label: isBranchingNode(node) ? t("workflow_node.action.rename_branch") : t("workflow_node.action.rename_node"),
              icon: <FormOutlinedIcon />,
              onClick: () => {
                nameRef.current = node.name;

                const dialog = modalApi.confirm({
                  title: isBranchingNode(node) ? t("workflow_node.action.rename_branch") : t("workflow_node.action.rename_node"),
                  content: (
                    <div className="pb-2 pt-4">
                      <Input
                        ref={nameInputRef}
                        autoFocus
                        defaultValue={node.name}
                        onChange={(e) => (nameRef.current = e.target.value)}
                        onPressEnter={async () => {
                          await handleRenameConfirm();
                          dialog.destroy();
                        }}
                      />
                    </div>
                  ),
                  icon: null,
                  okText: t("common.button.save"),
                  onOk: handleRenameConfirm,
                });
                setTimeout(() => nameInputRef.current?.focus(), 1);
              },
            },
            {
              type: "divider",
            },
            {
              key: "remove",
              disabled: disabled || node.type === WorkflowNodeType.Start,
              label: isBranchingNode(node) ? t("workflow_node.action.remove_branch") : t("workflow_node.action.remove_node"),
              icon: <CloseCircleOutlinedIcon />,
              danger: true,
              onClick: handleDeleteClick,
            },
          ],
        }}
        trigger={["click"]}
      >
        {trigger}
      </Dropdown>
    </>
  );
};
// #endregion

// #region Wrapper
type SharedNodeBlockProps = SharedNodeProps & {
  children: React.ReactNode;
  onClick?: (e: React.MouseEvent) => void;
};

const SharedNodeBlock = ({ children, node, disabled, onClick }: SharedNodeBlockProps) => {
  const handleNodeClick = (e: React.MouseEvent) => {
    onClick?.(e);
  };

  return (
    <>
      <Popover
        classNames={{ root: "shadow-md" }}
        styles={{ body: { padding: 0 } }}
        arrow={false}
        content={<SharedNodeMenu node={node} disabled={disabled} trigger={<Button color="primary" icon={<MoreOutlinedIcon />} variant="text" />} />}
        placement="rightTop"
      >
        <Card className="relative w-[256px] overflow-hidden shadow-md" styles={{ body: { padding: 0 } }} hoverable>
          <div className="bg-primary flex h-[48px] flex-col items-center justify-center truncate px-4 py-2 text-white">
            <SharedNodeTitle
              className="focus:bg-background focus:text-foreground overflow-hidden outline-none focus:rounded-sm"
              node={node}
              disabled={disabled}
            />
          </div>

          <div className="flex cursor-pointer flex-col justify-center px-4 py-2" onClick={handleNodeClick}>
            <div className="overflow-hidden text-sm">{children}</div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} disabled={disabled} />
    </>
  );
};
// #endregion

// #region EditDrawer
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

    const oldValues = JSON.parse(JSON.stringify(node.config ?? {}));
    const newValues = JSON.parse(JSON.stringify(getFormValues()));
    const changed = !isEqual(oldValues, {}) && !isEqual(oldValues, newValues);

    const { promise, resolve, reject } = Promise.withResolvers();
    if (changed) {
      console.log(oldValues, newValues);
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
        afterOpenChange={setOpen}
        closable={!pending}
        destroyOnClose
        extra={
          <SharedNodeMenu
            node={node}
            disabled={disabled}
            trigger={<Button icon={<EllipsisOutlinedIcon />} type="text" />}
            afterDelete={() => {
              setOpen(false);
            }}
          />
        }
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
        loading={loading}
        maskClosable={!pending}
        open={open}
        title={<div className="max-w-[480px] truncate">{node.name}</div>}
        width={720}
        onClose={handleClose}
      >
        {children}
      </Drawer>
    </>
  );
};
// #endregion

export default {
  Title: memo(SharedNodeTitle),
  Menu: memo(SharedNodeMenu),
  Block: memo(SharedNodeBlock),
  ConfigDrawer: memo(SharedNodeConfigDrawer),
};
