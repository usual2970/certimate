import { AccessEdit } from "@/components/certimate/AccessEdit";
import AccessGroupEdit from "@/components/certimate/AccessGroupEdit";
import AccessGroupList from "@/components/certimate/AccessGroupList";
import XPagination from "@/components/certimate/XPagination";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Access as AccessType, accessTypeMap } from "@/domain/access";
import { convertZulu2Beijing } from "@/lib/time";
import { useConfig } from "@/providers/config";
import { remove } from "@/repository/access";
import { t } from "i18next";
import { Key } from "lucide-react";
import { useLocation, useNavigate } from "react-router-dom";

const Access = () => {
  const { config, deleteAccess } = useConfig();
  const { accesses } = config;

  const perPage = 10;

  const totalPages = Math.ceil(accesses.length / perPage);

  const navigate = useNavigate();
  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const page = query.get("page");
  const pageNumber = page ? Number(page) : 1;

  const tab = query.get("tab");

  const accessGroupId = query.get("accessGroupId");

  const startIndex = (pageNumber - 1) * perPage;
  const endIndex = startIndex + perPage;

  const handleDelete = async (data: AccessType) => {
    const rs = await remove(data);
    deleteAccess(rs.id);
  };

  const handleTabItemClick = (tab: string) => {
    query.set("tab", tab);
    navigate({ search: query.toString() });
  };

  return (
    <div className="">
      <div className="flex justify-between items-center">
        <div className="text-muted-foreground">{t('access.management')}</div>
        {tab != "access_group" ? (
          <AccessEdit trigger={<Button>{t('access.add')}</Button>} op="add" />
        ) : (
          <AccessGroupEdit trigger={<Button>{t('access.group.add')}</Button>} />
        )}
      </div>

      <Tabs
        defaultValue={tab ? tab : "access"}
        value={tab ? tab : "access"}
        className="w-full mt-5"
      >
        <TabsList className="space-x-5 px-3">
          <TabsTrigger
            value="access"
            onClick={() => {
              handleTabItemClick("access");
            }}
          >
            {t('access.management')}
          </TabsTrigger>
          <TabsTrigger
            value="access_group"
            onClick={() => {
              handleTabItemClick("access_group");
            }}
          >
            {t('access.group.management')}
          </TabsTrigger>
        </TabsList>
        <TabsContent value="access">
          {accesses.length === 0 ? (
            <div className="flex flex-col items-center mt-10">
              <span className="bg-orange-100 p-5 rounded-full">
                <Key size={40} className="text-primary" />
              </span>

              <div className="text-center text-sm text-muted-foreground mt-3">
                {t('access.empty')}
              </div>
              <AccessEdit
                trigger={<Button>{t('access.add')}</Button>}
                op="add"
                className="mt-3"
              />
            </div>
          ) : (
            <>
              <div className="hidden sm:flex sm:flex-row text-muted-foreground text-sm border-b dark:border-stone-500 sm:p-2 mt-5">
                <div className="w-48">{t('name')}</div>
                <div className="w-48">{t('access.type')}</div>

                <div className="w-52">{t('create.time')}</div>
                <div className="w-52">{t('update.time')}</div>
                <div className="grow">{t('operation')}</div>
              </div>
              <div className="sm:hidden flex text-sm text-muted-foreground">
                {t('access.list')}
              </div>
              {accesses
                .filter((item) => {
                  return accessGroupId ? item.group == accessGroupId : true;
                })
                .slice(startIndex, endIndex)
                .map((access) => (
                  <div
                    className="flex flex-col sm:flex-row text-secondary-foreground border-b dark:border-stone-500 sm:p-2 hover:bg-muted/50 text-sm"
                    key={access.id}
                  >
                    <div className="sm:w-48 w-full pt-1 sm:pt-0 flex items-center">
                      {access.name}
                    </div>
                    <div className="sm:w-48 w-full pt-1 sm:pt-0 flex  items-center space-x-2">
                      <img
                        src={accessTypeMap.get(access.configType)?.[1]}
                        className="w-6"
                      />
                      <div>{t(accessTypeMap.get(access.configType)?.[0] || '')}</div>
                    </div>

                    <div className="sm:w-52 w-full pt-1 sm:pt-0 flex items-center">
                      {t('created.in')}{" "}
                      {access.created && convertZulu2Beijing(access.created)}
                    </div>
                    <div className="sm:w-52 w-full pt-1 sm:pt-0 flex items-center">
                      {t('updated.in')}{" "}
                      {access.updated && convertZulu2Beijing(access.updated)}
                    </div>
                    <div className="flex items-center grow justify-start pt-1 sm:pt-0">
                      <AccessEdit
                        trigger={
                          <Button variant={"link"} className="p-0">
                            {t('edit')}
                          </Button>
                        }
                        op="edit"
                        data={access}
                      />
                      <Separator orientation="vertical" className="h-4 mx-2" />
                      <Button
                        variant={"link"}
                        className="p-0"
                        onClick={() => {
                          handleDelete(access);
                        }}
                      >
                        {t('delete')}
                      </Button>
                    </div>
                  </div>
                ))}
              <XPagination
                totalPages={totalPages}
                currentPage={pageNumber}
                onPageChange={(page) => {
                  query.set("page", page.toString());
                  navigate({ search: query.toString() });
                }}
              />
            </>
          )}
        </TabsContent>
        <TabsContent value="access_group">
          <AccessGroupList />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Access;
