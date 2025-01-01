import { memo, useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Link, Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
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
import { Button, Drawer, Dropdown, Layout, Menu, Tooltip, theme, type ButtonProps, type MenuProps } from "antd";

import Version from "@/components/core/Version";
import { useBrowserTheme } from "@/hooks";
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
    <Layout className="min-h-screen" hasSider>
      <Layout.Sider className="max-md:hidden max-md:static fixed top-0 left-0 h-full z-[20]" width="256px" theme="light">
        <div className="flex flex-col items-center justify-between w-full h-full overflow-hidden">
          <div className="w-full">
            <SiderMenu />
          </div>
          <div className="w-full py-2 text-center">
            <Version />
          </div>
        </div>
      </Layout.Sider>

      <Layout className="pl-[256px] max-md:pl-0">
        <Layout.Header className="sticky top-0 left-0 right-0 p-0 z-[19] shadow-sm" style={{ background: themeToken.colorBgContainer }}>
          <div className="flex items-center justify-between size-full px-4 overflow-hidden">
            <div className="flex items-center gap-4 size-full">
              <Button className="md:hidden" icon={<MenuOutlinedIcon />} size="large" onClick={handleSiderOpen} />
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
                <Button icon={<SettingOutlinedIcon />} size="large" onClick={handleSettingsClick} />
              </Tooltip>
              <Tooltip title={t("common.menu.logout")} mouseEnterDelay={2}>
                <Button danger icon={<LogoutOutlinedIcon />} size="large" onClick={handleLogoutClick} />
              </Tooltip>
            </div>
          </div>
        </Layout.Header>

        <Layout.Content style={{ overflow: "initial" }}>
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
  const menuItems: Required<MenuProps>["items"] = [
    [MENU_KEY_HOME, <HomeOutlinedIcon />, t("dashboard.page.title")],
    [MENU_KEY_WORKFLOWS, <NodeIndexOutlinedIcon />, t("workflow.page.title")],
    [MENU_KEY_CERTIFICATES, <SafetyOutlinedIcon />, t("certificate.page.title")],
    [MENU_KEY_ACCESSES, <CloudServerOutlinedIcon />, t("access.page.title")],
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

  const getActiveMenuItem = useCallback(() => {
    const item =
      menuItems.find((item) => item!.key === location.pathname) ??
      menuItems.find((item) => item!.key !== MENU_KEY_HOME && location.pathname.startsWith(item!.key as string));
    return item;
  }, [location.pathname, menuItems]);

  useEffect(() => {
    const item = getActiveMenuItem();
    if (item) {
      setMenuSelectedKey(item.key as string);
    } else {
      setMenuSelectedKey(undefined);
    }
  }, [location.pathname, getActiveMenuItem]);

  useEffect(() => {
    if (menuSelectedKey && menuSelectedKey !== getActiveMenuItem()?.key) {
      navigate(menuSelectedKey);
    }
  }, [menuSelectedKey, navigate, getActiveMenuItem]);

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
