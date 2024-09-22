import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from "@/components/ui/navigation-menu";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Toaster } from "@/components/ui/toaster";
import { KeyRound, Megaphone, UserRound } from "lucide-react";
import { useEffect, useState } from "react";

import { Outlet, useLocation, useNavigate } from "react-router-dom";

const SettingLayout = () => {
  const location = useLocation();
  const [tabValue, setTabValue] = useState("account");
  const navigate = useNavigate();

  useEffect(() => {
    const pathname = location.pathname;
    const tabValue = pathname.split("/")[2];
    setTabValue(tabValue);
  }, [location]);

  return (
    <div>
      <Toaster />
      <div className="text-muted-foreground border-b dark:border-stone-500 py-5">
        偏好设置
      </div>
      <div className="w-full mt-5 p-0 md:p-3 flex justify-center">
        <Tabs defaultValue="account" className="w-full" value={tabValue}>
          <TabsList className="mx-auto">
            <TabsTrigger
              value="account"
              onClick={() => {
                navigate("/setting/account");
              }}
              className="px-5"
            >
              <UserRound size={14} />
              <div className="ml-1">账户</div>
            </TabsTrigger>
            <TabsTrigger
              value="password"
              onClick={() => {
                navigate("/setting/password");
              }}
              className="px-5"
            >
              <KeyRound size={14} />
              <div className="ml-1">密码</div>
            </TabsTrigger>

            <TabsTrigger
              value="notify"
              onClick={() => {
                navigate("/setting/notify");
              }}
              className="px-5"
            >
              <Megaphone size={14} />
              <div className="ml-1">消息推送</div>
            </TabsTrigger>
          </TabsList>
          <TabsContent value={tabValue}>
            <div className="mt-5 w-full md:w-[45em]">
              <Outlet />
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default SettingLayout;
