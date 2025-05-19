import { memo, useRef, useState } from "react";
import { MoreOutlined as MoreOutlinedIcon } from "@ant-design/icons";
import { Button, Card, Popover } from "antd";

import SharedNode, { type SharedNodeProps } from "./_SharedNode";
import AddNode from "./AddNode";
import ConditionNodeConfigForm, { ConditionNodeConfigFormInstance } from "./ConditionNodeConfigForm";
import { WorkflowNodeConfigForCondition } from "@/domain/workflow";
import { produce } from "immer";
import { useWorkflowStore } from "@/stores/workflow";
import { useZustandShallowSelector } from "@/hooks";

export type ConditionNodeProps = SharedNodeProps & {
  branchId: string;
  branchIndex: number;
};

const ConditionNode = ({ node, disabled, branchId, branchIndex }: ConditionNodeProps) => {
  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const [formPending, setFormPending] = useState(false);
  const formRef = useRef<ConditionNodeConfigFormInstance>(null);

  const [drawerOpen, setDrawerOpen] = useState(false);

  const getFormValues = () => formRef.current!.getFieldsValue() as WorkflowNodeConfigForCondition;

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
      <Popover
        classNames={{ root: "shadow-md" }}
        styles={{ body: { padding: 0 } }}
        arrow={false}
        content={
          <SharedNode.Menu
            node={node}
            branchId={branchId}
            branchIndex={branchIndex}
            disabled={disabled}
            trigger={<Button color="primary" icon={<MoreOutlinedIcon />} variant="text" />}
          />
        }
        placement="rightTop"
      >
        <Card className="relative z-[1] mt-10 w-[256px] shadow-md" styles={{ body: { padding: 0 } }} hoverable onClick={() => setDrawerOpen(true)}>
          <div className="flex h-[48px] flex-col items-center justify-center truncate px-4 py-2">
            <SharedNode.Title
              className="focus:bg-background focus:text-foreground overflow-hidden outline-slate-200 focus:rounded-sm"
              node={node}
              disabled={disabled}
            />
          </div>
        </Card>

        <SharedNode.ConfigDrawer
          node={node}
          open={drawerOpen}
          pending={formPending}
          onConfirm={handleDrawerConfirm}
          onOpenChange={(open) => setDrawerOpen(open)}
          getFormValues={() => formRef.current!.getFieldsValue()}
        >
          <ConditionNodeConfigForm nodeId={node.id} ref={formRef} disabled={disabled} initialValues={node.config} />
        </SharedNode.ConfigDrawer>
      </Popover>

      <AddNode node={node} disabled={disabled} />
    </>
  );
};

export default memo(ConditionNode);

