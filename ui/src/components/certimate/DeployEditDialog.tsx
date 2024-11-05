import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Plus } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ScrollArea } from "@/components/ui/scroll-area";
import AccessEditDialog from "./AccessEditDialog";
import { Context as DeployEditContext, type DeployEditContext as DeployEditContextType } from "./DeployEdit";
import DeployToAliyunOSS from "./DeployToAliyunOSS";
import DeployToAliyunCDN from "./DeployToAliyunCDN";
import DeployToAliyunCLB from "./DeployToAliyunCLB";
import DeployToAliyunALB from "./DeployToAliyunALB";
import DeployToAliyunNLB from "./DeployToAliyunNLB";
import DeployToTencentCDN from "./DeployToTencentCDN";
import DeployToTencentCLB from "./DeployToTencentCLB";
import DeployToTencentCOS from "./DeployToTencentCOS";
import DeployToTencentTEO from "./DeployToTencentTEO";
import DeployToHuaweiCloudCDN from "./DeployToHuaweiCloudCDN";
import DeployToHuaweiCloudELB from "./DeployToHuaweiCloudELB";
import DeployToBaiduCloudCDN from "./DeployToBaiduCloudCDN";
import DeployToQiniuCDN from "./DeployToQiniuCDN";
import DeployToDogeCloudCDN from "./DeployToDogeCloudCDN";
import DeployToLocal from "./DeployToLocal";
import DeployToSSH from "./DeployToSSH";
import DeployToWebhook from "./DeployToWebhook";
import DeployToKubernetesSecret from "./DeployToKubernetesSecret";
import { deployTargetsMap, type DeployConfig } from "@/domain/domain";
import { accessProvidersMap } from "@/domain/access";
import { useConfigContext } from "@/providers/config";

type DeployEditDialogProps = {
  trigger: React.ReactNode;
  deployConfig?: DeployConfig;
  onSave: (deploy: DeployConfig) => void;
};

const DeployEditDialog = ({ trigger, deployConfig, onSave }: DeployEditDialogProps) => {
  const { t } = useTranslation();

  const {
    config: { accesses },
  } = useConfigContext();

  const [deployType, setDeployType] = useState("");

  const [locDeployConfig, setLocDeployConfig] = useState<DeployConfig>({
    access: "",
    type: "",
  });

  const [errors, setErrors] = useState<Record<string, string | undefined>>({});

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
    setDeployType(locDeployConfig.type);
    setErrors({});
  }, [locDeployConfig.type]);

  const setConfig = useCallback(
    (deploy: DeployConfig) => {
      if (deploy.type !== locDeployConfig.type) {
        setLocDeployConfig({ ...deploy, access: "", config: {} });
      } else {
        setLocDeployConfig({ ...deploy });
      }
    },
    [locDeployConfig.type]
  );

  const targetAccesses = accesses.filter((item) => {
    if (item.usage == "apply") {
      return false;
    }

    if (locDeployConfig.type == "") {
      return true;
    }

    return item.configType === deployTargetsMap.get(locDeployConfig.type)?.provider;
  });

  const handleSaveClick = () => {
    // 验证数据
    const newError = { ...errors };
    newError.type = locDeployConfig.type === "" ? t("domain.deployment.form.access.placeholder") : "";
    newError.access = locDeployConfig.access === "" ? t("domain.deployment.form.access.placeholder") : "";
    setErrors(newError);
    if (Object.values(newError).some((e) => !!e)) return;

    // 保存数据
    onSave(locDeployConfig);

    // 清理数据
    setLocDeployConfig({
      access: "",
      type: "",
    });
    setErrors({});

    // 关闭弹框
    setOpen(false);
  };

  let childComponent = <></>;
  switch (deployType) {
    case "aliyun-oss":
      childComponent = <DeployToAliyunOSS />;
      break;
    case "aliyun-cdn":
    case "aliyun-dcdn":
      childComponent = <DeployToAliyunCDN />;
      break;
    case "aliyun-clb":
      childComponent = <DeployToAliyunCLB />;
      break;
    case "aliyun-alb":
      childComponent = <DeployToAliyunALB />;
      break;
    case "aliyun-nlb":
      childComponent = <DeployToAliyunNLB />;
      break;
    case "tencent-cdn":
    case "tencent-ecdn":
      childComponent = <DeployToTencentCDN />;
      break;
    case "tencent-clb":
      childComponent = <DeployToTencentCLB />;
      break;
    case "tencent-cos":
      childComponent = <DeployToTencentCOS />;
      break;
    case "tencent-teo":
      childComponent = <DeployToTencentTEO />;
      break;
    case "huaweicloud-cdn":
      childComponent = <DeployToHuaweiCloudCDN />;
      break;
    case "huaweicloud-elb":
      childComponent = <DeployToHuaweiCloudELB />;
      break;
    case "baiducloud-cdn":
      childComponent = <DeployToBaiduCloudCDN />;
      break;
    case "qiniu-cdn":
      childComponent = <DeployToQiniuCDN />;
      break;
    case "dogecloud-cdn":
      childComponent = <DeployToDogeCloudCDN />;
      break;
    case "local":
      childComponent = <DeployToLocal />;
      break;
    case "ssh":
      childComponent = <DeployToSSH />;
      break;
    case "webhook":
      childComponent = <DeployToWebhook />;
      break;
    case "k8s-secret":
      childComponent = <DeployToKubernetesSecret />;
      break;
  }

  return (
    <DeployEditContext.Provider
      value={{
        config: locDeployConfig as DeployEditContextType["config"],
        setConfig: setConfig as DeployEditContextType["setConfig"],
        errors: errors as DeployEditContextType["errors"],
        setErrors: setErrors as DeployEditContextType["setErrors"],
      }}
    >
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger>{trigger}</DialogTrigger>
        <DialogContent
          className="dark:text-stone-200"
          onInteractOutside={(event) => {
            event.preventDefault();
          }}
        >
          <DialogHeader>
            <DialogTitle>{t("domain.deployment.tab")}</DialogTitle>
            <DialogDescription></DialogDescription>
          </DialogHeader>

          <ScrollArea className="max-h-[80vh]">
            <div className="container py-3">
              {/* 部署方式 */}
              <div>
                <Label>{t("domain.deployment.form.type.label")}</Label>

                <Select
                  value={locDeployConfig.type}
                  onValueChange={(val: string) => {
                    setConfig({ ...locDeployConfig, type: val });
                  }}
                >
                  <SelectTrigger className="mt-2">
                    <SelectValue placeholder={t("domain.deployment.form.type.placeholder")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectLabel>{t("domain.deployment.form.type.list")}</SelectLabel>
                      {Array.from(deployTargetsMap.entries()).map(([key, target]) => (
                        <SelectItem key={key} value={key}>
                          <div className="flex items-center space-x-2">
                            <img className="w-6" src={target.icon} />
                            <div>{t(target.name)}</div>
                          </div>
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>

                <div className="text-red-500 text-sm mt-1">{errors.type}</div>
              </div>

              {/* 授权配置 */}
              <div className="mt-8">
                <Label className="flex justify-between">
                  <div>{t("domain.deployment.form.access.label")}</div>
                  <AccessEditDialog
                    trigger={
                      <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                        <Plus size={14} />
                        {t("common.add")}
                      </div>
                    }
                    op="add"
                  />
                </Label>

                <Select
                  value={locDeployConfig.access}
                  onValueChange={(val: string) => {
                    setConfig({ ...locDeployConfig, access: val });
                  }}
                >
                  <SelectTrigger className="mt-2">
                    <SelectValue placeholder={t("domain.deployment.form.access.placeholder")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectLabel>{t("domain.deployment.form.access.list")}</SelectLabel>
                      {targetAccesses.map((item) => (
                        <SelectItem key={item.id} value={item.id}>
                          <div className="flex items-center space-x-2">
                            <img className="w-6" src={accessProvidersMap.get(item.configType)?.icon} />
                            <div>{item.name}</div>
                          </div>
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>

                <div className="text-red-500 text-sm mt-1">{errors.access}</div>
              </div>

              {/* 其他参数 */}
              <div className="mt-8">{childComponent}</div>
            </div>
          </ScrollArea>

          <DialogFooter>
            <Button
              onClick={(e) => {
                e.stopPropagation();
                handleSaveClick();
              }}
            >
              {t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </DeployEditContext.Provider>
  );
};

export default DeployEditDialog;
