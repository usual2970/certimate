import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";

import { useState } from "react";
import { useTranslation } from "react-i18next";

import AccessTencentForm from "./AccessTencentForm";

import { Label } from "../ui/label";

import { Access, accessTypeMap } from "@/domain/access";
import AccessAliyunForm from "./AccessAliyunForm";
import { cn } from "@/lib/utils";

import AccessSSHForm from "./AccessSSHForm";
import WebhookForm from "./AccessWebhookFrom";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import AccessCloudflareForm from "./AccessCloudflareForm";
import AccessQiniuForm from "./AccessQiniuForm";
import AccessNamesiloForm from "./AccessNamesiloForm";
import AccessGodaddyFrom from "./AccessGodaddyForm";
import AccessLocalForm from "./AccessLocalForm";

type TargetConfigEditProps = {
  op: "add" | "edit" | "copy";
  className?: string;
  trigger: React.ReactNode;
  data?: Access;
};
export function AccessEdit({
  trigger,
  op,
  data,
  className,
}: TargetConfigEditProps) {
  const [open, setOpen] = useState(false);
  const { t } = useTranslation();

  const typeKeys = Array.from(accessTypeMap.keys());

  const [configType, setConfigType] = useState(data?.configType || "");

  let form = <> </>;
  switch (configType) {
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
        <WebhookForm
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
        <AccessGodaddyFrom
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
          <DialogTitle>{op == "add" ? t('access.add') : op == "edit" ? t('access.edit') : t('access.copy')}</DialogTitle>
        </DialogHeader>
        <ScrollArea className="max-h-[80vh]">
          <div className="container py-3">
            <Label>{t('access.type')}</Label>

            <Select
              onValueChange={(val) => {
                setConfigType(val);
              }}
              defaultValue={configType}
            >
              <SelectTrigger className="mt-3">
                <SelectValue placeholder={t('access.type.not.empty')} />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>{t('access.type')}</SelectLabel>
                  {typeKeys.map((key) => (
                    <SelectItem value={key} key={key}>
                      <div
                        className={cn(
                          "flex items-center space-x-2 rounded cursor-pointer",
                          getOptionCls(key)
                        )}
                      >
                        <img
                          src={accessTypeMap.get(key)?.[1]}
                          className="h-6 w-6"
                        />
                        <div>{t(accessTypeMap.get(key)?.[0] || '')}</div>
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
}
