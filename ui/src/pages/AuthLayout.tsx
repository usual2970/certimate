import { Navigate, Outlet } from "react-router-dom";

import Version from "@/components/Version";
import { getAuthStore } from "@/repository/admin";

const AuthLayout = () => {
  const auth = getAuthStore();
  if (auth.isValid && auth.isSuperuser) {
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
