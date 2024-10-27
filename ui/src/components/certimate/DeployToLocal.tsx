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

const DeployToLocal = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          format: "pem",
          certPath: "/etc/nginx/ssl/nginx.crt",
          keyPath: "/etc/nginx/ssl/nginx.key",
          pfxPassword: "",
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

  const formSchema = z
    .object({
      format: z.union([z.literal("pem"), z.literal("pfx")], {
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

  useEffect(() => {
    if (data.config?.format === "pem") {
      if (data.config.certPath && data.config.certPath.endsWith(".pfx")) {
        const newData = produce(data, (draft) => {
          draft.config ??= {};
          draft.config.certPath = data.config!.certPath.replace(/.pfx$/, ".crt");
        });
        setDeploy(newData);
      }
    } else if (data.config?.format === "pfx") {
      if (data.config.certPath && data.config.certPath.endsWith(".crt")) {
        const newData = produce(data, (draft) => {
          draft.config ??= {};
          draft.config.certPath = data.config!.certPath.replace(/.crt$/, ".pfx");
        });
        setDeploy(newData);
      }
    }
  }, [data.config?.format]);

  const getOptionCls = (val: string) => {
    if (data.config?.shell === val) {
      return "border-primary dark:border-primary";
    }

    return "";
  };

  const handleUsePresetScript = (key: string) => {
    switch (key) {
      case "reload_nginx":
        {
          const newData = produce(data, (draft) => {
            draft.config ??= {};
            draft.config.shell = "sh";
            draft.config.command = "sudo service nginx reload";
          });
          setDeploy(newData);
        }
        break;

      case "binding_iis":
        {
          const newData = produce(data, (draft) => {
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
          });
          setDeploy(newData);
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
            value={data?.config?.format}
            onValueChange={(value) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.format = value;
              });
              setDeploy(newData);
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("domain.deployment.form.file_format.placeholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="pem">PEM</SelectItem>
                <SelectItem value="pfx">PFX</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          <div className="text-red-600 text-sm mt-1">{error?.format}</div>
        </div>

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

        {data.config?.format === "pem" ? (
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
        ) : (
          <></>
        )}

        {data.config?.format === "pfx" ? (
          <div>
            <Label>{t("domain.deployment.form.file_pfx_password.label")}</Label>
            <Input
              placeholder={t("domain.deployment.form.file_pfx_password.placeholder")}
              className="w-full mt-1"
              value={data?.config?.pfxPassword}
              onChange={(e) => {
                const newData = produce(data, (draft) => {
                  draft.config ??= {};
                  draft.config.pfxPassword = e.target.value?.trim();
                });
                setDeploy(newData);
              }}
            />
            <div className="text-red-600 text-sm mt-1">{error?.pfxPassword}</div>
          </div>
        ) : (
          <></>
        )}

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
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
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
