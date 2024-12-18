import { createHashRouter } from "react-router-dom";

import AuthLayout from "./pages/AuthLayout";
import ConsoleLayout from "./pages/ConsoleLayout";
import Login from "./pages/login/Login";
import Dashboard from "./pages/dashboard/Dashboard";
import AccessList from "./pages/accesses/AccessList";
import WorkflowList from "./pages/workflows/WorkflowList";
import WorkflowDetail from "./pages/workflows/WorkflowDetail";
import CertificateList from "./pages/certificates/CertificateList";
import Settings from "./pages/settings/Settings";
import SettingsAccount from "./pages/settings/Account";
import SettingsPassword from "./pages/settings/Password";
import SettingsNotification from "./pages/settings/Notification";
import SettingsSSLProvider from "./pages/settings/SSLProvider";

export const router = createHashRouter([
  {
    path: "/",
    element: <ConsoleLayout />,
    children: [
      {
        path: "/",
        element: <Dashboard />,
      },
      {
        path: "/accesses",
        element: <AccessList />,
      },
      {
        path: "/certificates",
        element: <CertificateList />,
      },
      {
        path: "/workflows",
        element: <WorkflowList />,
      },
      {
        path: "/workflows/detail",
        element: <WorkflowDetail />,
      },
      {
        path: "/settings",
        element: <Settings />,
        children: [
          {
            path: "/settings/account",
            element: <SettingsAccount />,
          },
          {
            path: "/settings/password",
            element: <SettingsPassword />,
          },
          {
            path: "/settings/notification",
            element: <SettingsNotification />,
          },
          {
            path: "/settings/ssl-provider",
            element: <SettingsSSLProvider />,
          },
        ],
      },
    ],
  },
  {
    path: "/login",
    element: <AuthLayout />,
    children: [
      {
        path: "/login",
        element: <Login />,
      },
    ],
  },
]);
