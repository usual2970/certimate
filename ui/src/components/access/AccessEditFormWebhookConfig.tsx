import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForWebhook } from "@/domain/access";

type AccessEditFormWebhookConfigFieldValues = Partial<AccessConfigForWebhook>;

export type AccessEditFormWebhookConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormWebhookConfigFieldValues;
  onValuesChange?: (values: AccessEditFormWebhookConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormWebhookConfigFieldValues => {
  return {
    url: "",
  };
};

const AccessEditFormWebhookConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormWebhookConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    url: z.string({ message: t("access.form.webhook_url.placeholder") }).url(t("common.errmsg.url_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormWebhookConfigFieldValues);
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues ?? initFormModel()} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item name="url" label={t("access.form.webhook_url.label")} rules={[formRule]}>
        <Input placeholder={t("access.form.webhook_url.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormWebhookConfig;
