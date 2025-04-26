import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input, Select, Switch } from "antd";
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
    method: "POST",
    headers: "Content-Type: application/json",
    allowInsecureConnections: false,
  };
};

const AccessFormWebhookConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormWebhookConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    url: z.string().url(t("common.errmsg.url_invalid")),
    method: z.union([z.literal("GET"), z.literal("POST"), z.literal("PUT"), z.literal("PATCH"), z.literal("DELETE")], {
      message: t("access.form.webhook_method.placeholder"),
    }),
    headers: z
      .string()
      .nullish()
      .refine((v) => {
        if (!v) return true;

        const lines = v.split(/\r?\n/);
        for (const line of lines) {
          if (line.split(":").length < 2) {
            return false;
          }
        }
        return true;
      }, t("access.form.webhook_headers.errmsg.invalid")),
    allowInsecureConnections: z.boolean().nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleWebhookHeadersBlur = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    let value = e.target.value;
    value = value.trim();
    value = value.replace(/(?<!\r)\n/g, "\r\n");
    formInst.setFieldValue("headers", value);
  };

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

      <Form.Item name="method" label={t("access.form.webhook_method.label")} rules={[formRule]}>
        <Select
          options={["GET", "POST", "PUT", "PATCH", "DELETE"].map((s) => ({ label: s, value: s }))}
          placeholder={t("access.form.webhook_method.placeholder")}
        />
      </Form.Item>

      <Form.Item
        name="headers"
        label={t("access.form.webhook_headers.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.webhook_headers.tooltip") }}></span>}
      >
        <Input.TextArea autoSize={{ minRows: 3, maxRows: 5 }} placeholder={t("access.form.webhook_headers.placeholder")} onBlur={handleWebhookHeadersBlur} />
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
