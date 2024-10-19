import { useState } from "react";
import { useTranslation } from "react-i18next";

import { cn } from "@/lib/utils";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import AccessAliyunForm from "./AccessAliyunForm";
import AccessTencentForm from "./AccessTencentForm";
import AccessHuaweiCloudForm from "./AccessHuaweicloudForm";
import AccessQiniuForm from "./AccessQiniuForm";
import AccessAwsForm from "./AccessAwsForm";
import AccessCloudflareForm from "./AccessCloudflareForm";
import AccessNamesiloForm from "./AccessNamesiloForm";
import AccessGodaddyForm from "./AccessGodaddyForm";
import AccessLocalForm from "./AccessLocalForm";
import AccessSSHForm from "./AccessSSHForm";
import AccessWebhookForm from "./AccessWebhookForm";
import AccessKubernetesForm from "./AccessKubernetesForm";
import { Access, accessTypeMap } from "@/domain/access";

type AccessEditProps = {
  op: "add" | "edit" | "copy";
  className?: string;
  trigger: React.ReactNode;
  data?: Access;
};

const AccessEdit = ({ trigger, op, data, className }: AccessEditProps) => {
  const [open, setOpen] = useState(false);
  const { t } = useTranslation();

  const typeKeys = Array.from(accessTypeMap.keys());

  const [configType, setConfigType] = useState(data?.configType || "");

  let form = <> </>;
  switch (configType) {
    case "aliyun":
      form = (
        <AccessAliyunForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "tencent":
      form = (
        <AccessTencentForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "huaweicloud":
      form = (
        <AccessHuaweiCloudForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "qiniu":
      form = (
        <AccessQiniuForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "aws":
      form = (
        <AccessAwsForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "cloudflare":
      form = (
        <AccessCloudflareForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "namesilo":
      form = (
        <AccessNamesiloForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "godaddy":
      form = (
        <AccessGodaddyForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "local":
      form = (
        <AccessLocalForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "ssh":
      form = (
        <AccessSSHForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "webhook":
      form = (
        <AccessWebhookForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
    case "k8s":
      form = (
        <AccessKubernetesForm
          data={data}
          op={op}
          onAfterReq={() => {
            setOpen(false);
          }}
        />
      );
      break;
  }

  const getOptionCls = (val: string) => {
    return val == configType ? "border-primary" : "";
  };

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild className={cn(className)}>
        {trigger}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[600px] w-full dark:text-stone-200">
        <DialogHeader>
          <DialogTitle>
            {op == "add" ? t("access.authorization.add") : op == "edit" ? t("access.authorization.edit") : t("access.authorization.copy")}
          </DialogTitle>
        </DialogHeader>
        <ScrollArea className="max-h-[80vh]">
          <div className="container py-3">
            <Label>{t("access.authorization.form.type.label")}</Label>

            <Select
              onValueChange={(val) => {
                setConfigType(val);
              }}
              defaultValue={configType}
            >
              <SelectTrigger className="mt-3">
                <SelectValue placeholder={t("access.authorization.form.type.placeholder")} />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>{t("access.authorization.form.type.list")}</SelectLabel>
                  {typeKeys.map((key) => (
                    <SelectItem value={key} key={key}>
                      <div className={cn("flex items-center space-x-2 rounded cursor-pointer", getOptionCls(key))}>
                        <img src={accessTypeMap.get(key)?.[1]} className="h-6 w-6" />
                        <div>{t(accessTypeMap.get(key)?.[0] || "")}</div>
                      </div>
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>

            {form}
          </div>
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
};

export default AccessEdit;
