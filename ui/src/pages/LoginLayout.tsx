import { Navigate, Outlet } from "react-router-dom";

import Version from "@/components/certimate/Version";
import { getPocketBase } from "@/repository/pocketbase";

const LoginLayout = () => {
  if (getPocketBase().authStore.isValid && getPocketBase().authStore.isAdmin) {
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
