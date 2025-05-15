import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input, InputNumber } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import TextFileInput from "@/components/TextFileInput";
import { type AccessConfigForSSH } from "@/domain/access";
import { validDomainName, validIPv4Address, validIPv6Address, validPortNumber } from "@/utils/validators";

type AccessFormSSHConfigFieldValues = Nullish<AccessConfigForSSH>;

export type AccessFormSSHConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormSSHConfigFieldValues;
  onValuesChange?: (values: AccessFormSSHConfigFieldValues) => void;
};

const initFormModel = (): AccessFormSSHConfigFieldValues => {
  return {
    host: "127.0.0.1",
    port: 22,
    username: "root",
  };
};

const AccessFormSSHConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormSSHConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    host: z.string().refine((v) => validDomainName(v) || validIPv4Address(v) || validIPv6Address(v), t("common.errmsg.host_invalid")),
    port: z.preprocess(
      (v) => Number(v),
      z
        .number()
        .int(t("access.form.ssh_port.placeholder"))
        .refine((v) => validPortNumber(v), t("common.errmsg.port_invalid"))
    ),
    username: z
      .string()
      .min(1, "access.form.ssh_username.placeholder")
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    password: z
      .string()
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .nullish(),
    key: z
      .string()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
    keyPassphrase: z
      .string()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish()
      .refine((v) => !v || formInst.getFieldValue("key"), t("access.form.ssh_key.placeholder")),
  });
  const formRule = createSchemaFieldRule(formSchema);

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
      <div className="flex space-x-2">
        <div className="w-2/3">
          <Form.Item name="host" label={t("access.form.ssh_host.label")} rules={[formRule]}>
            <Input placeholder={t("access.form.ssh_host.placeholder")} />
          </Form.Item>
        </div>

        <div className="w-1/3">
          <Form.Item name="port" label={t("access.form.ssh_port.label")} rules={[formRule]}>
            <InputNumber className="w-full" placeholder={t("access.form.ssh_port.placeholder")} min={1} max={65535} />
          </Form.Item>
        </div>
      </div>

      <Form.Item name="username" label={t("access.form.ssh_username.label")} rules={[formRule]}>
        <Input autoComplete="new-password" placeholder={t("access.form.ssh_username.placeholder")} />
      </Form.Item>

      <Form.Item
        name="password"
        label={t("access.form.ssh_password.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.ssh_password.tooltip") }}></span>}
      >
        <Input.Password allowClear autoComplete="new-password" placeholder={t("access.form.ssh_password.placeholder")} />
      </Form.Item>

      <Form.Item
        name="key"
        label={t("access.form.ssh_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.ssh_key.tooltip") }}></span>}
      >
        <TextFileInput allowClear autoSize={{ minRows: 1, maxRows: 5 }} placeholder={t("access.form.ssh_key.placeholder")} />
      </Form.Item>

      <Form.Item
        name="keyPassphrase"
        label={t("access.form.ssh_key_passphrase.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.ssh_key_passphrase.tooltip") }}></span>}
      >
        <Input.Password allowClear autoComplete="new-password" placeholder={t("access.form.ssh_key_passphrase.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessFormSSHConfig;
