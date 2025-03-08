import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input, Switch } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForWebhook } from "@/domain/access";

type AccessFormWebhookConfigFieldValues = Nullish<AccessConfigForWebhook>;

export type AccessFormWebhookConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormWebhookConfigFieldValues;
  onValuesChange?: (values: AccessFormWebhookConfigFieldValues) => void;
};

const initFormModel = (): AccessFormWebhookConfigFieldValues => {
  return {
    url: "",
  };
};

const AccessFormWebhookConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormWebhookConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    url: z.string({ message: t("access.form.webhook_url.placeholder") }).url(t("common.errmsg.url_invalid")),
    allowInsecureConnections: z.boolean().nullish(),
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
      <Form.Item name="url" label={t("access.form.webhook_url.label")} rules={[formRule]}>
        <Input placeholder={t("access.form.webhook_url.placeholder")} />
      </Form.Item>

      <Form.Item
        name="allowInsecureConnections"
        label={t("access.form.webhook_allow_insecure_conns.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.webhook_allow_insecure_conns.tooltip") }}></span>}
      >
        <Switch
          checkedChildren={t("access.form.webhook_allow_insecure_conns.switch.on")}
          unCheckedChildren={t("access.form.webhook_allow_insecure_conns.switch.off")}
        />
      </Form.Item>
    </Form>
  );
};

export default AccessFormWebhookConfig;
