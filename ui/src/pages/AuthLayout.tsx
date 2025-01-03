import { Navigate, Outlet } from "react-router-dom";

import Version from "@/components/Version";
import { getPocketBase } from "@/repository/pocketbase";

const AuthLayout = () => {
  const auth = getPocketBase().authStore;
  if (auth.isValid && auth.isAdmin) {
    return <Navigate to="/" />;
  }

  return (
    <div className="container">
      <Outlet />

      <Version className="fixed bottom-4 right-8" />
    </div>
  );
};

export default AuthLayout;
