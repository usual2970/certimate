import { cn } from "@/lib/utils";

import Show from "../Show";
import { useCallback, useEffect, useMemo, useState } from "react";
import { FormControl, FormLabel } from "../ui/form";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Input } from "../ui/input";
import { z } from "zod";
import { useTranslation } from "react-i18next";
import { Edit, Plus, Trash2 } from "lucide-react";

type StringListProps = {
  className?: string;
  value: string;
  valueType?: "domain" | "ip";
  onValueChange: (value: string) => void;
};

const titles: Record<string, string> = {
  domain: "domain",
  ip: "IP",
};

const StringList = ({
  value,
  className,
  onValueChange,
  valueType = "domain",
}: StringListProps) => {
  const [list, setList] = useState<string[]>([]);

  const { t } = useTranslation();

  useMemo(() => {
    if (value) {
      setList(value.split(";"));
    }
  }, [value]);

  useEffect(() => {
    const changeList = () => {
      onValueChange(list.join(";"));
    };
    changeList();
  }, [list]);

  const addVal = (val: string) => {
    if (list.includes(val)) {
      return;
    }
    setList([...list, val]);
  };

  const editVal = (index: number, val: string) => {
    const newList = [...list];
    newList[index] = val;
    setList(newList);
  };

  const onRemoveClick = (index: number) => {
    const newList = [...list];
    newList.splice(index, 1);
    setList(newList);
  };

  return (
    <>
      <div className={cn(className)}>
        <FormLabel className="flex justify-between items-center">
          <div>{t(titles[valueType])}</div>

          <Show when={list.length > 0}>
            <StringEdit
              op="add"
              onValueChange={(val: string) => {
                addVal(val);
              }}
              valueType={valueType}
              value={""}
              trigger={
                <div className="flex items-center text-primary">
                  <Plus size={16} className="cursor-pointer " />

                  <div className="text-sm ">{t("add")}</div>
                </div>
              }
            />
          </Show>
        </FormLabel>
        <FormControl>
          <Show
            when={list.length > 0}
            fallback={
              <div className="border rounded-md p-3 text-sm mt-2 flex flex-col items-center">
                <div className="text-muted-foreground">暂未添加域名</div>

                <StringEdit
                  value={""}
                  trigger={t("add")}
                  onValueChange={addVal}
                  valueType={valueType}
                />
              </div>
            }
          >
            <div className="border rounded-md p-3 text-sm mt-2 text-gray-700 space-y-2 dark:text-white dark:border-stone-700 dark:bg-stone-950">
              {list.map((item, index) => (
                <div key={index} className="flex justify-between items-center">
                  <div>{item}</div>
                  <div className="flex space-x-2">
                    <StringEdit
                      op="edit"
                      valueType={valueType}
                      trigger={
                        <Edit
                          size={16}
                          className="cursor-pointer text-gray-600 dark:text-white"
                        />
                      }
                      value={item}
                      onValueChange={(val: string) => {
                        editVal(index, val);
                      }}
                    />
                    <Trash2
                      size={16}
                      className="cursor-pointer"
                      onClick={() => {
                        onRemoveClick(index);
                      }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </Show>
        </FormControl>
      </div>
    </>
  );
};

export default StringList;

type ValueType = "domain" | "ip";

type StringEditProps = {
  value: string;
  trigger: React.ReactNode;
  onValueChange: (value: string) => void;
  valueType: ValueType;
  op?: "add" | "edit";
};

const StringEdit = ({
  trigger,
  value,
  onValueChange,
  op = "add",
  valueType,
}: StringEditProps) => {
  const [currentValue, setCurrentValue] = useState<string>("");
  const [open, setOpen] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const { t } = useTranslation();

  useEffect(() => {
    setCurrentValue(value);
  }, [value]);

  const domainSchema = z
    .string()
    .regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: t("domain.not.empty.verify.message"),
    });

  const ipSchema = z.string().ip({ message: t("ip.not.empty.verify.message") });

  const schedules: Record<ValueType, z.ZodString> = {
    domain: domainSchema,
    ip: ipSchema,
  };

  const onSaveClick = useCallback(() => {
    const schema = schedules[valueType];

    const resp = schema.safeParse(currentValue);
    if (!resp.success) {
      setError(JSON.parse(resp.error.message)[0].message);
      return;
    }

    setCurrentValue("");
    setOpen(false);
    setError("");

    onValueChange(currentValue);
  }, [currentValue]);

  return (
    <Dialog
      open={open}
      onOpenChange={(open) => {
        setOpen(open);
      }}
    >
      <DialogTrigger className="text-primary">{trigger}</DialogTrigger>
      <DialogContent className="dark:text-white">
        <DialogHeader>
          <DialogTitle className="dark:text-white">
            {t(titles[valueType])}
          </DialogTitle>
        </DialogHeader>
        <Input
          value={currentValue}
          className="dark:text-white"
          onChange={(e) => {
            setCurrentValue(e.target.value);
          }}
        />
        <Show when={error.length > 0}>
          <div className="text-red-500 text-sm">{error}</div>
        </Show>

        <DialogFooter>
          <Button
            onClick={() => {
              onSaveClick();
            }}
          >
            {op === "add" ? t("add") : t("confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
