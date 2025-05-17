import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

type NotifyNodeConfigFormPushoverConfigFieldValues = Nullish<{
  priority?: string;
}>;

export type NotifyNodeConfigFormPushoverConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: NotifyNodeConfigFormPushoverConfigFieldValues;
  onValuesChange?: (values: NotifyNodeConfigFormPushoverConfigFieldValues) => void;
};

const initFormModel = (): NotifyNodeConfigFormPushoverConfigFieldValues => {
  return {};
};

const NotifyNodeConfigFormPushoverConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: NotifyNodeConfigFormPushoverConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    priority: z.enum(["-2", "-1", "0", "1", "2"]).default("0"),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values);
  };

  return (
    <Form
      form={formInst}
      disabled={disabled}
      initialValues={initialValues ?? initFormModel()}
      layout="vertical"
      name={formName}
      onValuesChange={handleFormChange}
    >
      <Form.Item
        name="priority"
        label={t("workflow_node.notify.form.pushover_priority.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.notify.form.pushover_priority.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.notify.form.pushover_priority.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default NotifyNodeConfigFormPushoverConfig;
