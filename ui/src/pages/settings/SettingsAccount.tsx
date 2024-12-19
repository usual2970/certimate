import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Form, Input, message, notification } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { getErrMsg } from "@/utils/error";
import { getPocketBase } from "@/repository/pocketbase";

const SettingsAccount = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const formSchema = z.object({
    username: z.string({ message: "settings.account.form.email.placeholder" }).email({ message: t("common.errmsg.email_invalid") }),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const [form] = Form.useForm<z.infer<typeof formSchema>>();
  const [formPending, setFormPending] = useState(false);

  const [initialValues] = useState<Partial<z.infer<typeof formSchema>>>({
    username: getPocketBase().authStore.model?.email,
  });
  const [initialChanged, setInitialChanged] = useState(false);

  const handleInputChange = () => {
    setInitialChanged(form.getFieldValue("username") !== initialValues.username);
  };

  const handleFormFinish = async (fields: z.infer<typeof formSchema>) => {
    setFormPending(true);

    try {
      await getPocketBase().admins.update(getPocketBase().authStore.model?.id, {
        email: fields.username,
      });

      messageApi.success(t("common.text.operation_succeeded"));

      setTimeout(() => {
        getPocketBase().authStore.clear();
        navigate("/login");
      }, 500);
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: <>{getErrMsg(err)}</> });
    } finally {
      setFormPending(false);
    }
  };

  return (
    <>
      {MessageContextHolder}
      {NotificationContextHolder}

      <div className="w-full md:max-w-[35em]">
        <Form form={form} disabled={formPending} initialValues={initialValues} layout="vertical" onFinish={handleFormFinish}>
          <Form.Item name="username" label={t("settings.account.form.email.label")} rules={[formRule]}>
            <Input placeholder={t("settings.account.form.email.placeholder")} onChange={handleInputChange} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" disabled={!initialChanged} loading={formPending}>
              {t("common.button.save")}
            </Button>
          </Form.Item>
        </Form>
      </div>
    </>
  );
};

export default SettingsAccount;
