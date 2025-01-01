import { useTranslation } from "react-i18next";
import { Button, Dropdown, Form, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { DownOutlined as DownOutlinedIcon } from "@ant-design/icons";
import { z } from "zod";

import Show from "@/components/Show";

const FORMAT_PEM = "pem" as const;
const FORMAT_PFX = "pfx" as const;
const FORMAT_JKS = "jks" as const;

const SHELLENV_SH = "sh" as const;
const SHELLENV_CMD = "cmd" as const;
const SHELLENV_POWERSHELL = "powershell" as const;

const DeployNodeFormLocalFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    format: z.union([z.literal(FORMAT_PEM), z.literal(FORMAT_PFX), z.literal(FORMAT_JKS)], {
      message: t("workflow_node.deploy.form.local_format.placeholder"),
    }),
    certPath: z
      .string()
      .min(1, t("workflow_node.deploy.form.local_cert_path.tooltip"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim(),
    keyPath: z
      .string()
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_PEM || !!v?.trim(), { message: t("workflow_node.deploy.form.local_key_path.tooltip") }),
    pfxPassword: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_PFX || !!v?.trim(), { message: t("workflow_node.deploy.form.local_pfx_password.tooltip") }),
    jksAlias: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_JKS || !!v?.trim(), { message: t("workflow_node.deploy.form.local_jks_alias.tooltip") }),
    jksKeypass: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_JKS || !!v?.trim(), { message: t("workflow_node.deploy.form.local_jks_keypass.tooltip") }),
    jksStorepass: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_JKS || !!v?.trim(), { message: t("workflow_node.deploy.form.local_jks_storepass.tooltip") }),
    shellEnv: z.union([z.literal(SHELLENV_SH), z.literal(SHELLENV_CMD), z.literal(SHELLENV_POWERSHELL)], {
      message: t("workflow_node.deploy.form.local_shell_env.placeholder"),
    }),
    preCommand: z
      .string()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
    postCommand: z
      .string()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const formInst = Form.useFormInstance();

  const fieldFormat = Form.useWatch("format", formInst);
  const fieldCertPath = Form.useWatch("certPath", formInst);

  const handleFormatSelect = (value: string) => {
    if (fieldFormat === value) return;

    switch (value) {
      case FORMAT_PEM:
        {
          if (/(.pfx|.jks)$/.test(fieldCertPath)) {
            formInst.setFieldValue("certPath", fieldCertPath.replace(/(.pfx|.jks)$/, ".crt"));
          }
        }
        break;

      case FORMAT_PFX:
        {
          if (/(.crt|.jks)$/.test(fieldCertPath)) {
            formInst.setFieldValue("certPath", fieldCertPath.replace(/(.crt|.jks)$/, ".pfx"));
          }
        }
        break;

      case FORMAT_JKS:
        {
          if (/(.crt|.pfx)$/.test(fieldCertPath)) {
            formInst.setFieldValue("certPath", fieldCertPath.replace(/(.crt|.pfx)$/, ".jks"));
          }
        }
        break;
    }
  };

  const handlePresetScriptClick = (key: string) => {
    switch (key) {
      case "reload_nginx":
        {
          formInst.setFieldValue("shellEnv", "sh");
          formInst.setFieldValue("postCommand", "sudo service nginx reload");
        }
        break;

      case "binding_iis":
        {
          formInst.setFieldValue("shellEnv", "powershell");
          formInst.setFieldValue(
            "postCommand",
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
          formInst.setFieldValue("shellEnv", "powershell");
          formInst.setFieldValue(
            "postCommand",
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
    <>
      <Form.Item name="format" label={t("workflow_node.deploy.form.local_format.label")} rules={[formRule]} initialValue={FORMAT_PEM}>
        <Select placeholder={t("workflow_node.deploy.form.local_format.placeholder")} onSelect={handleFormatSelect}>
          <Select.Option key={FORMAT_PEM} value={FORMAT_PEM}>
            {t("workflow_node.deploy.form.local_format.option.pem.label")}
          </Select.Option>
          <Select.Option key={FORMAT_PFX} value={FORMAT_PFX}>
            {t("workflow_node.deploy.form.local_format.option.pfx.label")}
          </Select.Option>
          <Select.Option key={FORMAT_JKS} value={FORMAT_JKS}>
            {t("workflow_node.deploy.form.local_format.option.jks.label")}
          </Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        name="certPath"
        label={t("workflow_node.deploy.form.local_cert_path.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.local_cert_path.tooltip") }}></span>}
        initialValue="/etc/ssl/certs/cert.crt"
      >
        <Input placeholder={t("workflow_node.deploy.form.local_cert_path.placeholder")} />
      </Form.Item>

      <Show when={fieldFormat === FORMAT_PEM}>
        <Form.Item
          name="keyPath"
          label={t("workflow_node.deploy.form.local_key_path.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.local_key_path.tooltip") }}></span>}
          initialValue="/etc/ssl/certs/cert.key"
        >
          <Input placeholder={t("workflow_node.deploy.form.local_key_path.placeholder")} />
        </Form.Item>
      </Show>

      <Show when={fieldFormat === FORMAT_PFX}>
        <Form.Item name="pfxPassword" label={t("workflow_node.deploy.form.local_pfx_password.label")} rules={[formRule]}>
          <Input placeholder={t("workflow_node.deploy.form.local_pfx_password.placeholder")} />
        </Form.Item>
      </Show>

      <Show when={fieldFormat === FORMAT_JKS}>
        <Form.Item name="jksAlias" label={t("workflow_node.deploy.form.local_jks_alias.label")} rules={[formRule]}>
          <Input placeholder={t("workflow_node.deploy.form.local_jks_alias.placeholder")} />
        </Form.Item>

        <Form.Item name="jksKeypass" label={t("workflow_node.deploy.form.local_jks_keypass.label")} rules={[formRule]}>
          <Input placeholder={t("workflow_node.deploy.form.local_jks_keypass.placeholder")} />
        </Form.Item>

        <Form.Item name="jksStorepass" label={t("workflow_node.deploy.form.local_jks_storepass.label")} rules={[formRule]}>
          <Input placeholder={t("workflow_node.deploy.form.local_jks_storepass.placeholder")} />
        </Form.Item>
      </Show>

      <Form.Item name="shellEnv" label={t("workflow_node.deploy.form.local_shell_env.label")} rules={[formRule]} initialValue={SHELLENV_SH}>
        <Select placeholder={t("workflow_node.deploy.form.local_shell_env.placeholder")}>
          <Select.Option key={SHELLENV_SH} value={SHELLENV_SH}>
            {t("workflow_node.deploy.form.local_shell_env.option.sh.label")}
          </Select.Option>
          <Select.Option key={SHELLENV_CMD} value={SHELLENV_CMD}>
            {t("workflow_node.deploy.form.local_shell_env.option.cmd.label")}
          </Select.Option>
          <Select.Option key={SHELLENV_POWERSHELL} value={SHELLENV_POWERSHELL}>
            {t("workflow_node.deploy.form.local_shell_env.option.powershell.label")}
          </Select.Option>
        </Select>
      </Form.Item>

      <Form.Item name="preCommand" label={t("workflow_node.deploy.form.local_pre_command.label")} rules={[formRule]}>
        <Input.TextArea autoSize={{ minRows: 1, maxRows: 5 }} placeholder={t("workflow_node.deploy.form.local_pre_command.placeholder")} />
      </Form.Item>

      <Form.Item className="mb-0">
        <label className="block mb-1">
          <div className="flex items-center justify-between gap-4 w-full">
            <div className="flex-grow max-w-full truncate">
              <span>{t("workflow_node.deploy.form.local_post_command.label")}</span>
            </div>
            <div className="text-right">
              <Dropdown
                menu={{
                  items: [
                    {
                      key: "reload_nginx",
                      label: t("workflow_node.deploy.form.local_preset_scripts.option.reload_nginx.label"),
                      onClick: () => handlePresetScriptClick("reload_nginx"),
                    },
                    {
                      key: "binding_iis",
                      label: t("workflow_node.deploy.form.local_preset_scripts.option.binding_iis.label"),
                      onClick: () => handlePresetScriptClick("binding_iis"),
                    },
                    {
                      key: "binding_netsh",
                      label: t("workflow_node.deploy.form.local_preset_scripts.option.binding_netsh.label"),
                      onClick: () => handlePresetScriptClick("binding_netsh"),
                    },
                  ],
                }}
                trigger={["click"]}
              >
                <Button size="small" type="link">
                  {t("workflow_node.deploy.form.local_preset_scripts.button")}
                  <DownOutlinedIcon />
                </Button>
              </Dropdown>
            </div>
          </div>
        </label>
        <Form.Item name="postCommand" rules={[formRule]}>
          <Input.TextArea autoSize={{ minRows: 1, maxRows: 5 }} placeholder={t("workflow_node.deploy.form.local_post_command.placeholder")} />
        </Form.Item>
      </Form.Item>
    </>
  );
};

export default DeployNodeFormLocalFields;
