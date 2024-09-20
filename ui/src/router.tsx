import { createHashRouter } from "react-router-dom";

import DashboardLayout from "./pages/DashboardLayout";
import Home from "./pages/domains/Home";
import Edit from "./pages/domains/Edit";
import Access from "./pages/access/Access";
import History from "./pages/history/History";
import Login from "./pages/login/Login";
import LoginLayout from "./pages/LoginLayout";
import Password from "./pages/setting/Password";
import SettingLayout from "./pages/SettingLayout";
import Dashboard from "./pages/dashboard/Dashboard";
import Account from "./pages/setting/Account";

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
        path: "/domains",
        element: <Home />,
      },
      {
        path: "/edit",
        element: <Edit />,
      },
      {
        path: "/access",
        element: <Access />,
      },
      {
        path: "/history",
        element: <History />,
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
    path: "/about",
    element: <div>About</div>,
  },
]);
