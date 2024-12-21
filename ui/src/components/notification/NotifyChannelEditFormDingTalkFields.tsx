import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

const NotifyChannelEditFormDingTalkFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessToken: z
      .string({ message: t("settings.notification.channel.form.dingtalk_access_token.placeholder") })
      .min(1, t("settings.notification.channel.form.dingtalk_access_token.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
    secret: z
      .string({ message: t("settings.notification.channel.form.dingtalk_secret.placeholder") })
      .min(1, t("settings.notification.channel.form.dingtalk_secret.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="accessToken"
        label={t("settings.notification.channel.form.dingtalk_access_token.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.notification.channel.form.dingtalk_access_token.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("settings.notification.channel.form.dingtalk_access_token.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secret"
        label={t("settings.notification.channel.form.dingtalk_secret.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.notification.channel.form.dingtalk_secret.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("settings.notification.channel.form.dingtalk_secret.placeholder")} />
      </Form.Item>
    </>
  );
};

export default NotifyChannelEditFormDingTalkFields;
