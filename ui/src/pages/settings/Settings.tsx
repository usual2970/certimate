import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import {
  ApiOutlined as ApiOutlinedIcon,
  DatabaseOutlined as DatabaseOutlinedIcon,
  LockOutlined as LockOutlinedIcon,
  SendOutlined as SendOutlinedIcon,
  UserOutlined as UserOutlinedIcon,
} from "@ant-design/icons";
import { PageHeader } from "@ant-design/pro-components";
import { Card, Space } from "antd";

const Settings = () => {
  const location = useLocation();
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [tabValue, setTabValue] = useState("account");
  useEffect(() => {
    const path = location.pathname.split("/")[2];
    if (!path) {
      navigate("/settings/account");
      return;
    }

    setTabValue(path);
  }, [location.pathname]);

  return (
    <div className="p-4">
      <PageHeader title={t("settings.page.title")} />

      <Card
        tabList={[
          {
            key: "account",
            label: (
              <Space>
                <UserOutlinedIcon />
                <label>{t("settings.account.tab")}</label>
              </Space>
            ),
          },
          {
            key: "password",
            label: (
              <Space>
                <LockOutlinedIcon />
                <label>{t("settings.password.tab")}</label>
              </Space>
            ),
          },
          {
            key: "notification",
            label: (
              <Space>
                <SendOutlinedIcon />
                <label>{t("settings.notification.tab")}</label>
              </Space>
            ),
          },
          {
            key: "ssl-provider",
            label: (
              <Space>
                <ApiOutlinedIcon />
                <label>{t("settings.sslprovider.tab")}</label>
              </Space>
            ),
          },
          {
            key: "persistence",
            label: (
              <Space>
                <DatabaseOutlinedIcon />
                <label>{t("settings.persistence.tab")}</label>
              </Space>
            ),
          },
        ]}
        activeTabKey={tabValue}
        onTabChange={(key) => {
          setTabValue(key);
          navigate(`/settings/${key}`);
        }}
      >
        <Outlet />
      </Card>
    </div>
  );
};

export default Settings;
