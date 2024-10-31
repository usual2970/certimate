import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";

type DeployToSSHConfigParams = {
  format?: string;
  certPath?: string;
  keyPath?: string;
  pfxPassword?: string;
  jksAlias?: string;
  jksKeypass?: string;
  jksStorepass?: string;
  shell?: string;
  preCommand?: string;
  command?: string;
};

const DeployToSSH = () => {
  const { t } = useTranslation();

  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToSSHConfigParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          format: "pem",
          certPath: "/etc/nginx/ssl/nginx.crt",
          keyPath: "/etc/nginx/ssl/nginx.key",
          command: "sudo service nginx reload",
        },
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  const formSchema = z
    .object({
      format: z.union([z.literal("pem"), z.literal("pfx"), z.literal("jks")], {
        message: t("domain.deployment.form.file_format.placeholder"),
      }),
      certPath: z
        .string()
        .min(1, t("domain.deployment.form.file_cert_path.placeholder"))
        .max(255, t("common.errmsg.string_max", { max: 255 })),
      keyPath: z
        .string()
        .min(0, t("domain.deployment.form.file_key_path.placeholder"))
        .max(255, t("common.errmsg.string_max", { max: 255 })),
      pfxPassword: z.string().optional(),
      jksAlias: z.string().optional(),
      jksKeypass: z.string().optional(),
      jksStorepass: z.string().optional(),
      preCommand: z.string().optional(),
      command: z.string().optional(),
    })
    .refine((data) => (data.format === "pem" ? !!data.keyPath?.trim() : true), {
      message: t("domain.deployment.form.file_key_path.placeholder"),
      path: ["keyPath"],
    })
    .refine((data) => (data.format === "pfx" ? !!data.pfxPassword?.trim() : true), {
      message: t("domain.deployment.form.file_pfx_password.placeholder"),
      path: ["pfxPassword"],
    })
    .refine((data) => (data.format === "jks" ? !!data.jksAlias?.trim() : true), {
      message: t("domain.deployment.form.file_jks_alias.placeholder"),
      path: ["jksAlias"],
    })
    .refine((data) => (data.format === "jks" ? !!data.jksKeypass?.trim() : true), {
      message: t("domain.deployment.form.file_jks_keypass.placeholder"),
      path: ["jksKeypass"],
    })
    .refine((data) => (data.format === "jks" ? !!data.jksStorepass?.trim() : true), {
      message: t("domain.deployment.form.file_jks_storepass.placeholder"),
      path: ["jksStorepass"],
    });

  useEffect(() => {
    const res = formSchema.safeParse(config.config);
    setErrors({
      ...errors,
      format: res.error?.errors?.find((e) => e.path[0] === "format")?.message,
      certPath: res.error?.errors?.find((e) => e.path[0] === "certPath")?.message,
      keyPath: res.error?.errors?.find((e) => e.path[0] === "keyPath")?.message,
      pfxPassword: res.error?.errors?.find((e) => e.path[0] === "pfxPassword")?.message,
      jksAlias: res.error?.errors?.find((e) => e.path[0] === "jksAlias")?.message,
      jksKeypass: res.error?.errors?.find((e) => e.path[0] === "jksKeypass")?.message,
      jksStorepass: res.error?.errors?.find((e) => e.path[0] === "jksStorepass")?.message,
      preCommand: res.error?.errors?.find((e) => e.path[0] === "preCommand")?.message,
      command: res.error?.errors?.find((e) => e.path[0] === "command")?.message,
    });
  }, [config]);

  useEffect(() => {
    if (config.config?.format === "pem") {
      if (/(.pfx|.jks)$/.test(config.config.certPath!)) {
        setConfig(
          produce(config, (draft) => {
            draft.config ??= {};
            draft.config.certPath = config.config!.certPath!.replace(/(.pfx|.jks)$/, ".crt");
          })
        );
      }
    } else if (config.config?.format === "pfx") {
      if (/(.crt|.jks)$/.test(config.config.certPath!)) {
        setConfig(
          produce(config, (draft) => {
            draft.config ??= {};
            draft.config.certPath = config.config!.certPath!.replace(/(.crt|.jks)$/, ".pfx");
          })
        );
      }
    } else if (config.config?.format === "jks") {
      if (/(.crt|.pfx)$/.test(config.config.certPath!)) {
        setConfig(
          produce(config, (draft) => {
            draft.config ??= {};
            draft.config.certPath = config.config!.certPath!.replace(/(.crt|.pfx)$/, ".jks");
          })
        );
      }
    }
  }, [config.config?.format]);

  return (
    <>
      <div className="flex flex-col space-y-8">
        <div>
          <Label>{t("domain.deployment.form.file_format.label")}</Label>
          <Select
            value={config?.config?.format}
            onValueChange={(value) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.format = value;
              });
              setConfig(nv);
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("domain.deployment.form.file_format.placeholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="pem">PEM</SelectItem>
                <SelectItem value="pfx">PFX</SelectItem>
                <SelectItem value="jks">JKS</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          <div className="text-red-600 text-sm mt-1">{errors?.format}</div>
        </div>

        <div>
          <Label>{t("domain.deployment.form.file_cert_path.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.file_cert_path.label")}
            className="w-full mt-1"
            value={config?.config?.certPath}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.certPath = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{errors?.certPath}</div>
        </div>

        {config.config?.format === "pem" ? (
          <div>
            <Label>{t("domain.deployment.form.file_key_path.label")}</Label>
            <Input
              placeholder={t("domain.deployment.form.file_key_path.placeholder")}
              className="w-full mt-1"
              value={config?.config?.keyPath}
              onChange={(e) => {
                const nv = produce(config, (draft) => {
                  draft.config ??= {};
                  draft.config.keyPath = e.target.value?.trim();
                });
                setConfig(nv);
              }}
            />
            <div className="text-red-600 text-sm mt-1">{errors?.keyPath}</div>
          </div>
        ) : (
          <></>
        )}

        {config.config?.format === "pfx" ? (
          <div>
            <Label>{t("domain.deployment.form.file_pfx_password.label")}</Label>
            <Input
              placeholder={t("domain.deployment.form.file_pfx_password.placeholder")}
              className="w-full mt-1"
              value={config?.config?.pfxPassword}
              onChange={(e) => {
                const nv = produce(config, (draft) => {
                  draft.config ??= {};
                  draft.config.pfxPassword = e.target.value?.trim();
                });
                setConfig(nv);
              }}
            />
            <div className="text-red-600 text-sm mt-1">{errors?.pfxPassword}</div>
          </div>
        ) : (
          <></>
        )}

        {config.config?.format === "jks" ? (
          <>
            <div>
              <Label>{t("domain.deployment.form.file_jks_alias.label")}</Label>
              <Input
                placeholder={t("domain.deployment.form.file_jks_alias.placeholder")}
                className="w-full mt-1"
                value={config?.config?.jksAlias}
                onChange={(e) => {
                  const nv = produce(config, (draft) => {
                    draft.config ??= {};
                    draft.config.jksAlias = e.target.value?.trim();
                  });
                  setConfig(nv);
                }}
              />
              <div className="text-red-600 text-sm mt-1">{errors?.jksAlias}</div>
            </div>

            <div>
              <Label>{t("domain.deployment.form.file_jks_keypass.label")}</Label>
              <Input
                placeholder={t("domain.deployment.form.file_jks_keypass.placeholder")}
                className="w-full mt-1"
                value={config?.config?.jksKeypass}
                onChange={(e) => {
                  const nv = produce(config, (draft) => {
                    draft.config ??= {};
                    draft.config.jksKeypass = e.target.value?.trim();
                  });
                  setConfig(nv);
                }}
              />
              <div className="text-red-600 text-sm mt-1">{errors?.jksKeypass}</div>
            </div>

            <div>
              <Label>{t("domain.deployment.form.file_jks_storepass.label")}</Label>
              <Input
                placeholder={t("domain.deployment.form.file_jks_storepass.placeholder")}
                className="w-full mt-1"
                value={config?.config?.jksStorepass}
                onChange={(e) => {
                  const nv = produce(config, (draft) => {
                    draft.config ??= {};
                    draft.config.jksStorepass = e.target.value?.trim();
                  });
                  setConfig(nv);
                }}
              />
              <div className="text-red-600 text-sm mt-1">{errors?.jksStorepass}</div>
            </div>
          </>
        ) : (
          <></>
        )}

        <div>
          <Label>{t("domain.deployment.form.shell_pre_command.label")}</Label>
          <Textarea
            className="mt-1"
            value={config?.config?.preCommand}
            placeholder={t("domain.deployment.form.shell_pre_command.placeholder")}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.preCommand = e.target.value;
              });
              setConfig(nv);
            }}
          ></Textarea>
          <div className="text-red-600 text-sm mt-1">{errors?.preCommand}</div>
        </div>

        <div>
          <Label>{t("domain.deployment.form.shell_command.label")}</Label>
          <Textarea
            className="mt-1"
            value={config?.config?.command}
            placeholder={t("domain.deployment.form.shell_command.placeholder")}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.command = e.target.value;
              });
              setConfig(nv);
            }}
          ></Textarea>
          <div className="text-red-600 text-sm mt-1">{errors?.command}</div>
        </div>
      </div>
    </>
  );
};

export default DeployToSSH;
