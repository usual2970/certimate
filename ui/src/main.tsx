import React from "react";
import ReactDOM from "react-dom/client";
import dayjs from "dayjs";
import dayjsUtc from "dayjs/plugin/utc";
import "dayjs/locale/zh-cn";

import App from "./App";
import "./i18n";
import "./global.css";

dayjs.extend(dayjsUtc);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
