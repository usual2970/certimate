import { useTranslation } from "react-i18next";
import { DownOutlined as DownOutlinedIcon } from "@ant-design/icons";
import { Button, Dropdown, Form, type FormInstance, Input, Select, Switch } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import Show from "@/components/Show";
import { CERTIFICATE_FORMATS } from "@/domain/certificate";

type DeployNodeConfigFormSSHConfigFieldValues = Nullish<{
  format: string;
  certPath: string;
  keyPath?: string | null;
  pfxPassword?: string | null;
  jksAlias?: string | null;
  jksKeypass?: string | null;
  jksStorepass?: string | null;
  preCommand?: string | null;
  postCommand?: string | null;
  useSCP?: boolean;
}>;

export type DeployNodeConfigFormSSHConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormSSHConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormSSHConfigFieldValues) => void;
};

const FORMAT_PEM = CERTIFICATE_FORMATS.PEM;
const FORMAT_PFX = CERTIFICATE_FORMATS.PFX;
const FORMAT_JKS = CERTIFICATE_FORMATS.JKS;

const initFormModel = (): DeployNodeConfigFormSSHConfigFieldValues => {
  return {
    format: FORMAT_PEM,
    certPath: "/etc/ssl/certimate/cert.crt",
    keyPath: "/etc/ssl/certimate/cert.key",
  };
};

const DeployNodeConfigFormSSHConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: DeployNodeConfigFormSSHConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    format: z.union([z.literal(FORMAT_PEM), z.literal(FORMAT_PFX), z.literal(FORMAT_JKS)], {
      message: t("workflow_node.deploy.form.ssh_format.placeholder"),
    }),
    certPath: z
      .string()
      .min(1, t("workflow_node.deploy.form.ssh_cert_path.tooltip"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim(),
    keyPath: z
      .string()
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_PEM || !!v?.trim(), { message: t("workflow_node.deploy.form.ssh_key_path.tooltip") }),
    pfxPassword: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_PFX || !!v?.trim(), { message: t("workflow_node.deploy.form.ssh_pfx_password.tooltip") }),
    jksAlias: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_JKS || !!v?.trim(), { message: t("workflow_node.deploy.form.ssh_jks_alias.tooltip") }),
    jksKeypass: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_JKS || !!v?.trim(), { message: t("workflow_node.deploy.form.ssh_jks_keypass.tooltip") }),
    jksStorepass: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish()
      .refine((v) => fieldFormat !== FORMAT_JKS || !!v?.trim(), { message: t("workflow_node.deploy.form.ssh_jks_storepass.tooltip") }),
    preCommand: z
      .string()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
    postCommand: z
      .string()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
    useSCP: z.boolean().nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

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

  const handlePresetPreScriptClick = (key: string) => {
    switch (key) {
      case "backup_files":
        {
          formInst.setFieldValue(
            "preCommand",
            `# 请将以下路径替换为实际值
cp "${formInst.getFieldValue("certPath")}" "${formInst.getFieldValue("certPath")}.bak" 2>/dev/null || :
cp "${formInst.getFieldValue("keyPath")}" "${formInst.getFieldValue("keyPath")}.bak" 2>/dev/null || :
            `.trim()
          );
        }
        break;
    }
  };

  const handlePresetPostScriptClick = (key: string) => {
    switch (key) {
      case "reload_nginx":
        {
          formInst.setFieldValue("postCommand", "sudo service nginx reload");
        }
        break;
    }
  };

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values);
  };

  return (
    <Form
      form={formInst}
      disabled={disabled}
      initialValues={initialValues ?? initFormModel()}
      layout="vertical"
      name={formName}
      onValuesChange={handleFormChange}
    >
      <Form.Item name="format" label={t("workflow_node.deploy.form.ssh_format.label")} rules={[formRule]}>
        <Select placeholder={t("workflow_node.deploy.form.ssh_format.placeholder")} onSelect={handleFormatSelect}>
          <Select.Option key={FORMAT_PEM} value={FORMAT_PEM}>
            {t("workflow_node.deploy.form.ssh_format.option.pem.label")}
          </Select.Option>
          <Select.Option key={FORMAT_PFX} value={FORMAT_PFX}>
            {t("workflow_node.deploy.form.ssh_format.option.pfx.label")}
          </Select.Option>
          <Select.Option key={FORMAT_JKS} value={FORMAT_JKS}>
            {t("workflow_node.deploy.form.ssh_format.option.jks.label")}
          </Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        name="certPath"
        label={t("workflow_node.deploy.form.ssh_cert_path.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.ssh_cert_path.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.ssh_cert_path.placeholder")} />
      </Form.Item>

      <Show when={fieldFormat === FORMAT_PEM}>
        <Form.Item
          name="keyPath"
          label={t("workflow_node.deploy.form.ssh_key_path.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.ssh_key_path.tooltip") }}></span>}
        >
          <Input placeholder={t("workflow_node.deploy.form.ssh_key_path.placeholder")} />
        </Form.Item>
      </Show>

      <Show when={fieldFormat === FORMAT_PFX}>
        <Form.Item
          name="pfxPassword"
          label={t("workflow_node.deploy.form.ssh_pfx_password.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.ssh_pfx_password.tooltip") }}></span>}
        >
          <Input placeholder={t("workflow_node.deploy.form.ssh_pfx_password.placeholder")} />
        </Form.Item>
      </Show>

      <Show when={fieldFormat === FORMAT_JKS}>
        <Form.Item
          name="jksAlias"
          label={t("workflow_node.deploy.form.ssh_jks_alias.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.ssh_jks_alias.tooltip") }}></span>}
        >
          <Input placeholder={t("workflow_node.deploy.form.ssh_jks_alias.placeholder")} />
        </Form.Item>

        <Form.Item
          name="jksKeypass"
          label={t("workflow_node.deploy.form.ssh_jks_keypass.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.ssh_jks_keypass.tooltip") }}></span>}
        >
          <Input placeholder={t("workflow_node.deploy.form.ssh_jks_keypass.placeholder")} />
        </Form.Item>

        <Form.Item
          name="jksStorepass"
          label={t("workflow_node.deploy.form.ssh_jks_storepass.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.ssh_jks_storepass.tooltip") }}></span>}
        >
          <Input placeholder={t("workflow_node.deploy.form.ssh_jks_storepass.placeholder")} />
        </Form.Item>
      </Show>

      <Form.Item label={t("workflow_node.deploy.form.ssh_shell_env.label")}>
        <Select options={[{ value: t("workflow_node.deploy.form.ssh_shell_env.value") }]} value={t("workflow_node.deploy.form.ssh_shell_env.value")} />
      </Form.Item>

      <Form.Item className="mb-0">
        <label className="mb-1 block">
          <div className="flex w-full items-center justify-between gap-4">
            <div className="max-w-full grow truncate">
              <span>{t("workflow_node.deploy.form.ssh_pre_command.label")}</span>
            </div>
            <div className="text-right">
              <Dropdown
                menu={{
                  items: [
                    {
                      key: "backup_files",
                      label: t("workflow_node.deploy.form.ssh_preset_scripts.option.backup_files.label"),
                      onClick: () => handlePresetPreScriptClick("backup_files"),
                    },
                  ],
                }}
                trigger={["click"]}
              >
                <Button size="small" type="link">
                  {t("workflow_node.deploy.form.ssh_preset_scripts.button")}
                  <DownOutlinedIcon />
                </Button>
              </Dropdown>
            </div>
          </div>
        </label>
        <Form.Item name="preCommand" rules={[formRule]}>
          <Input.TextArea autoSize={{ minRows: 1, maxRows: 5 }} placeholder={t("workflow_node.deploy.form.ssh_pre_command.placeholder")} />
        </Form.Item>
      </Form.Item>

      <Form.Item className="mb-0">
        <label className="mb-1 block">
          <div className="flex w-full items-center justify-between gap-4">
            <div className="max-w-full grow truncate">
              <span>{t("workflow_node.deploy.form.ssh_post_command.label")}</span>
            </div>
            <div className="text-right">
              <Dropdown
                menu={{
                  items: [
                    {
                      key: "reload_nginx",
                      label: t("workflow_node.deploy.form.ssh_preset_scripts.option.reload_nginx.label"),
                      onClick: () => handlePresetPostScriptClick("reload_nginx"),
                    },
                  ],
                }}
                trigger={["click"]}
              >
                <Button size="small" type="link">
                  {t("workflow_node.deploy.form.ssh_preset_scripts.button")}
                  <DownOutlinedIcon />
                </Button>
              </Dropdown>
            </div>
          </div>
        </label>
        <Form.Item name="postCommand" rules={[formRule]}>
          <Input.TextArea autoSize={{ minRows: 1, maxRows: 5 }} placeholder={t("workflow_node.deploy.form.ssh_post_command.placeholder")} />
        </Form.Item>
      </Form.Item>

      <Form.Item
        name="useSCP"
        label={t("workflow_node.deploy.form.ssh_use_scp.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.ssh_use_scp.tooltip") }}></span>}
      >
        <Switch />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigFormSSHConfig;
