import { useTranslation } from "react-i18next";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "../ui/select";
import { Button } from "../ui/button";
import { DeployFormProps } from "./DeployForm";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useShallow } from "zustand/shallow";
import { usePanel } from "./PanelProvider";
import { useEffect, useState } from "react";
import i18n from "@/i18n";
import { WorkflowNode } from "@/domain/workflow";
import { Textarea } from "../ui/textarea";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";
import AccessSelect from "./AccessSelect";
import AccessEditDialog from "../certimate/AccessEditDialog";
import { Plus } from "lucide-react";

const selectState = (state: WorkflowState) => ({
  updateNode: state.updateNode,
  getWorkflowOuptutBeforeId: state.getWorkflowOuptutBeforeId,
});

const t = i18n.t;

const formSchema = z
  .object({
    providerType: z.string(),
    access: z.string().min(1, t("domain.deployment.form.access.placeholder")),
    certificate: z.string().min(1),
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
    shell: z.union([z.literal("sh"), z.literal("cmd"), z.literal("powershell")], {
      message: t("domain.deployment.form.shell.placeholder"),
    }),
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

const DeployToLocal = ({ data }: DeployFormProps) => {
  const { updateNode, getWorkflowOuptutBeforeId } = useWorkflowStore(useShallow(selectState));
  const { hidePanel } = usePanel();
  const { t } = useTranslation();

  const [beforeOutput, setBeforeOutput] = useState<WorkflowNode[]>([]);

  useEffect(() => {
    const rs = getWorkflowOuptutBeforeId(data.id, "certificate");
    console.log(rs);
    setBeforeOutput(rs);
  }, [data]);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      providerType: "local",
      access: data.config?.access as string,
      certificate: data.config?.certificate as string,
      format: (data.config?.format as "pem" | "pfx" | "jks") || "pem",
      certPath: (data.config?.certPath as string) || "/etc/ssl/certs/cert.crt",
      keyPath: (data.config?.keyPath as string) || "/etc/ssl/private/cert.key",
      pfxPassword: (data.config?.pfxPassword as string) || "",
      jksAlias: (data.config?.jksAlias as string) || "",
      jksKeypass: (data.config?.jksKeypass as string) || "",
      jksStorepass: (data.config?.jksStorepass as string) || "",
      preCommand: (data.config?.preCommand as string) || "",
      command: (data.config?.command as string) || "service nginx reload",
      shell: (data.config?.shell as "sh" | "cmd" | "powershell") || "sh",
    },
  });

  const format = form.watch("format");
  const certPath = form.watch("certPath");

  useEffect(() => {
    if (format === "pem" && /(.pfx|.jks)$/.test(certPath)) {
      form.setValue("certPath", certPath.replace(/(.pfx|.jks)$/, ".crt"));
    } else if (format === "pfx" && /(.crt|.jks)$/.test(certPath)) {
      form.setValue("certPath", certPath.replace(/(.crt|.jks)$/, ".pfx"));
    } else if (format === "jks" && /(.crt|.pfx)$/.test(certPath)) {
      form.setValue("certPath", certPath.replace(/(.crt|.pfx)$/, ".jks"));
    }
  }, [format]);

  const onSubmit = async (config: z.infer<typeof formSchema>) => {
    updateNode({ ...data, config, validated: true });
    hidePanel();
  };

  const handleUsePresetScript = (key: string) => {
    switch (key) {
      case "reload_nginx":
        {
          form.setValue("shell", "sh");
          form.setValue("command", "sudo service nginx reload");
        }
        break;

      case "binding_iis":
        {
          form.setValue("shell", "powershell");
          form.setValue(
            "command",
            `# 请将以下变量替换为实际值
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
            `.trim()
          );
        }
        break;

      case "binding_netsh":
        {
          form.setValue("shell", "powershell");
          form.setValue(
            "command",
            `# 请将以下变量替换为实际值
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
            `.trim()
          );
        }
        break;
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="access"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="flex justify-between">
                <div>{t("domain.deployment.form.access.label")}</div>

                <AccessEditDialog
                  trigger={
                    <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                      <Plus size={14} />
                      {t("common.add")}
                    </div>
                  }
                  op="add"
                  outConfigType="local"
                />
              </FormLabel>
              <FormControl>
                <AccessSelect
                  {...field}
                  value={field.value}
                  onValueChange={(value) => {
                    form.setValue("access", value);
                  }}
                  providerType="local"
                />
              </FormControl>

              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="certificate"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("workflow.common.certificate.label")}</FormLabel>
              <FormControl>
                <Select
                  {...field}
                  value={field.value}
                  onValueChange={(value) => {
                    form.setValue("certificate", value);
                  }}
                >
                  <SelectTrigger>
                    <SelectValue placeholder={t("workflow.common.certificate.placeholder")} />
                  </SelectTrigger>
                  <SelectContent>
                    {beforeOutput.map((item) => (
                      <>
                        <SelectGroup key={item.id}>
                          <SelectLabel>{item.name}</SelectLabel>
                          {item.output?.map((output) => (
                            <SelectItem key={output.name} value={`${item.id}#${output.name}`}>
                              <div>
                                {item.name}-{output.label}
                              </div>
                            </SelectItem>
                          ))}
                        </SelectGroup>
                      </>
                    ))}
                  </SelectContent>
                </Select>
              </FormControl>

              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="format"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("domain.deployment.form.file_format.label")}</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t("domain.deployment.form.file_format.placeholder")} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="pem">PEM</SelectItem>
                  <SelectItem value="pfx">PFX</SelectItem>
                  <SelectItem value="jks">JKS</SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="certPath"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("domain.deployment.form.file_cert_path.label")}</FormLabel>
              <FormControl>
                <Input placeholder={t("domain.deployment.form.file_cert_path.label")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="keyPath"
          render={({ field }) => (
            <FormItem>
              <FormLabel>密钥路径</FormLabel>
              <FormControl>
                <Input placeholder="输入密钥路径" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {format === "pfx" && (
          <FormField
            control={form.control}
            name="pfxPassword"
            render={({ field }) => (
              <FormItem>
                <FormLabel>PFX 密码</FormLabel>
                <FormControl>
                  <Input type="password" placeholder="输入 PFX 密码" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        )}

        {format === "jks" && (
          <>
            <FormField
              control={form.control}
              name="jksAlias"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>JKS 别名</FormLabel>
                  <FormControl>
                    <Input placeholder="输入 JKS 别名" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="jksKeypass"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>JKS Keypass</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="输入 JKS Keypass" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="jksStorepass"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>JKS Storepass</FormLabel>
                  <FormControl>
                    <Input type="password" placeholder="输入 JKS Storepass" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </>
        )}

        <FormField
          control={form.control}
          name="shell"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("domain.deployment.form.shell.label")}</FormLabel>
              <Select onValueChange={field.onChange} value={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t("domain.deployment.form.shell.placeholder")} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="sh">POSIX Bash (Linux)</SelectItem>
                  <SelectItem value="cmd">CMD (Windows)</SelectItem>
                  <SelectItem value="powershell">PowerShell (Windows)</SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="preCommand"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("domain.deployment.form.shell_pre_command.label")}</FormLabel>
              <FormControl>
                <Textarea placeholder={t("domain.deployment.form.shell_pre_command.placeholder")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="command"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="flex justify-between items-center">
                <div>{t("domain.deployment.form.shell_command.label")}</div>
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
              </FormLabel>
              <FormControl>
                <Textarea placeholder={t("domain.deployment.form.shell_command.placeholder")} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex justify-end">
          <Button type="submit">{t("common.save")}</Button>
        </div>
      </form>
    </Form>
  );
};

export default DeployToLocal;
