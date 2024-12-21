import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

const NotifyChannelEditFormWebhookFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    url: z
      .string({ message: t("settings.notification.channel.form.webhook_url.placeholder") })
      .min(1, t("settings.notification.channel.form.webhook_url.placeholder"))
      .url({ message: t("common.errmsg.url_invalid") }),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <div>
      <Form.Item name="url" label={t("settings.notification.channel.form.webhook_url.label")} rules={[formRule]}>
        <Input placeholder={t("settings.notification.channel.form.webhook_url.placeholder")} />
      </Form.Item>
    </div>
  );
};

export default NotifyChannelEditFormWebhookFields;
