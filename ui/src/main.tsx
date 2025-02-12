import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import dayjs from "dayjs";
import dayjsUtc from "dayjs/plugin/utc";

import App from "./App";
import "./i18n";
import "./index.css";
import "./global.css";

dayjs.extend(dayjsUtc);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
