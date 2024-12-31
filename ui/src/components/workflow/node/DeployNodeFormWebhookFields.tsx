import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
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
    }, t("workflow_node.deploy.form.webhook_data.placeholder")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const formInst = Form.useFormInstance();

  const handleWebhookDataBlur = (e: React.FocusEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    try {
      const json = JSON.stringify(JSON.parse(value), null, 2);
      formInst.setFieldValue("webhookData", json);
    } catch {
      return;
    }
  };

  return (
    <>
      <Form.Item
        name="webhookData"
        label={t("workflow_node.deploy.form.webhook_data.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.webhook_data.tooltip") }}></span>}
      >
        <Input.TextArea
          autoSize={{ minRows: 3, maxRows: 10 }}
          placeholder={t("workflow_node.deploy.form.webhook_data.placeholder")}
          onBlur={handleWebhookDataBlur}
        />
      </Form.Item>
    </>
  );
};

export default DeployNodeFormWebhookFields;
