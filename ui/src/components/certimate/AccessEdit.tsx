import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { useState } from "react";

import AccessTencentForm from "./AccessTencentForm";

import { Label } from "../ui/label";

import { Access, accessTypeMap } from "@/domain/access";
import AccessAliyunForm from "./AccessAliyunForm";
import { cn } from "@/lib/utils";
import { RadioGroup, RadioGroupItem } from "../ui/radio-group";

type TargetConfigEditProps = {
  op: "add" | "edit";
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

  const typeKeys = Array.from(accessTypeMap.keys());

  const [configType, setConfigType] = useState(data?.configType || "");

  let form = <> </>;
  switch (configType) {
    case "tencent":
      form = (
        <AccessTencentForm
          data={data}
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
      <DialogContent className="sm:max-w-[600px] w-full">
        <DialogHeader>
          <DialogTitle>{op == "add" ? "添加" : "编辑"}授权</DialogTitle>
        </DialogHeader>
        <div className="container">
          <Label>服务商</Label>
          <RadioGroup
            value={configType}
            className="flex mt-3 space-x-2"
            onValueChange={setConfigType}
          >
            {typeKeys.map((key) => (
              <div className="flex items-center space-x-2" key={key}>
                <RadioGroupItem value={key} id={key} hidden />
                <Label htmlFor={key}>
                  <div
                    className={cn(
                      "flex items-center space-x-2 border p-2 rounded cursor-pointer",
                      getOptionCls(key)
                    )}
                  >
                    <img src={accessTypeMap.get(key)?.[1]} className="h-6" />
                    <div>{accessTypeMap.get(key)?.[0]}</div>
                  </div>
                </Label>
              </div>
            ))}
          </RadioGroup>

          {form}
        </div>
      </DialogContent>
    </Dialog>
  );
}
