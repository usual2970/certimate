import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { nanoid } from "nanoid";
import { EditIcon, Trash2 } from "lucide-react";

import Show from "@/components/Show";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import DeployEditDialog from "./DeployEditDialog";
import { DeployConfig } from "@/domain/domain";
import { accessProvidersMap } from "@/domain/access";
import { useConfigContext } from "@/providers/config";

type DeployItemProps = {
  item: DeployConfig;
  onDelete: () => void;
  onSave: (deploy: DeployConfig) => void;
};

const DeployItem = ({ item, onDelete, onSave }: DeployItemProps) => {
  const {
    config: { accesses },
  } = useConfigContext();
  const { t } = useTranslation();

  const access = accesses.find((access) => access.id === item.access);

  const getTypeIcon = () => {
    if (!access) {
      return "";
    }

    return accessProvidersMap.get(access.configType)?.icon || "";
  };

  const getTypeName = () => {
    if (!access) {
      return "";
    }

    return t(accessProvidersMap.get(access.configType)?.name || "");
  };

  return (
    <div className="flex justify-between text-sm p-3 items-center text-stone-700 dark:text-stone-200">
      <div className="flex space-x-2 items-center">
        <div>
          <img src={getTypeIcon()} className="w-9"></img>
        </div>
        <div className="text-stone-600 flex-col flex space-y-0 dark:text-stone-200">
          <div>{getTypeName()}</div>
          <div>{access?.name}</div>
        </div>
      </div>
      <div className="flex space-x-2">
        <DeployEditDialog
          trigger={<EditIcon size={16} className="cursor-pointer" />}
          deployConfig={item}
          onSave={(deploy: DeployConfig) => {
            onSave(deploy);
          }}
        />

        <Trash2
          size={16}
          className="cursor-pointer"
          onClick={() => {
            onDelete();
          }}
        />
      </div>
    </div>
  );
};

type DeployListProps = {
  deploys: DeployConfig[];
  onChange: (deploys: DeployConfig[]) => void;
};

const DeployList = ({ deploys, onChange }: DeployListProps) => {
  const [list, setList] = useState<DeployConfig[]>([]);

  const { t } = useTranslation();

  useEffect(() => {
    setList(deploys);
  }, [deploys]);

  const handleAdd = (deploy: DeployConfig) => {
    deploy.id = nanoid();

    const newList = [...list, deploy];

    setList(newList);

    onChange(newList);
  };

  const handleDelete = (id: string) => {
    const newList = list.filter((item) => item.id !== id);

    setList(newList);

    onChange(newList);
  };

  const handleSave = (deploy: DeployConfig) => {
    const newList = list.map((item) => {
      if (item.id === deploy.id) {
        return { ...deploy };
      }
      return item;
    });

    setList(newList);

    onChange(newList);
  };

  return (
    <>
      <Show
        when={list.length > 0}
        fallback={
          <Alert className="w-full border dark:border-stone-400">
            <AlertDescription className="flex flex-col items-center">
              <div>{t("domain.deployment.nodata")}</div>
              <div className="flex justify-end mt-2">
                <DeployEditDialog
                  onSave={(config: DeployConfig) => {
                    handleAdd(config);
                  }}
                  trigger={<Button size={"sm"}>{t("common.add")}</Button>}
                />
              </div>
            </AlertDescription>
          </Alert>
        }
      >
        <div className="flex justify-end py-2 border-b dark:border-stone-400">
          <DeployEditDialog
            trigger={<Button size={"sm"}>{t("common.add")}</Button>}
            onSave={(config: DeployConfig) => {
              handleAdd(config);
            }}
          />
        </div>

        <div className="w-full md:w-[35em] rounded mt-5 border dark:border-stone-400 dark:text-stone-200">
          <div className="">
            {list.map((item) => (
              <DeployItem
                key={item.id}
                item={item}
                onDelete={() => {
                  handleDelete(item.id ?? "");
                }}
                onSave={(deploy: DeployConfig) => {
                  handleSave(deploy);
                }}
              />
            ))}
          </div>
        </div>
      </Show>
    </>
  );
};

export default DeployList;
