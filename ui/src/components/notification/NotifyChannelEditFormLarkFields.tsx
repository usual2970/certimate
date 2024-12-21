import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

const NotifyChannelEditFormLarkFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    webhookUrl: z
      .string({ message: t("settings.notification.channel.form.lark_webhook_url.placeholder") })
      .min(1, t("settings.notification.channel.form.lark_webhook_url.placeholder"))
      .url({ message: t("common.errmsg.url_invalid") }),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="webhookUrl"
        label={t("settings.notification.channel.form.lark_webhook_url.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.notification.channel.form.lark_webhook_url.tooltip") }}></span>}
      >
        <Input placeholder={t("settings.notification.channel.form.lark_webhook_url.placeholder")} />
      </Form.Item>
    </>
  );
};

export default NotifyChannelEditFormLarkFields;
