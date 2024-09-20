import DeployProgress from "@/components/certimate/DeployProgress";
import DeployState from "@/components/certimate/DeployState";
import XPagination from "@/components/certimate/XPagination";
import Show from "@/components/Show";
import {
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialog,
  AlertDialogContent,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";

import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import { Toaster } from "@/components/ui/toaster";
import { Tooltip, TooltipTrigger } from "@/components/ui/tooltip";
import { useToast } from "@/components/ui/use-toast";
import { Domain } from "@/domain/domain";
import { CustomFile, saveFiles2ZIP } from "@/lib/file";
import { convertZulu2Beijing, getDate } from "@/lib/time";
import {
  list,
  remove,
  save,
  subscribeId,
  unsubscribeId,
} from "@/repository/domains";

import { TooltipContent, TooltipProvider } from "@radix-ui/react-tooltip";
import { CircleCheck, CircleX, Earth } from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";

const Home = () => {
  const toast = useToast();

  const navigate = useNavigate();

  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const page = query.get("page");

  const state = query.get("state");

  const [totalPage, setTotalPage] = useState(0);

  const handleCreateClick = () => {
    navigate("/edit");
  };

  const setPage = (newPage: number) => {
    query.set("page", newPage.toString());
    navigate(`?${query.toString()}`);
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
      const data = await list({
        page: page ? Number(page) : 1,
        perPage: 10,
        state: state ? state : "",
      });

      setDomains(data.items);
      setTotalPage(data.totalPages);
    };
    fetchData();
  }, [page, state]);

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

  const handleRightNowClick = async (domain: Domain) => {
    try {
      unsubscribeId(domain.id);
      subscribeId(domain.id, (resp) => {
        console.log(resp);
        const updatedDomains = domains.map((domain) => {
          if (domain.id === resp.id) {
            return { ...resp };
          }
          return domain;
        });
        setDomains(updatedDomains);
      });
      domain.rightnow = true;

      await save(domain);

      toast.toast({
        title: "操作成功",
        description: "已发起部署，请稍后查看部署日志。",
      });
    } catch (e) {
      toast.toast({
        title: "执行失败",
        description: (
          <>
            执行失败，请查看
            <Link
              to={`/history?domain=${domain.id}`}
              className="underline text-blue-500"
            >
              部署日志
            </Link>
            查看详情。
          </>
        ),
        variant: "destructive",
      });
    }
  };

  const handleForceClick = async (domain: Domain) => {
    await handleRightNowClick({ ...domain, deployed: false });
  };

  const handleDownloadClick = async (domain: Domain) => {
    const zipName = `${domain.id}-${domain.domain}.zip`;
    const files: CustomFile[] = [
      {
        name: `${domain.domain}.pem`,
        content: domain.certificate ? domain.certificate : "",
      },
      {
        name: `${domain.domain}.key`,
        content: domain.privateKey ? domain.privateKey : "",
      },
    ];

    await saveFiles2ZIP(zipName, files);
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
            <div className="hidden sm:flex sm:flex-row text-muted-foreground text-sm border-b dark:border-stone-500 sm:p-2 mt-5">
              <div className="w-36">域名</div>
              <div className="w-40">有效期限</div>
              <div className="w-32">最近执行状态</div>
              <div className="w-64">最近执行阶段</div>
              <div className="w-40 sm:ml-2">最近执行时间</div>
              <div className="w-24">是否启用</div>
              <div className="grow">操作</div>
            </div>
            <div className="sm:hidden flex text-sm text-muted-foreground">
              域名
            </div>

            {domains.map((domain) => (
              <div
                className="flex flex-col sm:flex-row text-secondary-foreground border-b  dark:border-stone-500 sm:p-2 hover:bg-muted/50 text-sm"
                key={domain.id}
              >
                <div className="sm:w-36 w-full pt-1 sm:pt-0 flex items-center">
                  {domain.domain}
                </div>
                <div className="sm:w-40 w-full pt-1 sm:pt-0 flex  items-center">
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
                      <DeployState deployment={domain.expand.lastDeployment} />
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
                <div className="sm:w-24 flex items-center">
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
                  <Show when={domain.enabled ? true : false}>
                    <Separator orientation="vertical" className="h-4 mx-2" />
                    <Button
                      variant={"link"}
                      className="p-0"
                      onClick={() => handleRightNowClick(domain)}
                    >
                      立即部署
                    </Button>
                  </Show>

                  <Show
                    when={
                      (domain.enabled ? true : false) && domain.deployed
                        ? true
                        : false
                    }
                  >
                    <Separator orientation="vertical" className="h-4 mx-2" />
                    <Button
                      variant={"link"}
                      className="p-0"
                      onClick={() => handleForceClick(domain)}
                    >
                      强行部署
                    </Button>
                  </Show>

                  <Show when={domain.expiredAt ? true : false}>
                    <Separator orientation="vertical" className="h-4 mx-2" />
                    <Button
                      variant={"link"}
                      className="p-0"
                      onClick={() => handleDownloadClick(domain)}
                    >
                      下载
                    </Button>
                  </Show>

                  {!domain.enabled && (
                    <>
                      <Separator orientation="vertical" className="h-4 mx-2" />
                      <AlertDialog>
                        <AlertDialogTrigger asChild>
                          <Button variant={"link"} className="p-0">
                            删除
                          </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>删除域名</AlertDialogTitle>
                            <AlertDialogDescription>
                              确定要删除域名吗？
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>取消</AlertDialogCancel>
                            <AlertDialogAction
                              onClick={() => {
                                handleDeleteClick(domain.id);
                              }}
                            >
                              确认
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>

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

            <XPagination
              totalPages={totalPage}
              currentPage={page ? Number(page) : 1}
              onPageChange={(page) => {
                setPage(page);
              }}
            />
          </>
        )}
      </div>
    </>
  );
};

export default Home;
