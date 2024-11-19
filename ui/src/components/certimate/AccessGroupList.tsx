import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Group } from "lucide-react";

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
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useToast } from "@/components/ui/use-toast";
import AccessGroupEdit from "./AccessGroupEdit";
import { accessProvidersMap } from "@/domain/access";
import { getErrMessage } from "@/lib/error";
import { useConfigContext } from "@/providers/config";
import { remove } from "@/repository/access_group";

const AccessGroupList = () => {
  const {
    config: { accessGroups },
    reloadAccessGroups,
  } = useConfigContext();

  const { toast } = useToast();

  const navigate = useNavigate();
  const { t } = useTranslation();

  const handleRemoveClick = async (id: string) => {
    try {
      await remove(id);
      reloadAccessGroups();
    } catch (e) {
      toast({
        title: t("common.delete.failed.message"),
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
    <div className="mt-10">
      <Show when={accessGroups.length == 0}>
        <>
          <div className="flex flex-col items-center mt-10">
            <span className="bg-orange-100 p-5 rounded-full">
              <Group size={40} className="text-primary" />
            </span>

            <div className="text-center text-sm text-muted-foreground mt-3">{t("access.group.domains.nodata")}</div>
            <AccessGroupEdit trigger={<Button>{t("access.group.add")}</Button>} className="mt-3" />
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
                  {t("access.group.total", {
                    total: accessGroup.expand ? accessGroup.expand.access.length : 0,
                  })}
                </CardDescription>
              </CardHeader>
              <CardContent className="min-h-[180px]">
                {accessGroup.expand ? (
                  <>
                    {accessGroup.expand.access.slice(0, 3).map((access) => (
                      <div key={access.id} className="flex flex-col mb-3">
                        <div className="flex items-center">
                          <div className="">
                            <img src={accessProvidersMap.get(access.configType)!.icon} alt="provider" className="w-8 h-8"></img>
                          </div>
                          <div className="ml-3">
                            <div className="text-sm font-semibold text-gray-700 dark:text-gray-200">{access.name}</div>
                            <div className="text-xs text-muted-foreground">{access.configType}</div>
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
                      <div className="ml-2">{t("access.group.nodata")}</div>
                    </div>
                  </>
                )}
              </CardContent>
              <CardFooter>
                <div className="flex justify-end w-full">
                  <Show when={accessGroup.expand && accessGroup.expand.access.length > 0 ? true : false}>
                    <div>
                      <Button
                        size="sm"
                        variant={"link"}
                        onClick={() => {
                          navigate(`/access?accessGroupId=${accessGroup.id}&tab=access`, {
                            replace: true,
                          });
                        }}
                      >
                        {t("access.group.domains")}
                      </Button>
                    </div>
                  </Show>

                  <Show when={!accessGroup.expand || accessGroup.expand.access.length == 0 ? true : false}>
                    <div>
                      <Button size="sm" onClick={handleAddAccess}>
                        {t("access.authorization.add")}
                      </Button>
                    </div>
                  </Show>

                  <div className="ml-3">
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant={"destructive"} size={"sm"}>
                          {t("common.delete")}
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle className="dark:text-gray-200">{t("access.group.delete")}</AlertDialogTitle>
                          <AlertDialogDescription>{t("access.group.delete.confirm")}</AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel className="dark:text-gray-200">{t("common.cancel")}</AlertDialogCancel>
                          <AlertDialogAction
                            onClick={() => {
                              handleRemoveClick(accessGroup.id ? accessGroup.id : "");
                            }}
                          >
                            {t("common.confirm")}
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
  );
};

export default AccessGroupList;
