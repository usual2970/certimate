import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

const NotifyChannelEditFormServerChanFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    url: z
      .string({ message: t("settings.notification.channel.form.serverchan_url.placeholder") })
      .min(1, t("settings.notification.channel.form.serverchan_url.placeholder"))
      .url({ message: t("common.errmsg.url_invalid") }),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="url"
        label={t("settings.notification.channel.form.serverchan_url.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.notification.channel.form.serverchan_url.tooltip") }}></span>}
      >
        <Input placeholder={t("settings.notification.channel.form.serverchan_url.placeholder")} />
      </Form.Item>
    </>
  );
};

export default NotifyChannelEditFormServerChanFields;