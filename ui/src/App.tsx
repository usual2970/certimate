import { useEffect, useLayoutEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { RouterProvider } from "react-router-dom";
import { App, ConfigProvider, type ThemeConfig, theme } from "antd";
import { type Locale } from "antd/es/locale";
import AntdLocaleEnUs from "antd/locale/en_US";
import AntdLocaleZhCN from "antd/locale/zh_CN";
import dayjs from "dayjs";
import "dayjs/locale/zh-cn";

import { useBrowserTheme } from "@/hooks";
import { localeNames } from "@/i18n";
import { router } from "@/router.tsx";

const RootApp = () => {
  const { i18n } = useTranslation();

  const { theme: browserTheme } = useBrowserTheme();

  const antdLocalesMap: Record<string, Locale> = useMemo(
    () => ({
      [localeNames.ZH]: AntdLocaleZhCN,
      [localeNames.EN]: AntdLocaleEnUs,
    }),
    []
  );
  const [antdLocale, setAntdLocale] = useState(antdLocalesMap[i18n.language]);
  const handleLanguageChanged = () => {
    setAntdLocale(antdLocalesMap[i18n.language]);
    dayjs.locale(i18n.language);
  };
  i18n.on("languageChanged", handleLanguageChanged);
  useLayoutEffect(handleLanguageChanged, [antdLocalesMap, i18n]);

  const antdThemesMap: Record<string, ThemeConfig> = useMemo(
    () => ({
      ["light"]: { algorithm: theme.defaultAlgorithm },
      ["dark"]: { algorithm: theme.darkAlgorithm },
    }),
    []
  );
  const [antdTheme, setAntdTheme] = useState(antdThemesMap[browserTheme]);
  useEffect(() => {
    setAntdTheme(antdThemesMap[browserTheme]);

    const root = window.document.documentElement;
    root.classList.remove("light", "dark");
    root.classList.add(browserTheme);
  }, [antdThemesMap, browserTheme]);

  return (
    <ConfigProvider
      locale={antdLocale}
      theme={{
        ...antdTheme,
        token: {
          /* @see tailwind.config.js */
          colorPrimary: browserTheme === "dark" ? "hsl(20.5 90.2% 48.2%)" : "hsl(24.6 95% 53.1%)",
          colorLink: browserTheme === "dark" ? "hsl(20.5 90.2% 48.2%)" : "hsl(24.6 95% 53.1%)",
        },
      }}
    >
      <App>
        <RouterProvider router={router} />
      </App>
    </ConfigProvider>
  );
};

export default RootApp;
