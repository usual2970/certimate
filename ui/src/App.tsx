import { useLayoutEffect, useState } from "react";
import { RouterProvider } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { App as AntdApp, ConfigProvider as AntdConfigProvider } from "antd";
import { type Locale } from "antd/es/locale";
import AntdLocaleEnUs from "antd/locale/en_US";
import AntdLocaleZhCN from "antd/locale/zh_CN";
import dayjs from "dayjs";
import "dayjs/locale/zh-cn";

import { localeNames } from "./i18n";
import { router } from "./router.tsx";
import { ThemeProvider } from "./components/ThemeProvider.tsx";

const App = () => {
  const { i18n } = useTranslation();

  const antdLocalesMap: Record<string, Locale> = {
    [localeNames.ZH]: AntdLocaleZhCN,
    [localeNames.EN]: AntdLocaleEnUs,
  };
  const [antdLocale, setAntdLocale] = useState(antdLocalesMap[i18n.language]);

  const handleLanguageChanged = () => {
    setAntdLocale(antdLocalesMap[i18n.language]);
    dayjs.locale(i18n.language);
  };
  i18n.on("languageChanged", handleLanguageChanged);
  useLayoutEffect(handleLanguageChanged, [i18n]);

  return (
    <AntdConfigProvider
      locale={antdLocale}
      theme={{
        token: {
          colorPrimary: "hsl(24.6 95% 53.1%)",
        },
      }}
    >
      <AntdApp>
        <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
          <RouterProvider router={router} />
        </ThemeProvider>
      </AntdApp>
    </AntdConfigProvider>
  );
};

export default App;
