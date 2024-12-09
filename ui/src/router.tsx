import { createHashRouter } from "react-router-dom";

import AuthLayout from "./pages/AuthLayout";
import ConsoleLayout from "./pages/ConsoleLayout";
import SettingsLayout from "./pages/SettingsLayout";
import Login from "./pages/login/Login";
import Account from "./pages/settings/Account";
import Password from "./pages/settings/Password";
import Notify from "./pages/settings/Notify";
import SSLProvider from "./pages/settings/SSLProvider";
import Dashboard from "./pages/dashboard/Dashboard";
import AccessList from "./pages/accesses/AccessList";
import WorkflowList from "./pages/workflows/WorkflowList";
import WorkflowDetail from "./pages/workflows/WorkflowDetail";
import CertificateList from "./pages/certificates/CertificateList";

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
        element: <SettingsLayout />,
        children: [
          {
            path: "/settings/password",
            element: <Password />,
          },
          {
            path: "/settings/account",
            element: <Account />,
          },
          {
            path: "/settings/notify",
            element: <Notify />,
          },
          {
            path: "/settings/ssl-provider",
            element: <SSLProvider />,
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
