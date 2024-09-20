import DeployProgress from "@/components/certimate/DeployProgress";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Deployment, DeploymentListReq, Log } from "@/domain/deployment";
import { convertZulu2Beijing } from "@/lib/time";
import { list } from "@/repository/deployment";
import { CircleCheck, CircleX, Smile } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

const History = () => {
  const navigate = useNavigate();
  const [deployments, setDeployments] = useState<Deployment[]>();
  const [searchParams] = useSearchParams();
  const domain = searchParams.get("domain");

  useEffect(() => {
    const fetchData = async () => {
      const param: DeploymentListReq = {};
      if (domain) {
        param.domain = domain;
      }
      const data = await list(param);
      setDeployments(data.items);
    };
    fetchData();
  }, [domain]);

  return (
    <ScrollArea className="h-[80vh] overflow-hidden">
      <div className="text-muted-foreground">部署历史</div>
      {!deployments?.length ? (
        <>
          <Alert className="max-w-[40em] mx-auto mt-20">
            <AlertTitle>暂无数据</AlertTitle>
            <AlertDescription>
              <div className="flex items-center mt-5">
                <div>
                  <Smile className="text-yellow-400" size={36} />
                </div>
                <div className="ml-2">
                  {" "}
                  你暂未创建任何部署，请先添加域名进行部署吧！
                </div>
              </div>
              <div className="mt-2 flex justify-end">
                <Button
                  onClick={() => {
                    navigate("/");
                  }}
                >
                  添加域名
                </Button>
              </div>
            </AlertDescription>
          </Alert>
        </>
      ) : (
        <>
          <div className="hidden sm:flex sm:flex-row text-muted-foreground text-sm border-b dark:border-stone-500 sm:p-2 mt-5">
            <div className="w-48">域名</div>

            <div className="w-24">状态</div>
            <div className="w-56">阶段</div>
            <div className="w-56 sm:ml-2 text-center">最近执行时间</div>

            <div className="grow">操作</div>
          </div>
          <div className="sm:hidden flex text-sm text-muted-foreground">
            部署历史
          </div>

          {deployments?.map((deployment) => (
            <div
              key={deployment.id}
              className="flex flex-col sm:flex-row text-secondary-foreground border-b  dark:border-stone-500 sm:p-2 hover:bg-muted/50 text-sm"
            >
              <div className="sm:w-48 w-full pt-1 sm:pt-0 flex items-center">
                {deployment.expand.domain?.domain}
              </div>
              <div className="sm:w-24 w-full pt-1 sm:pt-0 flex items-center">
                {deployment.phase === "deploy" && deployment.phaseSuccess ? (
                  <CircleCheck size={16} className="text-green-700" />
                ) : (
                  <CircleX size={16} className="text-red-700" />
                )}
              </div>
              <div className="sm:w-56 w-full pt-1 sm:pt-0 flex items-center">
                <DeployProgress
                  phase={deployment.phase}
                  phaseSuccess={deployment.phaseSuccess}
                />
              </div>
              <div className="sm:w-56 w-full pt-1 sm:pt-0 flex items-center sm:justify-center">
                {convertZulu2Beijing(deployment.deployedAt)}
              </div>
              <div className="flex items-center grow justify-start pt-1 sm:pt-0 sm:ml-2">
                <Sheet>
                  <SheetTrigger asChild>
                    <Button variant={"link"} className="p-0">
                      日志
                    </Button>
                  </SheetTrigger>
                  <SheetContent className="sm:max-w-5xl">
                    <SheetHeader>
                      <SheetTitle>
                        {deployment.expand.domain?.domain}-{deployment.id}
                        部署详情
                      </SheetTitle>
                    </SheetHeader>
                    <div className="bg-gray-950 text-stone-100 p-5 text-sm h-[80dvh]">
                      {deployment.log.check && (
                        <>
                          {deployment.log.check.map((item: Log) => {
                            return (
                              <div className="flex flex-col mt-2">
                                <div className="flex">
                                  <div>[{item.time}]</div>
                                  <div className="ml-2">{item.message}</div>
                                </div>
                                {item.error && (
                                  <div className="mt-1 text-red-600">
                                    {item.error}
                                  </div>
                                )}
                              </div>
                            );
                          })}
                        </>
                      )}

                      {deployment.log.apply && (
                        <>
                          {deployment.log.apply.map((item: Log) => {
                            return (
                              <div className="flex flex-col mt-2">
                                <div className="flex">
                                  <div>[{item.time}]</div>
                                  <div className="ml-2">{item.message}</div>
                                </div>
                                {item.info &&
                                  item.info.map((info: string) => {
                                    return (
                                      <div className="mt-1 text-green-600">
                                        {info}
                                      </div>
                                    );
                                  })}
                                {item.error && (
                                  <div className="mt-1 text-red-600">
                                    {item.error}
                                  </div>
                                )}
                              </div>
                            );
                          })}
                        </>
                      )}

                      {deployment.log.deploy && (
                        <>
                          {deployment.log.deploy.map((item: Log) => {
                            return (
                              <div className="flex flex-col mt-2">
                                <div className="flex">
                                  <div>[{item.time}]</div>
                                  <div className="ml-2">{item.message}</div>
                                </div>
                                {item.error && (
                                  <div className="mt-1 text-red-600">
                                    {item.error}
                                  </div>
                                )}
                              </div>
                            );
                          })}
                        </>
                      )}
                    </div>
                  </SheetContent>
                </Sheet>
              </div>
            </div>
          ))}
        </>
      )}
    </ScrollArea>
  );
};

export default History;
