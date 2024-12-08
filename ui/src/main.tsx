import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { App, ConfigProvider } from "antd";
import AntdLocaleZhCN from "antd/locale/zh_CN";
import "dayjs/locale/zh-cn";

import { router } from "./router.tsx";
import { ThemeProvider } from "./components/ThemeProvider.tsx";
import "./i18n";
import "./global.css";

// TODO: antd i18n
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App>
      <ConfigProvider
        locale={AntdLocaleZhCN}
        theme={{
          token: {
            colorPrimary: "hsl(24.6 95% 53.1%)",
          },
        }}
      >
        <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
          <RouterProvider router={router} />
        </ThemeProvider>
      </ConfigProvider>
    </App>
  </React.StrictMode>
);
