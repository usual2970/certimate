import { createContext, useContext, useEffect, useState } from "react";
import { Button } from "../ui/button";
import { EditIcon, Plus, Trash2 } from "lucide-react";
import { DeployConfig, targetTypeKeys, targetTypeMap } from "@/domain/domain";
import Show from "../Show";
import { Alert, AlertDescription } from "../ui/alert";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";

import { Label } from "../ui/label";
import { useConfig } from "@/providers/config";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { Access, accessTypeMap } from "@/domain/access";
import { useTranslation } from "react-i18next";
import { AccessEdit } from "./AccessEdit";
import { Input } from "../ui/input";
import { Textarea } from "../ui/textarea";
import KVList from "./KVList";

type DeployListContextProps = {
  deploys: DeployConfig[];
};

const DeployListContext = createContext<DeployListContextProps>({
  deploys: [],
});

export const useDeployListContext = () => {
  return useContext(DeployListContext);
};

type DeployListProps = {
  deploys: DeployConfig[];
};

const DeployList = ({ deploys }: DeployListProps) => {
  const [list, setList] = useState<DeployConfig[]>([]);

  useEffect(() => {
    setList(deploys);
  }, [deploys]);

  return (
    <>
      <DeployListContext.Provider value={{ deploys: deploys }}>
        <Show
          when={list.length > 0}
          fallback={
            <Alert className="w-full">
              <AlertDescription className="flex flex-col items-center">
                <div>暂无部署配置，请添加后开始部署证书吧</div>
                <div className="flex justify-end mt-2">
                  <DeployEditDialog
                    trigger={<Button size={"sm"}>添加部署</Button>}
                  />
                </div>
              </AlertDescription>
            </Alert>
          }
        >
          <div className="flex justify-end py-2 border-b">
            <DeployEditDialog trigger={<Button size={"sm"}>添加部署</Button>} />
          </div>

          <div className="w-full md:w-[35em] rounded mt-5 border">
            <div className="">
              <div className="flex justify-between text-sm p-3 items-center text-stone-700">
                <div className="flex space-x-2 items-center">
                  <div>
                    <img src="/imgs/providers/ssh.svg" className="w-9"></img>
                  </div>
                  <div className="text-stone-600 flex-col flex space-y-0">
                    <div>ssh部署</div>
                    <div>业务服务器</div>
                  </div>
                </div>
                <div className="flex space-x-2">
                  <EditIcon size={16} className="cursor-pointer" />
                  <Trash2 size={16} className="cursor-pointer" />
                </div>
              </div>
            </div>
          </div>
        </Show>
      </DeployListContext.Provider>
    </>
  );
};

type DeployEditDialogProps = {
  trigger: React.ReactNode;
};
const DeployEditDialog = ({ trigger }: DeployEditDialogProps) => {
  const {
    config: { accesses },
  } = useConfig();

  const [access, setAccess] = useState<Access>();
  const [deployType, setDeployType] = useState<TargetType>();
  const [accessType, setAccessType] = useState("");

  useEffect(() => {
    const temp = accessType.split("-");
    console.log(temp);
    let t;
    if (temp && temp.length > 1) {
      t = temp[1];
    } else {
      t = accessType;
    }
    setDeployType(t as TargetType);
  }, [accessType]);

  const { t } = useTranslation();

  const targetAccesses = accesses.filter((item) => {
    if (item.usage == "apply") {
      return false;
    }

    if (accessType == "") {
      return true;
    }
    const types = accessType.split("-");
    return item.configType === types[0];
  });

  return (
    <Dialog>
      <DialogTrigger>{trigger}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>部署</DialogTitle>
          <DialogDescription></DialogDescription>
        </DialogHeader>
        {/* 授权类型 */}
        <div>
          <Label>授权类型</Label>

          <Select
            value={accessType}
            onValueChange={(val: string) => {
              setAccessType(val);
            }}
          >
            <SelectTrigger className="mt-2">
              <SelectValue
                placeholder={t(
                  "domain.management.edit.access.not.empty.message"
                )}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>
                  {t("domain.management.edit.access.label")}
                </SelectLabel>
                {targetTypeKeys.map((item) => (
                  <SelectItem key={item} value={item}>
                    <div className="flex items-center space-x-2">
                      <img className="w-6" src={targetTypeMap.get(item)?.[1]} />
                      <div>{t(targetTypeMap.get(item)?.[0] ?? "")}</div>
                    </div>
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        {/* 授权 */}
        <div>
          <Label className="flex justify-between">
            <div>授权配置</div>
            <AccessEdit
              trigger={
                <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                  <Plus size={14} />
                  {t("add")}
                </div>
              }
              op="add"
            />
          </Label>

          <Select
            value={access?.id}
            onValueChange={(val: string) => {
              const temp = accesses.find((access) => {
                return access.id == val;
              });
              setAccess(temp!);
            }}
          >
            <SelectTrigger className="mt-2">
              <SelectValue
                placeholder={t(
                  "domain.management.edit.access.not.empty.message"
                )}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>
                  {t("domain.management.edit.access.label")}
                </SelectLabel>
                {targetAccesses.map((item) => (
                  <SelectItem key={item.id} value={item.id}>
                    <div className="flex items-center space-x-2">
                      <img
                        className="w-6"
                        src={accessTypeMap.get(item.configType)?.[1]}
                      />
                      <div>{item.name}</div>
                    </div>
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>

        <DeployEdit type={deployType!} />

        <DialogFooter>
          <Button>保存</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

type TargetType = "ssh" | "cdn" | "webhook" | "local" | "oss" | "dcdn";

type DeployEditProps = {
  type: TargetType;
  data?: DeployConfig;
};
const DeployEdit = ({ type, data }: DeployEditProps) => {
  const getDeploy = () => {
    switch (type) {
      case "ssh":
        return <DeploySSH data={data} />;
      case "local":
        return <DeploySSH data={data} />;
      case "cdn":
        return <DeployCDN data={data} />;
      case "dcdn":
        return <DeployCDN data={data} />;
      case "oss":
        return <DeployCDN data={data} />;
      case "webhook":
        return <DeployWebhook data={data} />;
      default:
        return <DeployCDN data={data} />;
    }
  };
  return getDeploy();
};

type DeployProps = {
  data?: DeployConfig;
};

const DeploySSH = ({ data }: DeployProps) => {
  const { t } = useTranslation();
  return (
    <>
      <div className="flex flex-col space-y-2">
        <div>
          <Label>{t("access.form.ssh.cert.path")}</Label>
          <Input
            placeholder={t("access.form.ssh.cert.path")}
            className="w-full mt-1"
            value={data?.config?.certPath}
          />
        </div>
        <div>
          <Label>{t("access.form.ssh.key.path")}</Label>
          <Input
            placeholder={t("access.form.ssh.key.path")}
            className="w-full mt-1"
            value={data?.config?.keyPath}
          />
        </div>

        <div>
          <Label>前置命令</Label>
          <Textarea className="mt-1"></Textarea>
        </div>

        <div>
          <Label>后置命令</Label>
          <Textarea className="mt-1"></Textarea>
        </div>
      </div>
    </>
  );
};

const DeployCDN = ({ data }: DeployProps) => {
  return (
    <div className="flex flex-col space-y-2">
      <div>
        <Label>部署至域名</Label>
        <Input
          placeholder="部署至域名"
          className="w-full mt-1"
          value={data?.config?.domain}
        />
      </div>
    </div>
  );
};

const DeployWebhook = ({ data }: DeployProps) => {
  return (
    <>
      <KVList variables={data?.config?.variables} />
    </>
  );
};

export default DeployList;
