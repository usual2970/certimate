import DeployProgress from "@/components/certimate/DeployProgress";
import { Button } from "@/components/ui/button";

import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import { Toaster } from "@/components/ui/toaster";
import { Tooltip, TooltipTrigger } from "@/components/ui/tooltip";
import { Domain } from "@/domain/domain";
import { convertZulu2Beijing, getDate } from "@/lib/time";
import { list, remove, save } from "@/repository/domains";
import { TooltipContent, TooltipProvider } from "@radix-ui/react-tooltip";
import { CircleCheck, CircleX, Earth } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const Home = () => {
  const navigate = useNavigate();
  const handleCreateClick = () => {
    navigate("/edit");
  };

  const handleEditClick = (id: string) => {
    navigate(`/edit?id=${id}`);
  };

  const handleHistoryClick = (id: string) => {
    navigate(`/history?domain=${id}`);
  };

  const handleDeleteClick = async (id: string) => {
    try {
      await remove(id);
      setDomains(domains.filter((domain) => domain.id !== id));
    } catch (error) {
      console.error("Error deleting domain:", error);
    }
  };

  const [domains, setDomains] = useState<Domain[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      const data = await list();
      setDomains(data);
    };
    fetchData();
  }, []);

  const handelCheckedChange = async (id: string) => {
    const checkedDomains = domains.filter((domain) => domain.id === id);
    const isChecked = checkedDomains[0].enabled;

    const data = checkedDomains[0];
    data.enabled = !isChecked;

    await save(data);

    const updatedDomains = domains.map((domain) => {
      if (domain.id === id) {
        return { ...domain, checked: !isChecked };
      }
      return domain;
    });
    setDomains(updatedDomains);
  };

  return (
    <>
      <div className="">
        <Toaster />
        <div className="flex justify-between items-center">
          <div className="text-muted-foreground">域名列表</div>
          <Button onClick={handleCreateClick}>新增域名</Button>
        </div>

        {!domains.length ? (
          <>
            <div className="flex flex-col items-center mt-10">
              <span className="bg-orange-100 p-5 rounded-full">
                <Earth size={40} className="text-primary" />
              </span>

              <div className="text-center text-sm text-muted-foreground mt-3">
                请添加域名开始部署证书吧。
              </div>
              <Button onClick={handleCreateClick} className="mt-3">
                添加域名
              </Button>
            </div>
          </>
        ) : (
          <>
            <div className="hidden sm:flex sm:flex-row text-muted-foreground text-sm border-b sm:p-2 mt-5">
              <div className="w-40">域名</div>
              <div className="w-48">有效期限</div>
              <div className="w-32">最近执行状态</div>
              <div className="w-64">最近执行阶段</div>
              <div className="w-40 sm:ml-2">最近执行时间</div>
              <div className="w-32">是否启用</div>
              <div className="grow">操作</div>
            </div>
            <div className="sm:hidden flex text-sm text-muted-foreground">
              域名
            </div>

            {domains.map((domain) => (
              <div
                className="flex flex-col sm:flex-row text-secondary-foreground border-b sm:p-2 hover:bg-muted/50 text-sm"
                key={domain.id}
              >
                <div className="sm:w-40 w-full pt-1 sm:pt-0 flex items-center">
                  {domain.domain}
                </div>
                <div className="sm:w-48 w-full pt-1 sm:pt-0 flex  items-center">
                  <div>
                    {domain.expiredAt ? (
                      <>
                        <div>有效期90天</div>
                        <div>{getDate(domain.expiredAt)}到期</div>
                      </>
                    ) : (
                      "---"
                    )}
                  </div>
                </div>
                <div className="sm:w-32 w-full pt-1 sm:pt-0 flex items-center">
                  {domain.lastDeployedAt && domain.expand?.lastDeployment ? (
                    <>
                      {domain.expand.lastDeployment?.phase === "deploy" &&
                      domain.expand.lastDeployment?.phaseSuccess ? (
                        <CircleCheck size={16} className="text-green-700" />
                      ) : (
                        <CircleX size={16} className="text-red-700" />
                      )}
                    </>
                  ) : (
                    "---"
                  )}
                </div>
                <div className="sm:w-64 w-full pt-1 sm:pt-0 flex items-center">
                  {domain.lastDeployedAt && domain.expand?.lastDeployment ? (
                    <DeployProgress
                      phase={domain.expand.lastDeployment?.phase}
                      phaseSuccess={domain.expand.lastDeployment?.phaseSuccess}
                    />
                  ) : (
                    "---"
                  )}
                </div>
                <div className="sm:w-40 pt-1 sm:pt-0 sm:ml-2 flex items-center">
                  {domain.lastDeployedAt
                    ? convertZulu2Beijing(domain.lastDeployedAt)
                    : "---"}
                </div>
                <div className="sm:w-32 flex items-center">
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger>
                        <Switch
                          checked={domain.enabled}
                          onCheckedChange={() => {
                            handelCheckedChange(domain.id);
                          }}
                        ></Switch>
                      </TooltipTrigger>
                      <TooltipContent>
                        <div className="border rounded-sm px-3 bg-background text-muted-foreground text-xs">
                          {domain.enabled ? "禁用" : "启用"}
                        </div>
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </div>
                <div className="flex items-center grow justify-start pt-1 sm:pt-0">
                  <Button
                    variant={"link"}
                    className="p-0"
                    onClick={() => handleHistoryClick(domain.id)}
                  >
                    部署历史
                  </Button>
                  {!domain.enabled && (
                    <>
                      <Separator orientation="vertical" className="h-4 mx-2" />
                      <Button
                        variant={"link"}
                        className="p-0"
                        onClick={() => handleDeleteClick(domain.id)}
                      >
                        删除
                      </Button>
                      <Separator orientation="vertical" className="h-4 mx-2" />
                      <Button
                        variant={"link"}
                        className="p-0"
                        onClick={() => handleEditClick(domain.id)}
                      >
                        编辑
                      </Button>
                    </>
                  )}
                </div>
              </div>
            ))}
          </>
        )}
      </div>
    </>
  );
};

export default Home;
