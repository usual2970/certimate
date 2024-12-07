import { useEffect, useState } from "react";
import { Link, Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Dropdown, Layout, Menu, Tooltip, theme, type ButtonProps, type MenuProps } from "antd";
import {
  Languages as LanguagesIcon,
  LogOut as LogOutIcon,
  Home as HomeIcon,
  Menu as MenuIcon,
  Server as ServerIcon,
  Settings as SettingsIcon,
  ShieldCheck as ShieldCheckIcon,
  Sun as SunIcon,
  Workflow as WorkflowIcon,
} from "lucide-react";

import Version from "@/components/certimate/Version";
import { getPocketBase } from "@/repository/pocketbase";
import { ConfigProvider } from "@/providers/config";

const ConsoleLayout = () => {
  const location = useLocation();
  const navigate = useNavigate();

  const { t } = useTranslation();

  const {
    token: { colorBgContainer },
  } = theme.useToken();

  const menuItems: Required<MenuProps>["items"] = [
    {
      key: "/",
      icon: <HomeIcon size={16} />,
      label: t("dashboard.page.title"),
      onClick: () => navigate("/"),
    },
    {
      key: "/workflows",
      icon: <WorkflowIcon size={16} />,
      label: t("workflow.page.title"),
      onClick: () => navigate("/workflows"),
    },
    {
      key: "/certificates",
      icon: <ShieldCheckIcon size={16} />,
      label: t("certificate.page.title"),
      onClick: () => navigate("/certificates"),
    },
    {
      key: "/accesses",
      icon: <ServerIcon size={16} />,
      label: t("access.page.title"),
      onClick: () => navigate("/accesses"),
    },
  ];
  const [menuSelectedKey, setMenuSelectedKey] = useState<string>();

  useEffect(() => {
    const item =
      menuItems.find((item) => item!.key === location.pathname) ??
      menuItems.find((item) => item!.key !== "/" && location.pathname.startsWith(item!.key as string));
    console.log(item);
    if (item) {
      setMenuSelectedKey(item.key as string);
    } else {
      setMenuSelectedKey(undefined);
    }
  }, [location.pathname]);

  useEffect(() => {
    if (menuSelectedKey) {
      navigate(menuSelectedKey);
    }
  }, [menuSelectedKey]);

  // TODO: 响应式侧边栏菜单

  const handleLogoutClick = () => {
    auth.clear();
    navigate("/login");
  };

  const handleSettingsClick = () => {
    navigate("/settings/account");
  };

  const auth = getPocketBase().authStore;
  if (!auth.isValid || !auth.isAdmin) {
    return <Navigate to="/login" />;
  }

  return (
    <>
      <ConfigProvider>
        <Layout className="w-full min-h-screen">
          <Layout.Sider theme="light" width={256}>
            <div className="flex flex-col items-center justify-between w-full h-full overflow-hidden">
              <Link to="/" className="flex items-center gap-2 w-full px-4 font-semibold overflow-hidden">
                <img src="/logo.svg" className="w-[36px] h-[36px]" />
                <span className="w-[64px] h-[64px] leading-[64px] dark:text-white truncate">Certimate</span>
              </Link>
              <div className="flex-grow w-full overflow-x-hidden overflow-y-auto">
                <Menu
                  mode="vertical"
                  items={menuItems}
                  selectedKeys={menuSelectedKey ? [menuSelectedKey] : []}
                  onSelect={({ key }) => {
                    setMenuSelectedKey(key);
                  }}
                />
              </div>
              <div className="w-full py-2 text-center">
                <Version />
              </div>
            </div>
          </Layout.Sider>

          <Layout>
            <Layout.Header style={{ padding: 0, background: colorBgContainer }}>
              <div className="flex items-center justify-between size-full px-4 overflow-hidden">
                <div className="flex items-center gap-4 size-full">{/* <Button icon={<MenuIcon />} size="large" /> */}</div>
                <div className="flex-grow flex items-center justify-end gap-4 size-full overflow-hidden">
                  <Tooltip title={t("common.menu.theme")} mouseEnterDelay={2}>
                    <ThemeToggleButton size="large" />
                  </Tooltip>
                  <Tooltip title={t("common.menu.locale")} mouseEnterDelay={2}>
                    <LocaleToggleButton size="large" />
                  </Tooltip>
                  <Tooltip title={t("common.menu.settings")} mouseEnterDelay={2}>
                    <Button icon={<SettingsIcon size={18} />} size="large" onClick={handleSettingsClick} />
                  </Tooltip>
                  <Tooltip title={t("common.menu.logout")} mouseEnterDelay={2}>
                    <Button danger icon={<LogOutIcon size={18} />} size="large" onClick={handleLogoutClick} />
                  </Tooltip>
                </div>
              </div>
            </Layout.Header>

            <Layout.Content>
              <div className="p-4">
                <Outlet />
              </div>
            </Layout.Content>
          </Layout>
        </Layout>
      </ConfigProvider>
    </>
  );
};

const ThemeToggleButton = ({ size }: { size?: ButtonProps["size"] }) => {
  // TODO: 主题切换
  const items: Required<MenuProps>["items"] = [];

  return (
    <Dropdown menu={{ items }} trigger={["click"]}>
      <Button icon={<SunIcon size={18} />} size={size} onClick={() => alert("TODO")} />
    </Dropdown>
  );
};

const LocaleToggleButton = ({ size }: { size?: ButtonProps["size"] }) => {
  const { i18n } = useTranslation();

  const items: Required<MenuProps>["items"] = Object.keys(i18n.store.data).map((key) => {
    return {
      key: key,
      label: <>{i18n.store.data[key].name as string}</>,
      onClick: () => i18n.changeLanguage(key),
    };
  });

  return (
    <Dropdown menu={{ items }} trigger={["click"]}>
      <Button icon={<LanguagesIcon size={18} />} size={size} />
    </Dropdown>
  );
};

export default ConsoleLayout;
