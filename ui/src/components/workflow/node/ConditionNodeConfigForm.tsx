import { forwardRef, memo, useEffect, useImperativeHandle, useState } from "react";
import { Button, Card, Form, Input, Select, Space, Radio, DatePicker } from "antd";
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
  workflowNodeIOOptions,
  WorkflowNodeIoValueType,
} from "@/domain/workflow";
import { FormInstance } from "antd";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

// 表单内部使用的扁平结构 - 修改后只保留必要字段
export interface ConditionItem {
  leftSelector: string;
  operator: ComparisonOperator;
  rightValue: string;
}

export type ConditionNodeConfigFormFieldValues = {
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
        leftSelector: "",
        operator: "==",
        rightValue: "",
      },
    ],
    logicalOperator: "and",
  };
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
          leftSelector: `${expr.left.selector.id}#${expr.left.selector.name}#${expr.left.selector.type}`,
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

// 根据变量类型获取适当的操作符选项
const getOperatorsByType = (type: string): { value: ComparisonOperator; label: string }[] => {
  switch (type) {
    case "number":
    case "string":
      return [
        { value: "==", label: "等于 (==)" },
        { value: "!=", label: "不等于 (!=)" },
        { value: ">", label: "大于 (>)" },
        { value: ">=", label: "大于等于 (>=)" },
        { value: "<", label: "小于 (<)" },
        { value: "<=", label: "小于等于 (<=)" },
      ];
    case "boolean":
      return [{ value: "is", label: "为" }];
    default:
      return [];
  }
};

// 从选择器字符串中提取变量类型
const getVariableTypeFromSelector = (selector: string): string => {
  if (!selector) return "string";

  // 假设选择器格式为 "id#name#type"
  const parts = selector.split("#");
  if (parts.length >= 3) {
    return parts[2].toLowerCase() || "string";
  }
  return "string";
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
    const handleFormChange = (_: undefined, values: ConditionNodeConfigFormFieldValues) => {
      setFormModel(values);

      if (onValuesChange) {
        // 将表单值转换为表达式
        const expression = formToExpression(values);
        onValuesChange({ expression });
      }
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
                  <div className="flex items-center gap-2">
                    {/* 左侧变量选择器 */}
                    <Form.Item {...restField} name={[name, "leftSelector"]} className="mb-0 flex-1" rules={[{ required: true, message: "请选择变量" }]}>
                      <Select
                        placeholder="选择变量"
                        options={previousNodes.map((item) => {
                          return workflowNodeIOOptions(item);
                        })}
                      />
                    </Form.Item>

                    {/* 操作符 - 动态根据变量类型改变选项 */}
                    <Form.Item
                      noStyle
                      shouldUpdate={(prevValues, currentValues) => {
                        return prevValues.conditions?.[name]?.leftSelector !== currentValues.conditions?.[name]?.leftSelector;
                      }}
                    >
                      {({ getFieldValue }) => {
                        const leftSelector = getFieldValue(["conditions", name, "leftSelector"]);
                        const varType = getVariableTypeFromSelector(leftSelector);
                        const operators = getOperatorsByType(varType);

                        return (
                          <Form.Item {...restField} name={[name, "operator"]} className="mb-0 w-32" rules={[{ required: true, message: "请选择" }]}>
                            <Select options={operators} />
                          </Form.Item>
                        );
                      }}
                    </Form.Item>

                    {/* 右侧输入控件 - 根据变量类型使用不同的控件 */}
                    <Form.Item
                      noStyle
                      shouldUpdate={(prevValues, currentValues) => {
                        return prevValues.conditions?.[name]?.leftSelector !== currentValues.conditions?.[name]?.leftSelector;
                      }}
                    >
                      {({ getFieldValue }) => {
                        const leftSelector = getFieldValue(["conditions", name, "leftSelector"]);
                        const varType = getVariableTypeFromSelector(leftSelector);

                        return (
                          <Form.Item {...restField} name={[name, "rightValue"]} className="mb-0 flex-1" rules={[{ required: true, message: "请输入值" }]}>
                            {varType === "boolean" ? (
                              <Select placeholder="选择值">
                                <Select.Option value="true">是</Select.Option>
                                <Select.Option value="false">否</Select.Option>
                              </Select>
                            ) : varType === "number" ? (
                              <Input type="number" placeholder="输入数值" />
                            ) : varType === "date" ? (
                              <DatePicker style={{ width: "100%" }} placeholder="选择日期" format="YYYY-MM-DD" />
                            ) : (
                              <Input placeholder="输入值" />
                            )}
                          </Form.Item>
                        );
                      }}
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
                      leftSelector: "",
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

// 表单值转换为表达式结构 (需要添加)
const formToExpression = (values: ConditionNodeConfigFormFieldValues): Expr => {
  const createComparisonExpr = (condition: ConditionItem): Expr => {
    const [id, name, typeStr] = condition.leftSelector.split("#");

    const type = typeStr as WorkflowNodeIoValueType;

    const left: Expr = {
      type: "var",
      selector: { id, name, type },
    };

    let rightValue: any = condition.rightValue;
    if (type === "number") {
      rightValue = Number(condition.rightValue);
    } else if (type === "boolean") {
      rightValue = condition.rightValue === "true";
    }

    const right: Expr = {
      type: "const",
      value: rightValue,
    };

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

export default memo(ConditionNodeConfigForm);

