import { Link, Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { CircleUser, Earth, History, Home, Menu, Server } from "lucide-react";

import LocaleToggle from "@/components/LocaleToggle";
import { ThemeToggle } from "@/components/ThemeToggle";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { cn } from "@/lib/utils";
import { getPb } from "@/repository/api";
import { ConfigProvider } from "@/providers/config";

import Version from "@/components/certimate/Version";

export default function Dashboard() {
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();

  if (!getPb().authStore.isValid || !getPb().authStore.isAdmin) {
    return <Navigate to="/login" />;
  }

  const currentPath = location.pathname;
  const getClass = (path: string) => {
    if (path == currentPath) {
      return "bg-muted text-primary";
    }
    return "text-muted-foreground";
  };

  const handleLogoutClick = () => {
    getPb().authStore.clear();
    navigate("/login");
  };

  const handleSettingClick = () => {
    navigate("/setting/account");
  };

  return (
    <>
      <ConfigProvider>
        <div className="grid min-h-screen w-full md:grid-cols-[180px_1fr] lg:grid-cols-[200px_1fr] 2xl:md:grid-cols-[280px_1fr] ">
          <div className="hidden border-r dark:border-stone-500 bg-muted/40 md:block">
            <div className="flex h-full max-h-screen flex-col gap-2">
              <div className="flex h-14 items-center border-b dark:border-stone-500 px-4 lg:h-[60px] lg:px-6">
                <Link to="/" className="flex items-center gap-2 font-semibold">
                  <img src="/vite.svg" className="w-[36px] h-[36px]" />
                  <span className="dark:text-white">Certimate</span>
                </Link>
              </div>
              <div className="flex-1">
                <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
                  <Link to="/" className={cn("flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary", getClass("/"))}>
                    <Home className="h-4 w-4" />
                    {t("dashboard.page.title")}
                  </Link>
                  <Link to="/domains" className={cn("flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary", getClass("/domains"))}>
                    <Earth className="h-4 w-4" />
                    {t("domain.page.title")}
                  </Link>
                  <Link to="/access" className={cn("flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary", getClass("/access"))}>
                    <Server className="h-4 w-4" />
                    {t("access.page.title")}
                  </Link>

                  <Link to="/history" className={cn("flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary", getClass("/history"))}>
                    <History className="h-4 w-4" />
                    {t("history.page.title")}
                  </Link>
                </nav>
              </div>

              <div className="">
                <Version className="justify-center" />
              </div>
            </div>
          </div>
          <div className="flex flex-col">
            <header className="flex h-14 items-center gap-4 border-b dark:border-stone-500 bg-muted/40 px-4 lg:h-[60px] lg:px-6">
              <Sheet>
                <SheetTrigger asChild>
                  <Button variant="outline" size="icon" className="shrink-0 md:hidden">
                    <Menu className="h-5 w-5 dark:text-white" />
                    <span className="sr-only">Toggle navigation menu</span>
                  </Button>
                </SheetTrigger>
                <SheetContent side="left" className="flex flex-col">
                  <nav className="grid gap-2 text-lg font-medium">
                    <Link to="/" className="flex items-center gap-2 text-lg font-semibold">
                      <img src="/vite.svg" className="w-[36px] h-[36px]" />
                      <span className="dark:text-white">Certimate</span>
                      <span className="sr-only">Certimate</span>
                    </Link>
                    <Link to="/" className={cn("mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2  hover:text-foreground", getClass("/"))}>
                      <Home className="h-5 w-5" />
                      {t("dashboard.page.title")}
                    </Link>
                    <Link
                      to="/domains"
                      className={cn("mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2  hover:text-foreground", getClass("/domains"))}
                    >
                      <Earth className="h-5 w-5" />
                      {t("domain.page.title")}
                    </Link>
                    <Link
                      to="/access"
                      className={cn("mx-[-0.65rem] flex items-center gap-4 rounded-xl  px-3 py-2  hover:text-foreground", getClass("/access"))}
                    >
                      <Server className="h-5 w-5" />
                      {t("access.page.title")}
                    </Link>

                    <Link
                      to="/history"
                      className={cn("mx-[-0.65rem] flex items-center gap-4 rounded-xl  px-3 py-2  hover:text-foreground", getClass("/history"))}
                    >
                      <History className="h-5 w-5" />
                      {t("history.page.title")}
                    </Link>
                  </nav>
                  <div className="">
                    <Version className="justify-center" />
                  </div>
                </SheetContent>
              </Sheet>
              <div className="w-full flex-1"></div>
              <ThemeToggle />
              <LocaleToggle />
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="secondary" size="icon" className="rounded-full">
                    <CircleUser className="h-5 w-5" />
                    <span className="sr-only">Toggle user menu</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={handleSettingClick}>{t("common.menu.settings")}</DropdownMenuItem>
                  <DropdownMenuItem onClick={handleLogoutClick}>{t("common.menu.logout")}</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </header>
            <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6 relative">
              <Outlet />
            </main>
          </div>
        </div>
      </ConfigProvider>
    </>
  );
}

