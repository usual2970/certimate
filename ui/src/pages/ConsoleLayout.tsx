import { memo, useEffect, useState } from "react";
import { Link, Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Drawer, Dropdown, Layout, Menu, Tooltip, theme, type ButtonProps, type MenuProps } from "antd";
import {
  Languages as LanguagesIcon,
  LogOut as LogOutIcon,
  Home as HomeIcon,
  Menu as MenuIcon,
  Moon as MoonIcon,
  Server as ServerIcon,
  Settings as SettingsIcon,
  ShieldCheck as ShieldCheckIcon,
  Sun as SunIcon,
  Workflow as WorkflowIcon,
} from "lucide-react";

import Version from "@/components/certimate/Version";
import { useTheme } from "@/hooks";
import { getPocketBase } from "@/repository/pocketbase";

const ConsoleLayout = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [siderOpen, setSiderOpen] = useState(false);

  const handleSiderOpen = () => {
    setSiderOpen(true);
  };

  const handleSiderClose = () => {
    setSiderOpen(false);
  };

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
    <Layout className="w-full min-h-screen">
      <Layout.Sider className="max-md:hidden" theme="light" width={256}>
        <div className="flex flex-col items-center justify-between w-full h-full overflow-hidden">
          <div className="w-full">
            <SiderMenu />
          </div>
          <div className="w-full py-2 text-center">
            <Version />
          </div>
        </div>
      </Layout.Sider>

      <Layout>
        <Layout.Header style={{ padding: 0, background: themeToken.colorBgContainer }}>
          <div className="flex items-center justify-between size-full px-4 overflow-hidden">
            <div className="flex items-center gap-4 size-full">
              <Button className="md:hidden" icon={<MenuIcon />} size="large" onClick={handleSiderOpen} />
              <Drawer
                closable={false}
                destroyOnClose
                open={siderOpen}
                placement="left"
                styles={{
                  content: { paddingTop: themeToken.paddingSM, paddingBottom: themeToken.paddingSM },
                  body: { padding: 0 },
                }}
                onClose={handleSiderClose}
              >
                <SiderMenu onSelect={() => handleSiderClose()} />
              </Drawer>
            </div>
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
  );
};

const SiderMenu = memo(({ onSelect }: { onSelect?: (key: string) => void }) => {
  const location = useLocation();
  const navigate = useNavigate();

  const { t } = useTranslation();

  const MENU_KEY_HOME = "/";
  const MENU_KEY_WORKFLOWS = "/workflows";
  const MENU_KEY_CERTIFICATES = "/certificates";
  const MENU_KEY_ACCESSES = "/accesses";
  const menuItems: Required<MenuProps>["items"] = [
    [MENU_KEY_HOME, <HomeIcon size={16} />, t("dashboard.page.title")],
    [MENU_KEY_WORKFLOWS, <WorkflowIcon size={16} />, t("workflow.page.title")],
    [MENU_KEY_CERTIFICATES, <ShieldCheckIcon size={16} />, t("certificate.page.title")],
    [MENU_KEY_ACCESSES, <ServerIcon size={16} />, t("access.page.title")],
  ].map(([key, icon, label]) => {
    return {
      key: key as string,
      icon: icon,
      label: label,
      onClick: () => {
        navigate(key as string);
        onSelect?.(key as string);
      },
    };
  });
  const [menuSelectedKey, setMenuSelectedKey] = useState<string>();

  const getActiveMenuItem = () => {
    const item =
      menuItems.find((item) => item!.key === location.pathname) ??
      menuItems.find((item) => item!.key !== MENU_KEY_HOME && location.pathname.startsWith(item!.key as string));
    return item;
  };

  useEffect(() => {
    const item = getActiveMenuItem();
    if (item) {
      setMenuSelectedKey(item.key as string);
    } else {
      setMenuSelectedKey(undefined);
    }
  }, [location.pathname]);

  useEffect(() => {
    if (menuSelectedKey && menuSelectedKey !== getActiveMenuItem()?.key) {
      navigate(menuSelectedKey);
    }
  }, [menuSelectedKey]);

  return (
    <>
      <Link to="/" className="flex items-center gap-2 w-full px-4 font-semibold overflow-hidden">
        <img src="/logo.svg" className="w-[36px] h-[36px]" />
        <span className="w-[74px] h-[64px] leading-[64px] dark:text-white truncate">Certimate</span>
      </Link>
      <div className="flex-grow w-full overflow-x-hidden overflow-y-auto">
        <Menu
          items={menuItems}
          mode="vertical"
          selectedKeys={menuSelectedKey ? [menuSelectedKey] : []}
          onSelect={({ key }) => {
            setMenuSelectedKey(key);
          }}
        />
      </div>
    </>
  );
});

const ThemeToggleButton = memo(({ size }: { size?: ButtonProps["size"] }) => {
  const { t } = useTranslation();

  const { theme, setThemeMode } = useTheme();

  const items: Required<MenuProps>["items"] = [
    ["light", t("common.theme.light")],
    ["dark", t("common.theme.dark")],
    ["system", t("common.theme.system")],
  ].map(([key, label]) => {
    return {
      key: key as string,
      label: label,
      onClick: () => {
        setThemeMode(key as Parameters<typeof setThemeMode>[0]);
        window.location.reload();
      },
    };
  });

  return (
    <Dropdown menu={{ items }} trigger={["click"]}>
      <Button icon={theme === "dark" ? <MoonIcon size={18} /> : <SunIcon size={18} />} size={size} />
    </Dropdown>
  );
});

const LocaleToggleButton = memo(({ size }: { size?: ButtonProps["size"] }) => {
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
});

export default ConsoleLayout;
