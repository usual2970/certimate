import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
import {
  CloudServerOutlined as CloudServerOutlinedIcon,
  GlobalOutlined as GlobalOutlinedIcon,
  HomeOutlined as HomeOutlinedIcon,
  LogoutOutlined as LogoutOutlinedIcon,
  MenuOutlined as MenuOutlinedIcon,
  MoonOutlined as MoonOutlinedIcon,
  NodeIndexOutlined as NodeIndexOutlinedIcon,
  SafetyOutlined as SafetyOutlinedIcon,
  SettingOutlined as SettingOutlinedIcon,
  SunOutlined as SunOutlinedIcon,
} from "@ant-design/icons";
import { Alert, Button, type ButtonProps, Drawer, Dropdown, Layout, Menu, type MenuProps, Tooltip, theme } from "antd";

import Show from "@/components/Show";
import Version from "@/components/Version";
import { useBrowserTheme, useTriggerElement } from "@/hooks";
import { getAuthStore } from "@/repository/admin";
import { isBrowserHappy } from "@/utils/browser";

const ConsoleLayout = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const handleLogoutClick = () => {
    auth.clear();
    navigate("/login");
  };

  const auth = getAuthStore();
  if (!auth.isValid || !auth.isSuperuser) {
    return <Navigate to="/login" />;
  }

  return (
    <Layout className="h-screen" hasSider>
      <Layout.Sider className="fixed left-0 top-0 z-20 h-full max-md:static max-md:hidden" width="256px" theme="light">
        <div className="flex size-full flex-col items-center justify-between overflow-hidden">
          <div className="w-full">
            <SiderMenu />
          </div>
          <div className="w-full py-2 text-center">
            <Version />
          </div>
        </div>
      </Layout.Sider>

      <Layout className="flex flex-col overflow-hidden pl-[256px] max-md:pl-0">
        <Show when={!isBrowserHappy()}>
          <Alert message={t("common.text.happy_browser")} type="warning" showIcon closable />
        </Show>

        <Layout.Header className="p-0 shadow-sm" style={{ background: themeToken.colorBgContainer }}>
          <div className="flex size-full items-center justify-between overflow-hidden px-4">
            <div className="flex items-center gap-4">
              <SiderMenuDrawer trigger={<Button className="md:hidden" icon={<MenuOutlinedIcon />} size="large" />} />
            </div>
            <div className="flex size-full grow items-center justify-end gap-4 overflow-hidden">
              <Tooltip title={t("common.menu.theme")} mouseEnterDelay={2}>
                <ThemeToggleButton size="large" />
              </Tooltip>
              <Tooltip title={t("common.menu.locale")} mouseEnterDelay={2}>
                <LocaleToggleButton size="large" />
              </Tooltip>
              <Tooltip title={t("common.menu.logout")} mouseEnterDelay={2}>
                <Button danger icon={<LogoutOutlinedIcon />} size="large" onClick={handleLogoutClick} />
              </Tooltip>
            </div>
          </div>
        </Layout.Header>

        <Layout.Content className="flex-1 overflow-y-auto overflow-x-hidden">
          <Outlet />
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
  const MENU_KEY_SETTINGS = "/settings";
  const menuItems: Required<MenuProps>["items"] = [
    [MENU_KEY_HOME, <HomeOutlinedIcon />, t("dashboard.page.title")],
    [MENU_KEY_WORKFLOWS, <NodeIndexOutlinedIcon />, t("workflow.page.title")],
    [MENU_KEY_CERTIFICATES, <SafetyOutlinedIcon />, t("certificate.page.title")],
    [MENU_KEY_ACCESSES, <CloudServerOutlinedIcon />, t("access.page.title")],
    [MENU_KEY_SETTINGS, <SettingOutlinedIcon />, t("settings.page.title")],
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
      <div className="flex w-full items-center gap-2 overflow-hidden px-4 font-semibold">
        <img src="/logo.svg" className="size-[36px]" />
        <span className="h-[64px] w-[74px] truncate leading-[64px] dark:text-white">Certimate</span>
      </div>
      <div className="w-full grow overflow-y-auto overflow-x-hidden">
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

const SiderMenuDrawer = memo(({ trigger }: { trigger: React.ReactNode }) => {
  const { token: themeToken } = theme.useToken();

  const [siderOpen, setSiderOpen] = useState(false);

  const triggerEl = useTriggerElement(trigger, { onClick: () => setSiderOpen(true) });

  return (
    <>
      {triggerEl}

      <Drawer
        closable={false}
        destroyOnHidden
        open={siderOpen}
        placement="left"
        styles={{
          content: { paddingTop: themeToken.paddingSM, paddingBottom: themeToken.paddingSM },
          body: { padding: 0 },
        }}
        onClose={() => setSiderOpen(false)}
      >
        <SiderMenu onSelect={() => setSiderOpen(false)} />
      </Drawer>
    </>
  );
});

const ThemeToggleButton = memo(({ size }: { size?: ButtonProps["size"] }) => {
  const { t } = useTranslation();

  const { theme, themeMode, setThemeMode } = useBrowserTheme();

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
        if (key !== themeMode) {
          window.location.reload();
        }
      },
    };
  });

  return (
    <Dropdown menu={{ items }} trigger={["click"]}>
      <Button icon={theme === "dark" ? <MoonOutlinedIcon /> : <SunOutlinedIcon />} size={size} />
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
      <Button icon={<GlobalOutlinedIcon />} size={size} />
    </Dropdown>
  );
});

export default ConsoleLayout;
