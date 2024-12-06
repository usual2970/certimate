import { createHashRouter } from "react-router-dom";

import LoginLayout from "./pages/LoginLayout";
import DashboardLayout from "./pages/DashboardLayout";
import SettingLayout from "./pages/SettingLayout";
import Login from "./pages/login/Login";
import Account from "./pages/setting/Account";
import Password from "./pages/setting/Password";
import Notify from "./pages/setting/Notify";
import SSLProvider from "./pages/setting/SSLProvider";
import Dashboard from "./pages/dashboard/Dashboard";
import AccessList from "./pages/accesses/AccessList";
import WorkflowList from "./pages/workflows/WorkflowList";
import WorkflowDetail from "./pages/workflows/WorkflowDetail";
import CertificateList from "./pages/certificates/CertificateList";

export const router = createHashRouter([
  {
    path: "/",
    element: <DashboardLayout />,
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
        path: "/setting",
        element: <SettingLayout />,
        children: [
          {
            path: "/setting/password",
            element: <Password />,
          },
          {
            path: "/setting/account",
            element: <Account />,
          },
          {
            path: "/setting/notify",
            element: <Notify />,
          },
          {
            path: "/setting/ssl-provider",
            element: <SSLProvider />,
          },
        ],
      },
    ],
  },
  {
    path: "/login",
    element: <LoginLayout />,
    children: [
      {
        path: "/login",
        element: <Login />,
      },
    ],
  },
  {
    path: "/workflow/detail",
    element: <WorkflowDetail />,
  },
]);
