import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";
import { cn } from "@/lib/utils";

const DeployToLocal = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          certPath: "/etc/nginx/ssl/nginx.crt",
          keyPath: "/etc/nginx/ssl/nginx.key",
          shell: "sh",
          preCommand: "",
          command: "sudo service nginx reload",
        },
      });
    }
  }, []);

  useEffect(() => {
    setError({});
  }, []);

  const formSchema = z.object({
    certPath: z.string().min(1, t("domain.deployment.form.file_cert_path.placeholder")),
    keyPath: z.string().min(1, t("domain.deployment.form.file_key_path.placeholder")),
    shell: z.union([z.literal("sh"), z.literal("cmd"), z.literal("powershell")], {
      message: t("domain.deployment.form.shell.placeholder"),
    }),
    preCommand: z.string().optional(),
    command: z.string().optional(),
  });

  useEffect(() => {
    const res = formSchema.safeParse(data.config);
    if (!res.success) {
      setError({
        ...error,
        certPath: res.error.errors.find((e) => e.path[0] === "certPath")?.message,
        keyPath: res.error.errors.find((e) => e.path[0] === "keyPath")?.message,
        shell: res.error.errors.find((e) => e.path[0] === "shell")?.message,
        preCommand: res.error.errors.find((e) => e.path[0] === "preCommand")?.message,
        command: res.error.errors.find((e) => e.path[0] === "command")?.message,
      });
    } else {
      setError({
        ...error,
        certPath: undefined,
        keyPath: undefined,
        shell: undefined,
        preCommand: undefined,
        command: undefined,
      });
    }
  }, [data]);

  const getOptionCls = (val: string) => {
    if (data.config?.shell === val) {
      return "border-primary dark:border-primary";
    }

    return "";
  };

  return (
    <>
      <div className="flex flex-col space-y-8">
        <div>
          <Label>{t("domain.deployment.form.file_cert_path.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.file_cert_path.label")}
            className="w-full mt-1"
            value={data?.config?.certPath}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.certPath = e.target.value?.trim();
              });
              setDeploy(newData);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{error?.certPath}</div>
        </div>

        <div>
          <Label>{t("domain.deployment.form.file_key_path.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.file_key_path.placeholder")}
            className="w-full mt-1"
            value={data?.config?.keyPath}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.keyPath = e.target.value?.trim();
              });
              setDeploy(newData);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{error?.keyPath}</div>
        </div>

        <div>
          <Label>{t("domain.deployment.form.shell.label")}</Label>
          <RadioGroup
            className="flex mt-1"
            value={data?.config?.shell}
            onValueChange={(val) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.shell = val;
              });
              setDeploy(newData);
            }}
          >
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="sh" id="shellOptionSh" />
              <Label htmlFor="shellOptionSh">
                <div className={cn("flex items-center space-x-2 border p-2 rounded cursor-pointer dark:border-stone-700", getOptionCls("sh"))}>
                  <div>POSIX Bash (Linux)</div>
                </div>
              </Label>
            </div>
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="cmd" id="shellOptionCmd" />
              <Label htmlFor="shellOptionCmd">
                <div className={cn("border p-2 rounded cursor-pointer dark:border-stone-700", getOptionCls("cmd"))}>
                  <div>CMD (Windows)</div>
                </div>
              </Label>
            </div>
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="powershell" id="shellOptionPowerShell" />
              <Label htmlFor="shellOptionPowerShell">
                <div className={cn("border p-2 rounded cursor-pointer dark:border-stone-700", getOptionCls("powershell"))}>
                  <div>PowerShell (Windows)</div>
                </div>
              </Label>
            </div>
          </RadioGroup>
          <div className="text-red-600 text-sm mt-1">{error?.shell}</div>
        </div>

        <div>
          <Label>{t("domain.deployment.form.shell_pre_command.label")}</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.preCommand}
            placeholder={t("domain.deployment.form.shell_pre_command.placeholder")}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.preCommand = e.target.value;
              });
              setDeploy(newData);
            }}
          ></Textarea>
          <div className="text-red-600 text-sm mt-1">{error?.preCommand}</div>
        </div>

        <div>
          <Label>{t("domain.deployment.form.shell_command.label")}</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.command}
            placeholder={t("domain.deployment.form.shell_command.placeholder")}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.command = e.target.value;
              });
              setDeploy(newData);
            }}
          ></Textarea>
          <div className="text-red-600 text-sm mt-1">{error?.command}</div>
        </div>
      </div>
    </>
  );
};

export default DeployToLocal;
