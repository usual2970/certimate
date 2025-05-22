import { forwardRef, memo, useEffect, useImperativeHandle, useState } from "react";
import { Button, Card, Form, Input, Select, Radio } from "antd";
import { PlusOutlined, DeleteOutlined } from "@ant-design/icons";
import i18n from "@/i18n";

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
  ExprType,
} from "@/domain/workflow";
import { FormInstance } from "antd";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";
import { useTranslation } from "react-i18next";

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
    logicalOperator: LogicalOperator.And,
  };
};

// 递归提取表达式中的条件项
const expressionToForm = (expr?: Expr): ConditionNodeConfigFormFieldValues => {
  if (!expr) return initFormModel();

  const conditions: ConditionItem[] = [];
  let logicalOp: LogicalOperator = LogicalOperator.And;

  const extractComparisons = (expr: Expr): void => {
    if (expr.type === ExprType.Compare) {
      // 确保左侧是变量，右侧是常量
      if (isVarExpr(expr.left) && isConstExpr(expr.right)) {
        conditions.push({
          leftSelector: `${expr.left.selector.id}#${expr.left.selector.name}#${expr.left.selector.type}`,
          operator: expr.op,
          rightValue: String(expr.right.value),
        });
      }
    } else if (expr.type === ExprType.Logical) {
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
        { value: "==", label: i18n.t("workflow_node.condition.form.comparison.equal") },
        { value: "!=", label: i18n.t("workflow_node.condition.form.comparison.not_equal") },
        { value: ">", label: i18n.t("workflow_node.condition.form.comparison.greater_than") },
        { value: ">=", label: i18n.t("workflow_node.condition.form.comparison.greater_than_or_equal") },
        { value: "<", label: i18n.t("workflow_node.condition.form.comparison.less_than") },
        { value: "<=", label: i18n.t("workflow_node.condition.form.comparison.less_than_or_equal") },
      ];
    case "boolean":
      return [{ value: "is", label: i18n.t("workflow_node.condition.form.comparison.is") }];
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
    const { t } = useTranslation();
    const prefix = "workflow_node.condition.form";

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
                    <Form.Item
                      {...restField}
                      name={[name, "leftSelector"]}
                      className="mb-0 flex-1"
                      rules={[{ required: true, message: t(`${prefix}.variable.errmsg`) }]}
                    >
                      <Select
                        placeholder={t(`${prefix}.variable.placeholder`)}
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
                          <Form.Item
                            {...restField}
                            name={[name, "operator"]}
                            className="mb-0 w-32"
                            rules={[{ required: true, message: t(`${prefix}.operator.errmsg`) }]}
                          >
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
                          <Form.Item
                            {...restField}
                            name={[name, "rightValue"]}
                            className="mb-0 flex-1"
                            rules={[{ required: true, message: t(`${prefix}.value.errmsg`) }]}
                          >
                            {varType === "boolean" ? (
                              <Select placeholder={t(`${prefix}.value.boolean.placeholder`)}>
                                <Select.Option value="true">{t(`${prefix}.value.boolean.true`)}</Select.Option>
                                <Select.Option value="false">{t(`${prefix}.value.boolean.false`)}</Select.Option>
                              </Select>
                            ) : varType === "number" ? (
                              <Input type="number" placeholder={t(`${prefix}.value.number.placeholder`)} />
                            ) : (
                              <Input placeholder={t(`${prefix}.value.string.placeholder`)} />
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
                  {t(`${prefix}.add_condition.button`)}
                </Button>
              </Form.Item>
            </>
          )}
        </Form.List>

        {formModel.conditions && formModel.conditions.length > 1 && (
          <Form.Item name="logicalOperator" label={t(`${prefix}.logical_operator.label`)}>
            <Radio.Group buttonStyle="solid">
              <Radio.Button value="and">{t(`${prefix}.logical_operator.and`)}</Radio.Button>
              <Radio.Button value="or">{t(`${prefix}.logical_operator.or`)}</Radio.Button>
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
      type: ExprType.Var,
      selector: { id, name, type },
    };

    const right: Expr = {
      type: ExprType.Const,
      value: condition.rightValue,
      valueType: type,
    };

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

export default memo(ConditionNodeConfigForm);
