import { useTranslation } from "react-i18next";
import { Alert, Button, Form, type FormInstance, Input, Select, Switch } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import Show from "@/components/Show";
import { type AccessConfigForWebhook } from "@/domain/access";

type AccessFormWebhookConfigFieldValues = Nullish<AccessConfigForWebhook>;

export type AccessFormWebhookConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormWebhookConfigFieldValues;
  usage?: "deployment" | "notification" | "none";
  onValuesChange?: (values: AccessFormWebhookConfigFieldValues) => void;
};

const initFormModel = (): AccessFormWebhookConfigFieldValues => {
  return {
    url: "",
    method: "POST",
    headers: "Content-Type: application/json",
    allowInsecureConnections: false,
    defaultDataForDeployment: JSON.stringify(
      {
        name: "${DOMAINS}",
        cert: "${CERTIFICATE}",
        privkey: "${PRIVATE_KEY}",
      },
      null,
      2
    ),
    defaultDataForNotification: JSON.stringify(
      {
        subject: "${SUBJECT}",
        message: "${MESSAGE}",
      },
      null,
      2
    ),
  };
};

const AccessFormWebhookConfig = ({ form: formInst, formName, disabled, initialValues, usage, onValuesChange }: AccessFormWebhookConfigProps) => {
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
    defaultDataForDeployment: z
      .string()
      .nullish()
      .refine((v) => {
        if (usage && usage !== "deployment") return true;
        if (!v) return true;

        try {
          const obj = JSON.parse(v);
          return typeof obj === "object" && !Array.isArray(obj);
        } catch {
          return false;
        }
      }, t("access.form.webhook_default_data.errmsg.json_invalid")),
    defaultDataForNotification: z
      .string()
      .nullish()
      .refine((v) => {
        if (usage && usage !== "notification") return true;
        if (!v) return true;

        try {
          const obj = JSON.parse(v);
          return typeof obj === "object" && !Array.isArray(obj);
        } catch {
          return false;
        }
      }, t("access.form.webhook_default_data.errmsg.json_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleWebhookHeadersBlur = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    let value = e.target.value;
    value = value.trim();
    value = value.replace(/(?<!\r)\n/g, "\r\n");
    formInst.setFieldValue("headers", value);
  };

  const handleWebhookDataForDeploymentBlur = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    try {
      const json = JSON.stringify(JSON.parse(value), null, 2);
      formInst.setFieldValue("defaultDataForDeployment", json);
    } catch {
      return;
    }
  };

  const handleWebhookDataForNotificationBlur = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    try {
      const json = JSON.stringify(JSON.parse(value), null, 2);
      formInst.setFieldValue("defaultDataForNotification", json);
    } catch {
      return;
    }
  };

  const handlePresetDataForDeploymentClick = () => {
    formInst.setFieldValue("defaultDataForDeployment", initFormModel().defaultDataForDeployment);
  };

  const handlePresetDataForNotificationClick = () => {
    formInst.setFieldValue("defaultDataForNotification", initFormModel().defaultDataForNotification);
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

      <Show when={!usage || usage === "deployment"}>
        <Form.Item className="mb-0">
          <label className="mb-1 block">
            <div className="flex w-full items-center justify-between gap-4">
              <div className="max-w-full grow truncate">{t("access.form.webhook_default_data_for_deployment.label")}</div>
              <div className="text-right">
                <Button size="small" type="link" onClick={handlePresetDataForDeploymentClick}>
                  {t("access.form.webhook_default_data_preset.button")}
                </Button>
              </div>
            </div>
          </label>
          <Form.Item name="defaultDataForDeployment" rules={[formRule]}>
            <Input.TextArea
              allowClear
              autoSize={{ minRows: 3, maxRows: 10 }}
              placeholder={t("access.form.webhook_default_data_for_deployment.placeholder")}
              onBlur={handleWebhookDataForDeploymentBlur}
            />
          </Form.Item>
        </Form.Item>

        <Form.Item>
          <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("access.form.webhook_default_data_for_deployment.guide") }}></span>} />
        </Form.Item>
      </Show>

      <Show when={!usage || usage === "notification"}>
        <Form.Item className="mb-0">
          <label className="mb-1 block">
            <div className="flex w-full items-center justify-between gap-4">
              <div className="max-w-full grow truncate">{t("access.form.webhook_default_data_for_notification.label")}</div>
              <div className="text-right">
                <Button size="small" type="link" onClick={handlePresetDataForNotificationClick}>
                  {t("access.form.webhook_default_data_preset.button")}
                </Button>
              </div>
            </div>
          </label>
          <Form.Item name="defaultDataForNotification" rules={[formRule]}>
            <Input.TextArea
              allowClear
              autoSize={{ minRows: 3, maxRows: 10 }}
              placeholder={t("access.form.webhook_default_data_for_notification.placeholder")}
              onBlur={handleWebhookDataForNotificationBlur}
            />
          </Form.Item>
        </Form.Item>

        <Form.Item>
          <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("access.form.webhook_default_data_for_notification.guide") }}></span>} />
        </Form.Item>
      </Show>

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
