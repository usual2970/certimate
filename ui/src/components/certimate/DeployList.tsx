import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { Button } from "../ui/button";
import { EditIcon, Plus, Trash2 } from "lucide-react";
import {
  DeployConfig,
  KVType,
  targetTypeKeys,
  targetTypeMap,
} from "@/domain/domain";
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
import { accessTypeMap } from "@/domain/access";
import { useTranslation } from "react-i18next";
import { AccessEdit } from "./AccessEdit";
import { Input } from "../ui/input";
import { Textarea } from "../ui/textarea";
import KVList from "./KVList";
import { produce } from "immer";

type DeployEditContextProps = {
  deploy: DeployConfig;
  setDeploy: (deploy: DeployConfig) => void;
};

const DeployEditContext = createContext<DeployEditContextProps>(
  {} as DeployEditContextProps
);

export const useDeployEditContext = () => {
  return useContext(DeployEditContext);
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
    </>
  );
};

type DeployEditDialogProps = {
  trigger: React.ReactNode;
  deployConfig?: DeployConfig;
};
const DeployEditDialog = ({ trigger, deployConfig }: DeployEditDialogProps) => {
  const {
    config: { accesses },
  } = useConfig();

  const [deployType, setDeployType] = useState<TargetType>();

  const [locDeployConfig, setLocDeployConfig] = useState<DeployConfig>({
    access: "",
    type: "",
  });

  useEffect(() => {
    if (deployConfig) {
      setLocDeployConfig({ ...deployConfig });
    } else {
      setLocDeployConfig({
        access: "",
        type: "",
      });
    }
  }, [deployConfig]);

  useEffect(() => {
    const temp = locDeployConfig.type.split("-");
    console.log(temp);
    let t;
    if (temp && temp.length > 1) {
      t = temp[1];
    } else {
      t = locDeployConfig.type;
    }
    setDeployType(t as TargetType);
  }, [locDeployConfig.type]);

  const setDeploy = useCallback(
    (deploy: DeployConfig) => {
      if (deploy.type !== locDeployConfig.type) {
        setLocDeployConfig({ ...deploy, access: "", config: {} });
      } else {
        setLocDeployConfig({ ...deploy });
      }
    },
    [locDeployConfig.type]
  );

  const { t } = useTranslation();

  const targetAccesses = accesses.filter((item) => {
    if (item.usage == "apply") {
      return false;
    }

    if (locDeployConfig.type == "") {
      return true;
    }
    const types = locDeployConfig.type.split("-");
    return item.configType === types[0];
  });

  return (
    <DeployEditContext.Provider
      value={{
        deploy: locDeployConfig,
        setDeploy: setDeploy,
      }}
    >
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
              value={locDeployConfig.type}
              onValueChange={(val: string) => {
                setDeploy({ ...locDeployConfig, type: val });
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
                        <img
                          className="w-6"
                          src={targetTypeMap.get(item)?.[1]}
                        />
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
              value={locDeployConfig.access}
              onValueChange={(val: string) => {
                setDeploy({ ...locDeployConfig, access: val });
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
    </DeployEditContext.Provider>
  );
};

type TargetType = "ssh" | "cdn" | "webhook" | "local" | "oss" | "dcdn";

type DeployEditProps = {
  type: TargetType;
};
const DeployEdit = ({ type }: DeployEditProps) => {
  const getDeploy = () => {
    switch (type) {
      case "ssh":
        return <DeploySSH />;
      case "local":
        return <DeploySSH />;
      case "cdn":
        return <DeployCDN />;
      case "dcdn":
        return <DeployCDN />;
      case "oss":
        return <DeployCDN />;
      case "webhook":
        return <DeployWebhook />;
      default:
        return <DeployCDN />;
    }
  };
  return getDeploy();
};

const DeploySSH = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy } = useDeployEditContext();
  return (
    <>
      <div className="flex flex-col space-y-2">
        <div>
          <Label>{t("access.form.ssh.cert.path")}</Label>
          <Input
            placeholder={t("access.form.ssh.cert.path")}
            className="w-full mt-1"
            value={data?.config?.certPath}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.certPath = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>
        <div>
          <Label>{t("access.form.ssh.key.path")}</Label>
          <Input
            placeholder={t("access.form.ssh.key.path")}
            className="w-full mt-1"
            value={data?.config?.keyPath}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.keyPath = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>

        <div>
          <Label>前置命令</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.preCommand}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.preCommand = e.target.value;
              });
              setDeploy(newData);
            }}
          ></Textarea>
        </div>

        <div>
          <Label>命令</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.command}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.command = e.target.value;
              });
              setDeploy(newData);
            }}
          ></Textarea>
        </div>
      </div>
    </>
  );
};

const DeployCDN = () => {
  const { deploy: data, setDeploy } = useDeployEditContext();
  return (
    <div className="flex flex-col space-y-2">
      <div>
        <Label>部署至域名</Label>
        <Input
          placeholder="部署至域名"
          className="w-full mt-1"
          value={data?.config?.domain}
          onChange={(e) => {
            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.domain = e.target.value;
            });
            setDeploy(newData);
          }}
        />
      </div>
    </div>
  );
};

const DeployWebhook = () => {
  const { deploy: data, setDeploy } = useDeployEditContext();

  return (
    <>
      <KVList
        variables={data?.config?.variables}
        onValueChange={(variables: KVType[]) => {
          const newData = produce(data, (draft) => {
            if (!draft.config) {
              draft.config = {};
            }
            draft.config.variables = variables;
          });
          setDeploy(newData);
        }}
      />
    </>
  );
};

export default DeployList;
