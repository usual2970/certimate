import AccessGroupEdit from "@/components/certimate/AccessGroupEdit";
import Show from "@/components/Show";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { getProviderInfo } from "@/domain/access";
import { getErrMessage } from "@/lib/error";
import { useConfig } from "@/providers/config";
import { remove } from "@/repository/access_group";
import { Group } from "lucide-react";
import { useToast } from "@/components/ui/use-toast";
import { Toaster } from "@/components/ui/toaster";
import { useNavigate } from "react-router-dom";
const AccessGroups = () => {
  const {
    config: { accessGroups },
    reloadAccessGroups,
  } = useConfig();

  const { toast } = useToast();

  const navigate = useNavigate();

  const handleRemoveClick = async (id: string) => {
    try {
      await remove(id);
      reloadAccessGroups();
    } catch (e) {
      toast({
        title: "删除失败",
        description: getErrMessage(e),
        variant: "destructive",
      });
      return;
    }
  };

  const handleAddAccess = () => {
    navigate("/access");
  };
  return (
    <div>
      <Toaster />
      <div className="flex justify-between items-center">
        <div className="text-muted-foreground">部署授权组</div>

        <AccessGroupEdit trigger={<Button>新增授权组</Button>} />
      </div>

      <div className="mt-10">
        <Show when={accessGroups.length == 0}>
          <>
            <div className="flex flex-col items-center mt-10">
              <span className="bg-orange-100 p-5 rounded-full">
                <Group size={40} className="text-primary" />
              </span>

              <div className="text-center text-sm text-muted-foreground mt-3">
                请添加域名开始部署证书吧。
              </div>
              <AccessGroupEdit
                trigger={<Button>新增授权组</Button>}
                className="mt-3"
              />
            </div>
          </>
        </Show>

        <ScrollArea className="h-[75vh] overflow-hidden">
          <div className="flex gap-5 flex-wrap">
            {accessGroups.map((accessGroup) => (
              <Card className="w-full md:w-[350px]">
                <CardHeader>
                  <CardTitle>{accessGroup.name}</CardTitle>
                  <CardDescription>
                    共有
                    {accessGroup.expand ? accessGroup.expand.access.length : 0}
                    个部署授权配置
                  </CardDescription>
                </CardHeader>
                <CardContent className="min-h-[180px]">
                  {accessGroup.expand ? (
                    <>
                      {accessGroup.expand.access.slice(0, 3).map((access) => (
                        <div key={access.id} className="flex flex-col mb-3">
                          <div className="flex items-center">
                            <div className="">
                              <img
                                src={getProviderInfo(access.configType)![1]}
                                alt="provider"
                                className="w-8 h-8"
                              ></img>
                            </div>
                            <div className="ml-3">
                              <div className="text-sm font-semibold text-gray-700 dark:text-gray-200">
                                {access.name}
                              </div>
                              <div className="text-xs text-muted-foreground">
                                {getProviderInfo(access.configType)![0]}
                              </div>
                            </div>
                          </div>
                        </div>
                      ))}
                    </>
                  ) : (
                    <>
                      <div className="flex text-gray-700 dark:text-gray-200 items-center">
                        <div>
                          <Group size={40} />
                        </div>
                        <div className="ml-2">
                          暂无部署授权配置，请添加后开始使用吧
                        </div>
                      </div>
                    </>
                  )}
                </CardContent>
                <CardFooter>
                  <div className="flex justify-end w-full">
                    <Show
                      when={
                        accessGroup.expand &&
                        accessGroup.expand.access.length > 0
                          ? true
                          : false
                      }
                    >
                      <div>
                        <Button
                          size="sm"
                          variant={"link"}
                          onClick={() => {
                            navigate(`/access?accessGroupId=${accessGroup.id}`);
                          }}
                        >
                          所有授权
                        </Button>
                      </div>
                    </Show>

                    <Show
                      when={
                        !accessGroup.expand ||
                        accessGroup.expand.access.length == 0
                          ? true
                          : false
                      }
                    >
                      <div>
                        <Button size="sm" onClick={handleAddAccess}>
                          新增授权
                        </Button>
                      </div>
                    </Show>

                    <div className="ml-3">
                      <AlertDialog>
                        <AlertDialogTrigger asChild>
                          <Button variant={"destructive"} size={"sm"}>
                            删除
                          </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle className="dark:text-gray-200">
                              删除组
                            </AlertDialogTitle>
                            <AlertDialogDescription>
                              确定要删除部署授权组吗？
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel className="dark:text-gray-200">
                              取消
                            </AlertDialogCancel>
                            <AlertDialogAction
                              onClick={() => {
                                handleRemoveClick(
                                  accessGroup.id ? accessGroup.id : ""
                                );
                              }}
                            >
                              确认
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    </div>
                  </div>
                </CardFooter>
              </Card>
            ))}
          </div>
        </ScrollArea>
      </div>
    </div>
  );
};

export default AccessGroups;
