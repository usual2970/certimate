import { useEffect, useState } from "react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Card, Space } from "antd";
import { PageHeader } from "@ant-design/pro-components";
import { KeyRound as KeyRoundIcon, Megaphone as MegaphoneIcon, ShieldCheck as ShieldCheckIcon, UserRound as UserRoundIcon } from "lucide-react";

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
  }, [location]);

  return (
    <>
      <PageHeader title={t("settings.page.title")} />

      <Card
        tabList={[
          {
            key: "account",
            label: (
              <Space>
                <UserRoundIcon size={14} />
                <label>{t("settings.account.tab")}</label>
              </Space>
            ),
          },
          {
            key: "password",
            label: (
              <Space>
                <KeyRoundIcon size={14} />
                <label>{t("settings.password.tab")}</label>
              </Space>
            ),
          },
          {
            key: "notification",
            label: (
              <Space>
                <MegaphoneIcon size={14} />
                <label>{t("settings.notification.tab")}</label>
              </Space>
            ),
          },
          {
            key: "ssl-provider",
            label: (
              <Space>
                <ShieldCheckIcon size={14} />
                <label>{t("settings.ca.tab")}</label>
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
    </>
  );
};

export default Settings;
