import { Navigate, Outlet } from "react-router-dom";

import Version from "@/components/certimate/Version";
import { getPb } from "@/repository/api";

const LoginLayout = () => {
  if (getPb().authStore.isValid && getPb().authStore.isAdmin) {
    return <Navigate to="/" />;
  }

  return (
    <div className="container">
      <Outlet />

      <Version className="fixed right-0 bottom-0 justify-end pr-5" />
    </div>
  );
};

export default LoginLayout;

