import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Card, Form, Input, notification } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { getPocketBase } from "@/repository/pocketbase";

const Login = () => {
  const navigage = useNavigate();

  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const formSchema = z.object({
    username: z.string().email(t("login.username.errmsg.invalid")),
    password: z.string().min(10, t("login.password.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const [form] = Form.useForm();

  const [isPending, setIsPending] = useState(false);

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setIsPending(true);

    try {
      await getPocketBase().admins.authWithPassword(values.username, values.password);
      navigage("/");
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
    } finally {
      setIsPending(false);
    }
  };

  return (
    <>
      {NotificationContextHolder}

      <Card className="mx-auto mt-32 p-10 max-w-[35em] border dark:border-stone-500 rounded-md shadow-md">
        <div className="flex items-center justify-center mb-10">
          <img src="/logo.svg" className="w-16" />
        </div>

        <Form form={form} disabled={isPending} layout="vertical" onFinish={onSubmit}>
          <Form.Item name="username" label={t("login.username.label")} rules={[formRule]}>
            <Input placeholder={t("login.username.placeholder")} />
          </Form.Item>

          <Form.Item name="password" label={t("login.password.label")} rules={[formRule]}>
            <Input type="password" placeholder={t("login.password.placeholder")} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={isPending}>
              {t("login.submit")}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </>
  );
};

export default Login;
