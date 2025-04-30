import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

type NotifyNodeConfigFormTelegramConfigFieldValues = Nullish<{
  chatId?: string | number;
}>;

export type NotifyNodeConfigFormTelegramConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: NotifyNodeConfigFormTelegramConfigFieldValues;
  onValuesChange?: (values: NotifyNodeConfigFormTelegramConfigFieldValues) => void;
};

const initFormModel = (): NotifyNodeConfigFormTelegramConfigFieldValues => {
  return {};
};

const NotifyNodeConfigFormTelegramConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: NotifyNodeConfigFormTelegramConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    chatId: z
      .preprocess(
        (v) => (v == null || v === "" ? undefined : Number(v)),
        z
          .number()
          .nullish()
          .refine((v) => {
            if (v == null || v + "" === "") return true;
            return /^\d+$/.test(v + "") && +v! > 0;
          }, t("workflow_node.notify.form.telegram_chat_id.placeholder"))
      )
      .nullish(),
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
        name="chatId"
        label={t("workflow_node.notify.form.telegram_chat_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.notify.form.telegram_chat_id.tooltip") }}></span>}
      >
        <Input type="number" allowClear placeholder={t("workflow_node.notify.form.telegram_chat_id.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default NotifyNodeConfigFormTelegramConfig;
