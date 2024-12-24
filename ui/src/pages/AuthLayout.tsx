import { Navigate, Outlet } from "react-router-dom";

import Version from "@/components/ui/Version";
import { getPocketBase } from "@/repository/pocketbase";

const AuthLayout = () => {
  const auth = getPocketBase().authStore;
  if (auth.isValid && auth.isAdmin) {
    return <Navigate to="/" />;
  }

  return (
    <div className="container">
      <Outlet />

      <Version className="fixed right-8 bottom-4" />
    </div>
  );
};

export default AuthLayout;
