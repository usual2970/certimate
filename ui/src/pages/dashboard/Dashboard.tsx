import DeployProgress from "@/components/certimate/DeployProgress";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Deployment, DeploymentListReq, Log } from "@/domain/deployment";
import { Statistic } from "@/domain/domain";
import { convertZulu2Beijing } from "@/lib/time";
import { list } from "@/repository/deployment";
import { statistics } from "@/repository/domains";

import {
  Ban,
  CalendarX2,
  CircleCheck,
  CircleX,
  LoaderPinwheel,
  Smile,
  SquareSigma,
} from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

const Dashboard = () => {
  const [statistic, setStatistic] = useState<Statistic>();
  const [deployments, setDeployments] = useState<Deployment[]>();

  const navigate = useNavigate();

  useEffect(() => {
    const fetchStatistic = async () => {
      const data = await statistics();
      setStatistic(data);
    };

    fetchStatistic();
  }, []);

  useEffect(() => {
    const fetchData = async () => {
      const param: DeploymentListReq = {
        perPage: 8,
      };

      const data = await list(param);
      setDeployments(data.items);
    };
    fetchData();
  }, []);

  return (
    <div className="flex flex-col">
      <div className="flex justify-between items-center">
        <div className="text-muted-foreground">控制面板</div>
      </div>
      <div className="flex mt-10 gap-5 flex-col flex-wrap md:flex-row">
        <div className="w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg border">
          <div className="p-3">
            <SquareSigma size={48} strokeWidth={1} className="text-blue-400" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">所有</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.total ? (
                  <Link to="/domains" className="hover:underline">
                    {statistic?.total}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">个</div>
            </div>
          </div>
        </div>

        <div className="w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg border">
          <div className="p-3">
            <CalendarX2 size={48} strokeWidth={1} className="text-red-400" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">即将过期</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.expired ? (
                  <Link to="/domains?state=expired" className="hover:underline">
                    {statistic?.expired}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">个</div>
            </div>
          </div>
        </div>

        <div className="border w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg">
          <div className="p-3">
            <LoaderPinwheel
              size={48}
              strokeWidth={1}
              className="text-green-400"
            />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">启用中</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.enabled ? (
                  <Link to="/domains?state=enabled" className="hover:underline">
                    {statistic?.enabled}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">个</div>
            </div>
          </div>
        </div>

        <div className="border w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg">
          <div className="p-3">
            <Ban size={48} strokeWidth={1} className="text-gray-400" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">未启用</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.disabled ? (
                  <Link
                    to="/domains?state=disabled"
                    className="hover:underline"
                  >
                    {statistic?.disabled}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">个</div>
            </div>
          </div>
        </div>
      </div>

      <div>
        <div className="text-muted-foreground mt-5 text-sm">部署历史</div>

        {deployments?.length == 0 ? (
          <>
            <Alert className="max-w-[40em] mt-10">
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
      </div>
    </div>
  );
};

export default Dashboard;
