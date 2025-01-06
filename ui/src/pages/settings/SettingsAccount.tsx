import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { Button, Form, Input, message, notification } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { getPocketBase } from "@/repository/pocketbase";
import { getErrMsg } from "@/utils/error";

const SettingsAccount = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const formSchema = z.object({
    username: z.string({ message: "settings.account.form.email.placeholder" }).email({ message: t("common.errmsg.email_invalid") }),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const {
    form: formInst,
    formPending,
    formProps,
  } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: {
      username: getPocketBase().authStore.model?.email,
    },
    onSubmit: async (values) => {
      try {
        await getPocketBase().admins.update(getPocketBase().authStore.model?.id, {
          email: values.username,
        });

        messageApi.success(t("common.text.operation_succeeded"));

        setTimeout(() => {
          getPocketBase().authStore.clear();
          navigate("/login");
        }, 500);
      } catch (err) {
        notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

        throw err;
      }
    },
  });
  const [formChanged, setFormChanged] = useState(false);

  const handleInputChange = () => {
    setFormChanged(formInst.getFieldValue("username") !== formProps.initialValues?.username);
  };

  return (
    <>
      {MessageContextHolder}
      {NotificationContextHolder}

      <div className="md:max-w-[40rem]">
        <Form {...formProps} form={formInst} disabled={formPending} layout="vertical">
          <Form.Item name="username" label={t("settings.account.form.email.label")} rules={[formRule]}>
            <Input placeholder={t("settings.account.form.email.placeholder")} onChange={handleInputChange} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" disabled={!formChanged} loading={formPending}>
              {t("common.button.save")}
            </Button>
          </Form.Item>
        </Form>
      </div>
    </>
  );
};

export default SettingsAccount;
