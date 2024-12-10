import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { produce } from "immer";
import { Edit, Plus, Trash2 } from "lucide-react";

import Show from "@/components/Show";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { KVType } from "@/domain/domain";

type KVListProps = {
  variables?: KVType[];
  onValueChange?: (variables: KVType[]) => void;
};

const KVList = ({ variables, onValueChange }: KVListProps) => {
  const [locVariables, setLocVariables] = useState<KVType[]>([]);

  const { t } = useTranslation();

  useEffect(() => {
    if (variables) {
      setLocVariables(variables);
    }
  }, [variables]);

  const handleAddClick = (variable: KVType) => {
    // 查看是否存在key，存在则更新，不存在则添加
    const index = locVariables.findIndex((item) => {
      return item.key === variable.key;
    });

    const newList = produce(locVariables, (draft) => {
      if (index === -1) {
        draft.push(variable);
      } else {
        draft[index] = variable;
      }
    });

    setLocVariables(newList);

    onValueChange?.(newList);
  };

  const handleDeleteClick = (index: number) => {
    const newList = [...locVariables];
    newList.splice(index, 1);
    setLocVariables(newList);

    onValueChange?.(newList);
  };

  const handleEditClick = (index: number, variable: KVType) => {
    const newList = [...locVariables];
    newList[index] = variable;
    setLocVariables(newList);

    onValueChange?.(newList);
  };

  return (
    <>
      <div className="flex justify-between dark:text-stone-200">
        <Label>{t("domain.deployment.form.variables.label")}</Label>
        <Show when={!!locVariables?.length}>
          <KVEdit
            variable={{
              key: "",
              value: "",
            }}
            trigger={
              <div className="flex items-center text-primary">
                <Plus size={16} className="cursor-pointer " />

                <div className="text-sm ">{t("common.button.add")}</div>
              </div>
            }
            onSave={(variable) => {
              handleAddClick(variable);
            }}
          />
        </Show>
      </div>

      <Show
        when={!!locVariables?.length}
        fallback={
          <div className="border rounded-md p-3 text-sm flex flex-col items-center">
            <div className="text-muted-foreground">{t("domain.deployment.form.variables.empty")}</div>

            <KVEdit
              trigger={
                <div className="flex items-center text-primary">
                  <Plus size={16} className="cursor-pointer " />

                  <div className="text-sm ">{t("common.button.add")}</div>
                </div>
              }
              variable={{
                key: "",
                value: "",
              }}
              onSave={(variable) => {
                handleAddClick(variable);
              }}
            />
          </div>
        }
      >
        <div className="border p-3 rounded-md text-stone-700 text-sm dark:text-stone-200">
          {locVariables?.map((item, index) => (
            <div key={index} className="flex justify-between items-center">
              <div>
                {item.key}={item.value}
              </div>
              <div className="flex space-x-2">
                <KVEdit
                  trigger={<Edit size={16} className="cursor-pointer" />}
                  variable={item}
                  onSave={(variable) => {
                    handleEditClick(index, variable);
                  }}
                />

                <Trash2
                  size={16}
                  className="cursor-pointer"
                  onClick={() => {
                    handleDeleteClick(index);
                  }}
                />
              </div>
            </div>
          ))}
        </div>
      </Show>
    </>
  );
};

type KVEditProps = {
  variable?: KVType;
  trigger: React.ReactNode;
  onSave: (variable: KVType) => void;
};

const KVEdit = ({ variable, trigger, onSave }: KVEditProps) => {
  const [locVariable, setLocVariable] = useState<KVType>({
    key: "",
    value: "",
  });

  useEffect(() => {
    if (variable) setLocVariable(variable!);
  }, [variable]);

  const { t } = useTranslation();

  const [open, setOpen] = useState<boolean>(false);

  const [err, setErr] = useState<Record<string, string>>({});

  const handleSaveClick = () => {
    if (!locVariable.key) {
      setErr({
        key: t("domain.deployment.form.variables.key.required"),
      });
      return;
    }

    if (!locVariable.value) {
      setErr({
        value: t("domain.deployment.form.variables.value.required"),
      });
      return;
    }

    onSave?.(locVariable);

    setOpen(false);

    setErr({});
  };

  return (
    <Dialog
      open={open}
      onOpenChange={() => {
        setOpen(!open);
      }}
    >
      <DialogTrigger>{trigger}</DialogTrigger>
      <DialogContent className="dark:text-stone-200">
        <DialogHeader className="flex flex-col">
          <DialogTitle>{t("domain.deployment.form.variables.label")}</DialogTitle>

          <div className="pt-5 flex flex-col items-start">
            <Label>{t("domain.deployment.form.variables.key")}</Label>
            <Input
              placeholder={t("domain.deployment.form.variables.key.placeholder")}
              value={locVariable?.key}
              onChange={(e) => {
                setLocVariable({ ...locVariable, key: e.target.value });
              }}
              className="w-full mt-1"
            />
            <div className="text-red-500 text-sm mt-1">{err?.key}</div>
          </div>

          <div className="pt-2  flex flex-col items-start">
            <Label>{t("domain.deployment.form.variables.value")}</Label>
            <Input
              placeholder={t("domain.deployment.form.variables.value.placeholder")}
              value={locVariable?.value}
              onChange={(e) => {
                setLocVariable({ ...locVariable, value: e.target.value });
              }}
              className="w-full mt-1"
            />

            <div className="text-red-500 text-sm mt-1">{err?.value}</div>
          </div>
        </DialogHeader>
        <DialogFooter>
          <div className="flex justify-end">
            <Button
              onClick={() => {
                handleSaveClick();
              }}
            >
              {t("common.button.save")}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default KVList;
