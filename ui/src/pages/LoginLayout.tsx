import { getPb } from "@/repository/api";
import { Navigate, Outlet } from "react-router-dom";

const LoginLayout = () => {
  if (getPb().authStore.isValid && getPb().authStore.isAdmin) {
    return <Navigate to="/" />;
  }
  return (
    <div className="container">
      <Outlet />
    </div>
  );
};

export default LoginLayout;
