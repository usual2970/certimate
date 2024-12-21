import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useRequest } from "ahooks";
import { Button, Form, Input, message, notification, Skeleton } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";
import { ClientResponseError } from "pocketbase";

import { defaultNotifyTemplate, SETTINGS_NAMES, type NotifyTemplatesSettingsContent } from "@/domain/settings";
import { get as getSettings, save as saveSettings } from "@/repository/settings";
import { getErrMsg } from "@/utils/error";

export type NotifyTemplateFormProps = {
  className?: string;
  style?: React.CSSProperties;
};

const NotifyTemplateForm = ({ className, style }: NotifyTemplateFormProps) => {
  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const formSchema = z.object({
    subject: z
      .string()
      .trim()
      .min(1, t("settings.notification.template.form.subject.placeholder"))
      .max(1000, t("common.errmsg.string_max", { max: 1000 })),
    message: z
      .string()
      .trim()
      .min(1, t("settings.notification.template.form.message.placeholder"))
      .max(1000, t("common.errmsg.string_max", { max: 1000 })),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const [form] = Form.useForm<z.infer<typeof formSchema>>();
  const [formPending, setFormPending] = useState(false);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>();
  const [initialChanged, setInitialChanged] = useState(false);

  const { loading } = useRequest(
    () => {
      return getSettings<NotifyTemplatesSettingsContent>(SETTINGS_NAMES.NOTIFY_TEMPLATES);
    },
    {
      onError: (err) => {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);
      },
      onFinally: (_, resp) => {
        const template = resp?.content?.notifyTemplates?.[0] ?? defaultNotifyTemplate;
        setInitialValues({ ...template });
      },
    }
  );

  const handleInputChange = () => {
    setInitialChanged(true);
  };

  const handleFormFinish = async (fields: z.infer<typeof formSchema>) => {
    setFormPending(true);

    try {
      const settings = await getSettings<NotifyTemplatesSettingsContent>(SETTINGS_NAMES.NOTIFY_TEMPLATES);
      await saveSettings<NotifyTemplatesSettingsContent>({
        ...settings,
        content: {
          notifyTemplates: [fields],
        },
      });

      messageApi.success(t("common.text.operation_succeeded"));
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    } finally {
      setFormPending(false);
    }
  };

  return (
    <div className={className} style={style}>
      {MessageContextHolder}
      {NotificationContextHolder}

      {loading ? (
        <Skeleton active />
      ) : (
        <Form form={form} disabled={formPending} initialValues={initialValues} layout="vertical" onFinish={handleFormFinish}>
          <Form.Item
            name="subject"
            label={t("settings.notification.template.form.subject.label")}
            extra={t("settings.notification.template.form.subject.tooltip")}
            rules={[formRule]}
          >
            <Input placeholder={t("settings.notification.template.form.subject.placeholder")} onChange={handleInputChange} />
          </Form.Item>

          <Form.Item
            name="message"
            label={t("settings.notification.template.form.message.label")}
            extra={t("settings.notification.template.form.message.tooltip")}
            rules={[formRule]}
          >
            <Input.TextArea
              autoSize={{ minRows: 3, maxRows: 5 }}
              placeholder={t("settings.notification.template.form.message.placeholder")}
              onChange={handleInputChange}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" disabled={!initialChanged} loading={formPending}>
              {t("common.button.save")}
            </Button>
          </Form.Item>
        </Form>
      )}
    </div>
  );
};

export default NotifyTemplateForm;
