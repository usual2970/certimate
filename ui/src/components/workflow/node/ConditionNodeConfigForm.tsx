import { forwardRef, memo, useEffect, useImperativeHandle, useState } from "react";
import { Button, Card, Form, Input, Select, Space, Radio } from "antd";
import { PlusOutlined, DeleteOutlined } from "@ant-design/icons";

import {
  WorkflowNodeConfigForCondition,
  Expr,
  WorkflowNodeIOValueSelector,
  ComparisonOperator,
  LogicalOperator,
  isConstExpr,
  isVarExpr,
  WorkflowNode,
} from "@/domain/workflow";
import { FormInstance } from "antd";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

// 表单内部使用的扁平结构 - 修改后只保留必要字段
interface ConditionItem {
  leftSelector: WorkflowNodeIOValueSelector;
  operator: ComparisonOperator;
  rightValue: string;
}

type ConditionNodeConfigFormFieldValues = {
  conditions: ConditionItem[];
  logicalOperator: LogicalOperator;
};

export type ConditionNodeConfigFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: Partial<WorkflowNodeConfigForCondition>;
  onValuesChange?: (values: WorkflowNodeConfigForCondition) => void;
  availableSelectors?: WorkflowNodeIOValueSelector[];
  nodeId: string;
};

export type ConditionNodeConfigFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<ConditionNodeConfigFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<ConditionNodeConfigFormFieldValues>["resetFields"];
  validateFields: FormInstance<ConditionNodeConfigFormFieldValues>["validateFields"];
};

// 初始表单值
const initFormModel = (): ConditionNodeConfigFormFieldValues => {
  return {
    conditions: [
      {
        leftSelector: undefined as unknown as WorkflowNodeIOValueSelector,
        operator: "==",
        rightValue: "",
      },
    ],
    logicalOperator: "and",
  };
};

// 将表单值转换为表达式结构
const formToExpression = (values: ConditionNodeConfigFormFieldValues): Expr => {
  // 创建单个条件的表达式
  const createComparisonExpr = (condition: ConditionItem): Expr => {
    const left: Expr = { type: "var", selector: condition.leftSelector };
    const right: Expr = { type: "const", value: condition.rightValue || "" };

    return {
      type: "compare",
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
      type: "logical",
      op: values.logicalOperator,
      left: expr,
      right: createComparisonExpr(values.conditions[i]),
    };
  }

  return expr;
};

// 递归提取表达式中的条件项
const expressionToForm = (expr?: Expr): ConditionNodeConfigFormFieldValues => {
  if (!expr) return initFormModel();

  const conditions: ConditionItem[] = [];
  let logicalOp: LogicalOperator = "and";

  const extractComparisons = (expr: Expr): void => {
    if (expr.type === "compare") {
      // 确保左侧是变量，右侧是常量
      if (isVarExpr(expr.left) && isConstExpr(expr.right)) {
        conditions.push({
          leftSelector: expr.left.selector,
          operator: expr.op,
          rightValue: String(expr.right.value),
        });
      }
    } else if (expr.type === "logical") {
      logicalOp = expr.op;
      extractComparisons(expr.left);
      extractComparisons(expr.right);
    }
  };

  extractComparisons(expr);

  return {
    conditions: conditions.length > 0 ? conditions : initFormModel().conditions,
    logicalOperator: logicalOp,
  };
};

const ConditionNodeConfigForm = forwardRef<ConditionNodeConfigFormInstance, ConditionNodeConfigFormProps>(
  ({ className, style, disabled, initialValues, onValuesChange, nodeId }, ref) => {
    const { getWorkflowOuptutBeforeId } = useWorkflowStore(useZustandShallowSelector(["updateNode", "getWorkflowOuptutBeforeId"]));

    const [form] = Form.useForm<ConditionNodeConfigFormFieldValues>();
    const [formModel, setFormModel] = useState<ConditionNodeConfigFormFieldValues>(initFormModel());

    const [previousNodes, setPreviousNodes] = useState<WorkflowNode[]>([]);
    useEffect(() => {
      const previousNodes = getWorkflowOuptutBeforeId(nodeId);
      setPreviousNodes(previousNodes);
    }, [nodeId]);

    // 初始化表单值
    useEffect(() => {
      if (initialValues?.expression) {
        const formValues = expressionToForm(initialValues.expression);
        form.setFieldsValue(formValues);
        setFormModel(formValues);
      }
    }, [form, initialValues]);

    // 公开表单方法
    useImperativeHandle(
      ref,
      () => ({
        getFieldsValue: form.getFieldsValue,
        resetFields: form.resetFields,
        validateFields: form.validateFields,
      }),
      [form]
    );

    // 表单值变更处理
    const handleFormChange = (changedValues: any, values: ConditionNodeConfigFormFieldValues) => {
      setFormModel(values);

      // 转换为表达式结构并通知父组件
      const expression = formToExpression(values);
      onValuesChange?.({ expression });
    };

    return (
      <Form form={form} className={className} style={style} layout="vertical" disabled={disabled} initialValues={formModel} onValuesChange={handleFormChange}>
        <Form.List name="conditions">
          {(fields, { add, remove }) => (
            <>
              {fields.map(({ key, name, ...restField }) => (
                <Card
                  key={key}
                  size="small"
                  className="mb-3"
                  extra={fields.length > 1 ? <Button icon={<DeleteOutlined />} danger type="text" onClick={() => remove(name)} /> : null}
                >
                  {/* 将三个表单项放在一行 */}
                  <div className="flex items-center gap-2">
                    {/* 左侧变量选择器 */}
                    <Form.Item {...restField} name={[name, "leftSelector"]} className="mb-0 flex-1" rules={[{ required: true, message: "请选择变量" }]}>
                      <Select
                        placeholder="选择变量"
                        options={previousNodes.map((item) => {
                          return {
                            label: item.name,
                            options: item.outputs?.map((output) => {
                              return {
                                label: `${item.name} - ${output.label}`,
                                value: `${item.id}#${output.name}`,
                              };
                            }),
                          };
                        })}
                      ></Select>
                    </Form.Item>

                    {/* 操作符 */}
                    <Form.Item {...restField} name={[name, "operator"]} className="mb-0 w-32" rules={[{ required: true, message: "请选择" }]}>
                      <Select>
                        <Select.Option value="==">等于 (==)</Select.Option>
                        <Select.Option value="!=">不等于 (!=)</Select.Option>
                        <Select.Option value=">">大于 (&gt;)</Select.Option>
                        <Select.Option value=">=">大于等于 (&gt;=)</Select.Option>
                        <Select.Option value="<">小于 (&lt;)</Select.Option>
                        <Select.Option value="<=">小于等于 (&lt;=)</Select.Option>
                      </Select>
                    </Form.Item>

                    {/* 右侧常量输入框 */}
                    <Form.Item {...restField} name={[name, "rightValue"]} className="mb-0 flex-1" rules={[{ required: true, message: "请输入值" }]}>
                      <Input placeholder="输入值" />
                    </Form.Item>
                  </div>
                </Card>
              ))}

              {/* 添加条件按钮 */}
              <Form.Item>
                <Button
                  type="dashed"
                  onClick={() =>
                    add({
                      leftSelector: undefined as unknown as WorkflowNodeIOValueSelector,
                      operator: "==",
                      rightValue: "",
                    })
                  }
                  block
                  icon={<PlusOutlined />}
                >
                  添加条件
                </Button>
              </Form.Item>
            </>
          )}
        </Form.List>

        {formModel.conditions && formModel.conditions.length > 1 && (
          <Form.Item name="logicalOperator" label="条件逻辑">
            <Radio.Group buttonStyle="solid">
              <Radio.Button value="and">满足所有条件 (AND)</Radio.Button>
              <Radio.Button value="or">满足任一条件 (OR)</Radio.Button>
            </Radio.Group>
          </Form.Item>
        )}
      </Form>
    );
  }
);

export default memo(ConditionNodeConfigForm);

