import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Form, Input, message, notification } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { getErrMsg } from "@/utils/error";
import { getPocketBase } from "@/repository/pocketbase";

const SettingsPassword = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const formSchema = z
    .object({
      oldPassword: z
        .string({ message: t("settings.password.form.old_password.placeholder") })
        .min(10, { message: t("settings.password.form.password.errmsg.invalid") }),
      newPassword: z
        .string({ message: t("settings.password.form.new_password.placeholder") })
        .min(10, { message: t("settings.password.form.password.errmsg.invalid") }),
      confirmPassword: z
        .string({ message: t("settings.password.form.confirm_password.placeholder") })
        .min(10, { message: t("settings.password.form.password.errmsg.invalid") }),
    })
    .refine((data) => data.newPassword === data.confirmPassword, {
      message: t("settings.password.form.password.errmsg.not_matched"),
      path: ["confirmPassword"],
    });
  const formRule = createSchemaFieldRule(formSchema);
  const [form] = Form.useForm<z.infer<typeof formSchema>>();
  const [formPending, setFormPending] = useState(false);

  const [initialChanged, setInitialChanged] = useState(false);

  const handleInputChange = () => {
    const fields = form.getFieldsValue();
    setInitialChanged(!!fields.oldPassword && !!fields.newPassword && !!fields.confirmPassword);
  };

  const handleFormFinish = async (fields: z.infer<typeof formSchema>) => {
    setFormPending(true);

    try {
      await getPocketBase().admins.authWithPassword(getPocketBase().authStore.model?.email, fields.oldPassword);
      await getPocketBase().admins.update(getPocketBase().authStore.model?.id, {
        password: fields.newPassword,
        passwordConfirm: fields.confirmPassword,
      });

      messageApi.success(t("common.text.operation_succeeded"));

      setTimeout(() => {
        getPocketBase().authStore.clear();
        navigate("/login");
      }, 500);
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    } finally {
      setFormPending(false);
    }
  };

  return (
    <>
      {MessageContextHolder}
      {NotificationContextHolder}

      <div className="md:max-w-[40rem]">
        <Form form={form} disabled={formPending} layout="vertical" onFinish={handleFormFinish}>
          <Form.Item name="oldPassword" label={t("settings.password.form.old_password.label")} rules={[formRule]}>
            <Input.Password placeholder={t("settings.password.form.old_password.placeholder")} onChange={handleInputChange} />
          </Form.Item>

          <Form.Item name="newPassword" label={t("settings.password.form.new_password.label")} rules={[formRule]}>
            <Input.Password placeholder={t("settings.password.form.new_password.placeholder")} onChange={handleInputChange} />
          </Form.Item>

          <Form.Item name="confirmPassword" label={t("settings.password.form.confirm_password.label")} rules={[formRule]}>
            <Input.Password placeholder={t("settings.password.form.confirm_password.placeholder")} onChange={handleInputChange} />
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

export default SettingsPassword;
