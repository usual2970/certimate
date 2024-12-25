import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type WebhookAccessConfig } from "@/domain/access";

type AccessEditFormWebhookConfigModelValues = Partial<WebhookAccessConfig>;

export type AccessEditFormWebhookConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormWebhookConfigModelValues;
  onModelChange?: (model: AccessEditFormWebhookConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormWebhookConfigModelValues => {
  return {
    url: "",
  };
};

const AccessEditFormWebhookConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormWebhookConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    url: z
      .string()
      .min(1, { message: t("access.form.webhook_url.placeholder") })
      .url({ message: t("common.errmsg.url_invalid") }),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onModelChange?.(values as AccessEditFormWebhookConfigModelValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item name="url" label={t("access.form.webhook_url.label")} rules={[formRule]}>
        <Input placeholder={t("access.form.webhook_url.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormWebhookConfig;
