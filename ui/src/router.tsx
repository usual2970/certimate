import { createHashRouter } from "react-router-dom";

import DashboardLayout from "./pages/DashboardLayout";
import Home from "./pages/domains/Home";
import Edit from "./pages/domains/Edit";
import Access from "./pages/access/Access";
import History from "./pages/history/History";
import Login from "./pages/login/Login";
import LoginLayout from "./pages/LoginLayout";

export const router = createHashRouter([
  {
    path: "/",
    element: <DashboardLayout />,
    children: [
      {
        path: "/",
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
