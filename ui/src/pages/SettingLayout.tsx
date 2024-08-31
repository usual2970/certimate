import { Toaster } from "@/components/ui/toaster";
import { Outlet } from "react-router-dom";

const SettingLayout = () => {
  return (
    <div>
      <Toaster />
      <div className="text-muted-foreground border-b dark:border-stone-500 py-5">
        设置密码
      </div>
      <div className="w-full sm:w-[35em] mt-10 flex flex-col p-3 mx-auto">
        {/* <div className="text-muted-foreground">
          <span className="transition-all text-sm bg-gray-400 px-3 py-1 rounded-sm text-white cursor-pointer">
            密码
          </span>
        </div> */}
        <Outlet />
      </div>
    </div>
  );
};

export default SettingLayout;
