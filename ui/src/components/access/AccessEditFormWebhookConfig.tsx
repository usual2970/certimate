import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type WebhookAccessConfig } from "@/domain/access";

type AccessEditFormWebhookConfigModelType = Partial<WebhookAccessConfig>;

export type AccessEditFormWebhookConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  loading?: boolean;
  model?: AccessEditFormWebhookConfigModelType;
  onModelChange?: (model: AccessEditFormWebhookConfigModelType) => void;
};

const initModel = () => {
  return {
    url: "",
  } as AccessEditFormWebhookConfigModelType;
};

const AccessEditFormWebhookConfig = ({ form, formName, disabled, loading, model, onModelChange }: AccessEditFormWebhookConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    url: z
      .string()
      .min(1, { message: t("access.form.webhook_url.placeholder") })
      .url({ message: t("common.errmsg.url_invalid") }),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormWebhookConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item name="url" label={t("access.form.webhook_url.label")} rules={[formRule]}>
        <Input placeholder={t("access.form.webhook_url.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormWebhookConfig;
