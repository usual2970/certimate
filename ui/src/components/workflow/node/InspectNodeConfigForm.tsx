import { forwardRef, memo, useImperativeHandle } from "react";
import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type WorkflowNodeConfigForInspect } from "@/domain/workflow";
import { useAntdForm } from "@/hooks";

import { validDomainName, validIPv4Address, validPortNumber } from "@/utils/validators";

type InspectNodeConfigFormFieldValues = Partial<WorkflowNodeConfigForInspect>;

export type InspectNodeConfigFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: InspectNodeConfigFormFieldValues;
  onValuesChange?: (values: InspectNodeConfigFormFieldValues) => void;
};

export type InspectNodeConfigFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<InspectNodeConfigFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<InspectNodeConfigFormFieldValues>["resetFields"];
  validateFields: FormInstance<InspectNodeConfigFormFieldValues>["validateFields"];
};

const initFormModel = (): InspectNodeConfigFormFieldValues => {
  return {
    domain: "",
    port: "443",
    path: "",
    host: "",
  };
};

const InspectNodeConfigForm = forwardRef<InspectNodeConfigFormInstance, InspectNodeConfigFormProps>(
  ({ className, style, disabled, initialValues, onValuesChange }, ref) => {
    const { t } = useTranslation();

    const formSchema = z.object({
      host: z.string().refine((val) => validIPv4Address(val) || validDomainName(val), {
        message: t("workflow_node.inspect.form.host.placeholder"),
      }),
      domain: z.string().optional(),
      port: z.string().refine((val) => validPortNumber(val), {
        message: t("workflow_node.inspect.form.port.placeholder"),
      }),
      path: z.string().optional(),
    });
    const formRule = createSchemaFieldRule(formSchema);
    const { form: formInst, formProps } = useAntdForm({
      name: "workflowNodeInspectConfigForm",
      initialValues: initialValues ?? initFormModel(),
    });

    const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
      onValuesChange?.(values as InspectNodeConfigFormFieldValues);
    };

    useImperativeHandle(ref, () => {
      return {
        getFieldsValue: () => {
          return formInst.getFieldsValue(true);
        },
        resetFields: (fields) => {
          return formInst.resetFields(fields as (keyof InspectNodeConfigFormFieldValues)[]);
        },
        validateFields: (nameList, config) => {
          return formInst.validateFields(nameList, config);
        },
      } as InspectNodeConfigFormInstance;
    });

    return (
      <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
        <Form.Item name="host" label={t("workflow_node.inspect.form.host.label")} rules={[formRule]}>
          <Input variant="filled" placeholder={t("workflow_node.inspect.form.host.placeholder")} />
        </Form.Item>

        <Form.Item name="port" label={t("workflow_node.inspect.form.port.label")} rules={[formRule]}>
          <Input variant="filled" placeholder={t("workflow_node.inspect.form.port.placeholder")} />
        </Form.Item>

        <Form.Item name="domain" label={t("workflow_node.inspect.form.domain.label")} rules={[formRule]}>
          <Input variant="filled" placeholder={t("workflow_node.inspect.form.domain.placeholder")} />
        </Form.Item>

        <Form.Item name="path" label={t("workflow_node.inspect.form.path.label")} rules={[formRule]}>
          <Input variant="filled" placeholder={t("workflow_node.inspect.form.path.placeholder")} />
        </Form.Item>
      </Form>
    );
  }
);

export default memo(InspectNodeConfigForm);
