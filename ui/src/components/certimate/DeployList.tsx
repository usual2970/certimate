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
import { nanoid } from "nanoid";
import { z } from "zod";

type DeployEditContextProps = {
  deploy: DeployConfig;
  error: Record<string, string>;
  setDeploy: (deploy: DeployConfig) => void;
  setError: (error: Record<string, string>) => void;
};

const DeployEditContext = createContext<DeployEditContextProps>(
  {} as DeployEditContextProps
);

export const useDeployEditContext = () => {
  return useContext(DeployEditContext);
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
              <div>{t("deployment.not.added")}</div>
              <div className="flex justify-end mt-2">
                <DeployEditDialog
                  onSave={(config: DeployConfig) => {
                    handleAdd(config);
                  }}
                  trigger={<Button size={"sm"}>{t("add")}</Button>}
                />
              </div>
            </AlertDescription>
          </Alert>
        }
      >
        <div className="flex justify-end py-2 border-b dark:border-stone-400">
          <DeployEditDialog
            trigger={<Button size={"sm"}>{t("add")}</Button>}
            onSave={(config: DeployConfig) => {
              handleAdd(config);
            }}
          />
        </div>

        <div className="w-full md:w-[35em] rounded mt-5 border dark:border-stone-400">
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

type DeployItemProps = {
  item: DeployConfig;
  onDelete: () => void;
  onSave: (deploy: DeployConfig) => void;
};
const DeployItem = ({ item, onDelete, onSave }: DeployItemProps) => {
  const {
    config: { accesses },
  } = useConfig();
  const { t } = useTranslation();
  const access = accesses.find((access) => access.id === item.access);
  const getImg = () => {
    if (!access) {
      return "";
    }

    const accessType = accessTypeMap.get(access.configType);

    if (accessType) {
      return accessType[1];
    }

    return "";
  };

  const getTypeName = () => {
    if (!access) {
      return "";
    }

    const accessType = targetTypeMap.get(item.type);

    if (accessType) {
      return t(accessType[0]);
    }

    return "";
  };

  return (
    <div className="flex justify-between text-sm p-3 items-center text-stone-700">
      <div className="flex space-x-2 items-center">
        <div>
          <img src={getImg()} className="w-9"></img>
        </div>
        <div className="text-stone-600 flex-col flex space-y-0">
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

type DeployEditDialogProps = {
  trigger: React.ReactNode;
  deployConfig?: DeployConfig;
  onSave: (deploy: DeployConfig) => void;
};
const DeployEditDialog = ({
  trigger,
  deployConfig,
  onSave,
}: DeployEditDialogProps) => {
  const {
    config: { accesses },
  } = useConfig();

  const [deployType, setDeployType] = useState<TargetType>();

  const [locDeployConfig, setLocDeployConfig] = useState<DeployConfig>({
    access: "",
    type: "",
  });

  const [error, setError] = useState<Record<string, string>>({});

  const [open, setOpen] = useState(false);

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
    setError({});
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

  const handleSaveClick = () => {
    // 验证数据
    // 保存数据
    // 清理数据
    // 关闭弹框
    const newError = { ...error };
    if (locDeployConfig.type === "") {
      newError.type = t("domain.management.edit.access.not.empty.message");
    } else {
      newError.type = "";
    }

    if (locDeployConfig.access === "") {
      newError.access = t("domain.management.edit.access.not.empty.message");
    } else {
      newError.access = "";
    }
    setError(newError);

    for (const key in newError) {
      if (newError[key] !== "") {
        return;
      }
    }

    onSave(locDeployConfig);

    setLocDeployConfig({
      access: "",
      type: "",
    });

    setError({});

    setOpen(false);
  };

  return (
    <DeployEditContext.Provider
      value={{
        deploy: locDeployConfig,
        setDeploy: setDeploy,
        error: error,
        setError: setError,
      }}
    >
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger>{trigger}</DialogTrigger>
        <DialogContent className="dark:text-stone-200">
          <DialogHeader>
            <DialogTitle>{t("deployment")}</DialogTitle>
            <DialogDescription></DialogDescription>
          </DialogHeader>
          {/* 授权类型 */}
          <div>
            <Label>{t("deployment.access.type")}</Label>

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

            <div className="text-red-500 text-sm mt-1">{error.type}</div>
          </div>
          {/* 授权 */}
          <div>
            <Label className="flex justify-between">
              <div>{t("deployment.access.config")}</div>
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

            <div className="text-red-500 text-sm mt-1">{error.access}</div>
          </div>

          <DeployEdit type={deployType!} />

          <DialogFooter>
            <Button
              onClick={(e) => {
                e.stopPropagation();
                handleSaveClick();
              }}
            >
              {t("save")}
            </Button>
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
        return <DeployOSS />;
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
  const { setError } = useDeployEditContext();

  useEffect(() => {
    setError({});
  }, []);

  const { deploy: data, setDeploy } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          certPath: "/etc/nginx/ssl/nginx.crt",
          keyPath: "/etc/nginx/ssl/nginx.key",
          preCommand: "",
          command: "sudo service nginx reload",
        },
      });
    }
  }, []);
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
          <Label>{t("access.form.ssh.pre.command")}</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.preCommand}
            placeholder={t("access.form.ssh.pre.command.not.empty")}
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
          <Label>{t("access.form.ssh.command")}</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.command}
            placeholder={t("access.form.ssh.command.not.empty")}
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
  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  const { t } = useTranslation();

  useEffect(() => {
    setError({});
  }, []);

  useEffect(() => {
    const resp = domainSchema.safeParse(data.config?.domain);
    if (!resp.success) {
      setError({
        ...error,
        domain: JSON.parse(resp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        domain: "",
      });
    }
  }, [data]);

  const domainSchema = z
    .string()
    .regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: t("domain.not.empty.verify.message"),
    });

  return (
    <div className="flex flex-col space-y-2">
      <div>
        <Label>{t("deployment.access.cdn.deploy.to.domain")}</Label>
        <Input
          placeholder={t("deployment.access.cdn.deploy.to.domain")}
          className="w-full mt-1"
          value={data?.config?.domain}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = domainSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                domain: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                domain: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.domain = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.domain}</div>
      </div>
    </div>
  );
};

const DeployOSS = () => {
  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  const { t } = useTranslation();

  useEffect(() => {
    setError({});
  }, []);

  useEffect(() => {
    const resp = domainSchema.safeParse(data.config?.domain);
    if (!resp.success) {
      setError({
        ...error,
        domain: JSON.parse(resp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        domain: "",
      });
    }
  }, [data]);

  useEffect(() => {
    const bucketResp = bucketSchema.safeParse(data.config?.domain);
    if (!bucketResp.success) {
      setError({
        ...error,
        bucket: JSON.parse(bucketResp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        bucket: "",
      });
    }
  }, []);

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          endpoint: "oss-cn-hangzhou.aliyuncs.com",
          bucket: "",
          domain: "",
        },
      });
    }
  }, []);

  const domainSchema = z
    .string()
    .regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: t("domain.not.empty.verify.message"),
    });

  const bucketSchema = z.string().min(1, {
    message: t("deployment.access.oss.bucket.not.empty"),
  });

  return (
    <div className="flex flex-col space-y-2">
      <div>
        <Label>{t("deployment.access.oss.endpoint")}</Label>

        <Input
          className="w-full mt-1"
          value={data?.config?.endpoint}
          onChange={(e) => {
            const temp = e.target.value;

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.endpoint = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.endpoint}</div>

        <Label>{t("deployment.access.oss.bucket")}</Label>
        <Input
          placeholder={t("deployment.access.oss.bucket.not.empty")}
          className="w-full mt-1"
          value={data?.config?.bucket}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = bucketSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                bucket: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                bucket: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.bucket = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.bucket}</div>

        <Label>{t("deployment.access.cdn.deploy.to.domain")}</Label>
        <Input
          placeholder={t("deployment.access.cdn.deploy.to.domain")}
          className="w-full mt-1"
          value={data?.config?.domain}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = domainSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                domain: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                domain: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.domain = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.domain}</div>
      </div>
    </div>
  );
};

const DeployWebhook = () => {
  const { deploy: data, setDeploy } = useDeployEditContext();

  const { setError } = useDeployEditContext();

  useEffect(() => {
    setError({});
  }, []);

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
