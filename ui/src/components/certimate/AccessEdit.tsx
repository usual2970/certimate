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
import AccessAliyunForm from "./AccessAliyunForm";
import { cn } from "@/lib/utils";

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

          <Select
            onValueChange={(value) => {
              setConfigType(value);
            }}
            value={configType}
          >
            <SelectTrigger className="mt-3">
              <SelectValue placeholder="选择服务商" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>服务商</SelectLabel>
                {typeKeys.map((key) => (
                  <SelectItem value={key}>{accessTypeMap.get(key)}</SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>

          {form}
        </div>
      </DialogContent>
    </Dialog>
  );
}
