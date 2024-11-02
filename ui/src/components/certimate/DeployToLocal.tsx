import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";
import { cn } from "@/lib/utils";

type DeployToLocalConfigParams = {
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

const DeployToLocal = () => {
  const { t } = useTranslation();

  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToLocalConfigParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          format: "pem",
          certPath: "/etc/nginx/ssl/nginx.crt",
          keyPath: "/etc/nginx/ssl/nginx.key",
          shell: "sh",
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
      shell: z.union([z.literal("sh"), z.literal("cmd"), z.literal("powershell")], {
        message: t("domain.deployment.form.shell.placeholder"),
      }),
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
      shell: res.error?.errors?.find((e) => e.path[0] === "shell")?.message,
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

  const getOptionCls = (val: string) => {
    if (config.config?.shell === val) {
      return "border-primary dark:border-primary";
    }

    return "";
  };

  const handleUsePresetScript = (key: string) => {
    switch (key) {
      case "reload_nginx":
        {
          setConfig(
            produce(config, (draft) => {
              draft.config ??= {};
              draft.config.shell = "sh";
              draft.config.command = "sudo service nginx reload";
            })
          );
        }
        break;

      case "binding_iis":
        {
          setConfig(
            produce(config, (draft) => {
              draft.config ??= {};
              draft.config.shell = "powershell";
              draft.config.command = `
# 请将以下变量替换为实际值
$pfxPath = "<your-pfx-path>" # PFX 文件路径
$pfxPassword = "<your-pfx-password>" # PFX 密码
$siteName = "<your-site-name>" # IIS 网站名称
$domain = "<your-domain-name>" # 域名
$ipaddr = "<your-binding-ip>"  # 绑定 IP，“*”表示所有 IP 绑定
$port = "<your-binding-port>"  # 绑定端口


# 导入证书到本地计算机的个人存储区
$cert = Import-PfxCertificate -FilePath "$pfxPath" -CertStoreLocation Cert:\\LocalMachine\\My -Password (ConvertTo-SecureString -String "$pfxPassword" -AsPlainText -Force) -Exportable
# 获取 Thumbprint
$thumbprint = $cert.Thumbprint
# 导入 WebAdministration 模块
Import-Module WebAdministration
# 检查是否已存在 HTTPS 绑定
$existingBinding = Get-WebBinding -Name "$siteName" -Protocol "https" -Port $port -HostHeader "$domain" -ErrorAction SilentlyContinue
if (!$existingBinding) {
    # 添加新的 HTTPS 绑定
  New-WebBinding -Name "$siteName" -Protocol "https" -Port $port -IPAddress "$ipaddr" -HostHeader "$domain"
}
# 获取绑定对象
$binding = Get-WebBinding -Name "$siteName" -Protocol "https" -Port $port -IPAddress "$ipaddr" -HostHeader "$domain"
# 绑定 SSL 证书
$binding.AddSslCertificate($thumbprint, "My")
# 删除目录下的证书文件
Remove-Item -Path "$pfxPath" -Force
            `.trim();
            })
          );
        }
        break;

      case "binding_netsh":
        {
          setConfig(
            produce(config, (draft) => {
              draft.config ??= {};
              draft.config.shell = "powershell";
              draft.config.command = `
# 请将以下变量替换为实际值
$pfxPath = "<your-pfx-path>" # PFX 文件路径
$pfxPassword = "<your-pfx-password>" # PFX 密码
$ipaddr = "<your-binding-ip>"  # 绑定 IP，“0.0.0.0”表示所有 IP 绑定，可填入域名。
$port = "<your-binding-port>"  # 绑定端口

$addr = $ipaddr + ":" + $port

# 导入证书到本地计算机的个人存储区
$cert = Import-PfxCertificate -FilePath "$pfxPath" -CertStoreLocation Cert:\\LocalMachine\\My -Password (ConvertTo-SecureString -String "$pfxPassword" -AsPlainText -Force) -Exportable
# 获取 Thumbprint
$thumbprint = $cert.Thumbprint
# 检测端口是否绑定证书，如绑定则删除绑定
$isExist = netsh http show sslcert ipport=$addr
if ($isExist -like "*$addr*"){ netsh http delete sslcert ipport=$addr }
# 绑定到端口
netsh http add sslcert ipport=$addr certhash=$thumbprint
# 删除目录下的证书文件
Remove-Item -Path "$pfxPath" -Force
            `.trim();
            })
          );
        }
        break;
    }
  };

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
          <Label>{t("domain.deployment.form.shell.label")}</Label>
          <RadioGroup
            className="flex mt-1"
            value={config?.config?.shell}
            onValueChange={(val) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.shell = val;
              });
              setConfig(nv);
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
          <div className="text-red-600 text-sm mt-1">{errors?.shell}</div>
        </div>

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
          <div className="flex items-center justify-between">
            <Label>{t("domain.deployment.form.shell_command.label")}</Label>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <a className="text-xs text-blue-500 cursor-pointer">{t("domain.deployment.form.shell_preset_scripts.trigger")}</a>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem onClick={() => handleUsePresetScript("reload_nginx")}>
                  {t("domain.deployment.form.shell_preset_scripts.option.reload_nginx.label")}
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => handleUsePresetScript("binding_iis")}>
                  {t("domain.deployment.form.shell_preset_scripts.option.binding_iis.label")}
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => handleUsePresetScript("binding_netsh")}>
                  {t("domain.deployment.form.shell_preset_scripts.option.binding_netsh.label")}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
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

export default DeployToLocal;
