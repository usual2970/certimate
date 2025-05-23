import { memo, useRef, useState } from "react";
import { MoreOutlined as MoreOutlinedIcon } from "@ant-design/icons";
import { Button, Card, Popover } from "antd";

import SharedNode, { type SharedNodeProps } from "./_SharedNode";
import AddNode from "./AddNode";
import ConditionNodeConfigForm, { ConditionItem, ConditionNodeConfigFormFieldValues, ConditionNodeConfigFormInstance } from "./ConditionNodeConfigForm";
import { Expr, WorkflowNodeIoValueType, ExprType } from "@/domain/workflow";
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

  const getFormValues = () => formRef.current!.getFieldsValue() as ConditionNodeConfigFormFieldValues;

  // 将表单值转换为表达式结构
  const formToExpression = (values: ConditionNodeConfigFormFieldValues): Expr => {
    // 创建单个条件的表达式
    const createComparisonExpr = (condition: ConditionItem): Expr => {
      const selectors = condition.leftSelector.split("#");
      const t = selectors[2] as WorkflowNodeIoValueType;
      const left: Expr = {
        type: ExprType.Var,
        selector: {
          id: selectors[0],
          name: selectors[1],
          type: t,
        },
      };

      const right: Expr = { type: ExprType.Const, value: condition.rightValue, valueType: t };

      return {
        type: ExprType.Compare,
        op: condition.operator,
        left,
        right,
      };
    };

    // 如果只有一个条件，直接返回比较表达式
    if (values.conditions.length === 1) {
      return createComparisonExpr(values.conditions[0]);
    }

    // 多个条件，通过逻辑运算符连接
    let expr: Expr = createComparisonExpr(values.conditions[0]);

    for (let i = 1; i < values.conditions.length; i++) {
      expr = {
        type: ExprType.Logical,
        op: values.logicalOperator,
        left: expr,
        right: createComparisonExpr(values.conditions[i]),
      };
    }

    return expr;
  };

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
      const expression = formToExpression(newValues);
      const newNode = produce(node, (draft) => {
        draft.config = {
          expression,
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
