import { useTranslation } from "react-i18next";
import { Alert, Button, Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

type DeployNodeConfigFormWebhookConfigFieldValues = Nullish<{
  webhookData: string;
}>;

export type DeployNodeConfigFormWebhookConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormWebhookConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormWebhookConfigFieldValues) => void;
};

const initFormModel = (): DeployNodeConfigFormWebhookConfigFieldValues => {
  return {
    webhookData: JSON.stringify(
      {
        name: "${DOMAINS}",
        cert: "${CERTIFICATE}",
        privkey: "${PRIVATE_KEY}",
      },
      null,
      2
    ),
  };
};

const DeployNodeConfigFormWebhookConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: DeployNodeConfigFormWebhookConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    webhookData: z.string({ message: t("workflow_node.deploy.form.webhook_data.placeholder") }).refine((v) => {
      try {
        JSON.parse(v);
        return true;
      } catch {
        return false;
      }
    }, t("workflow_node.deploy.form.webhook_data.errmsg.json_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleWebhookDataBlur = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    try {
      const json = JSON.stringify(JSON.parse(value), null, 2);
      formInst.setFieldValue("webhookData", json);
    } catch {
      return;
    }
  };

  const handlePresetDataClick = () => {
    formInst.setFieldValue("webhookData", initFormModel().webhookData);
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
      <Form.Item className="mb-0">
        <label className="mb-1 block">
          <div className="flex w-full items-center justify-between gap-4">
            <div className="max-w-full grow truncate">{t("workflow_node.deploy.form.webhook_data.label")}</div>
            <div className="text-right">
              <Button size="small" type="link" onClick={handlePresetDataClick}>
                {t("workflow_node.deploy.form.webhook_data_preset.button")}
              </Button>
            </div>
          </div>
        </label>
        <Form.Item name="webhookData" rules={[formRule]}>
          <Input.TextArea
            autoSize={{ minRows: 3, maxRows: 10 }}
            placeholder={t("workflow_node.deploy.form.webhook_data.placeholder")}
            onBlur={handleWebhookDataBlur}
          />
        </Form.Item>
      </Form.Item>

      <Form.Item>
        <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.webhook_data.guide") }}></span>} />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigFormWebhookConfig;
