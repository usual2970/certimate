import { useTranslation } from "react-i18next";
import { Navigate, Outlet } from "react-router-dom";
import { Alert, Layout } from "antd";

import Show from "@/components/Show";
import Version from "@/components/Version";
import { getAuthStore } from "@/repository/admin";
import { isBrowserHappy } from "@/utils/browser";

const AuthLayout = () => {
  const { t } = useTranslation();

  const auth = getAuthStore();
  if (auth.isValid && auth.isSuperuser) {
    return <Navigate to="/" />;
  }

  return (
    <Layout className="h-screen">
      <Show when={!isBrowserHappy()}>
        <Alert message={t("common.text.happy_browser")} type="warning" showIcon closable />
      </Show>

      <div className="container">
        <Outlet />

        <Version className="fixed bottom-4 right-8" />
      </div>
    </Layout>
  );
};

export default AuthLayout;
