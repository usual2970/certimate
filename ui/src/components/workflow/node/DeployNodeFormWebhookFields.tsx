import { useTranslation } from "react-i18next";
import { Alert, Button, Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

const DeployNodeFormWebhookFields = () => {
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
  const formInst = Form.useFormInstance();

  const initialValues: Partial<z.infer<typeof formSchema>> = {
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
    formInst.setFieldValue("webhookData", initialValues.webhookData);
  };

  return (
    <>
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
        <Form.Item name="webhookData" rules={[formRule]} initialValue={initialValues.webhookData}>
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
    </>
  );
};

export default DeployNodeFormWebhookFields;
