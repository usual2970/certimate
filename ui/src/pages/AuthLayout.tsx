import { Navigate, Outlet } from "react-router-dom";
import { Layout } from "antd";

import Version from "@/components/Version";
import { getAuthStore } from "@/repository/admin";

const AuthLayout = () => {
  const auth = getAuthStore();
  if (auth.isValid && auth.isSuperuser) {
    return <Navigate to="/" />;
  }

  return (
    <Layout className="h-screen">
      <div className="container">
        <Outlet />

        <Version className="fixed bottom-4 right-8" />
      </div>
    </Layout>
  );
};

export default AuthLayout;
