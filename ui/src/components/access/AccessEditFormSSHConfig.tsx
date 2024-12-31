import { useState } from "react";
import { flushSync } from "react-dom";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Button, Form, Input, InputNumber, Upload, type FormInstance, type UploadFile, type UploadProps } from "antd";
import { UploadOutlined as UploadOutlinedIcon } from "@ant-design/icons";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type SSHAccessConfig } from "@/domain/access";
import { readFileContent } from "@/utils/file";
import { validDomainName, validIPv4Address, validIPv6Address } from "@/utils/validators";

type AccessEditFormSSHConfigFieldValues = Partial<SSHAccessConfig>;

export type AccessEditFormSSHConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormSSHConfigFieldValues;
  onValuesChange?: (values: AccessEditFormSSHConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormSSHConfigFieldValues => {
  return {
    host: "127.0.0.1",
    port: 22,
    username: "root",
  };
};

const AccessEditFormSSHConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormSSHConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    host: z.string().refine((v) => validDomainName(v) || validIPv4Address(v) || validIPv6Address(v), t("common.errmsg.host_invalid")),
    port: z.number().int().gte(1, t("common.errmsg.port_invalid")).lte(65535, t("common.errmsg.port_invalid")),
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
      .refine((v) => !v || form.getFieldValue("key"), t("access.form.ssh_key.placeholder")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const [keyFileList, setKeyFileList] = useState<UploadFile[]>([]);
  useDeepCompareEffect(() => {
    setKeyFileList(initialValues?.key?.trim() ? [{ uid: "-1", name: "sshkey", status: "done" }] : []);
  }, [initialValues]);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormSSHConfigFieldValues);
  };

  const handleKeyFileChange: UploadProps["onChange"] = async ({ file }) => {
    if (file && file.status !== "removed") {
      formInst.setFieldValue("key", await readFileContent(file.originFileObj ?? (file as unknown as File)));
      setKeyFileList([file]);
    } else {
      formInst.setFieldValue("key", "");
      setKeyFileList([]);
    }

    flushSync(() => onValuesChange?.(formInst.getFieldsValue(true)));
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
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

      <div className="flex space-x-2">
        <div className="w-1/2">
          <Form.Item name="username" label={t("access.form.ssh_username.label")} rules={[formRule]}>
            <Input autoComplete="new-password" placeholder={t("access.form.ssh_username.placeholder")} />
          </Form.Item>
        </div>

        <div className="w-1/2">
          <Form.Item
            name="password"
            label={t("access.form.ssh_password.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.ssh_password.tooltip") }}></span>}
          >
            <Input.Password autoComplete="new-password" placeholder={t("access.form.ssh_password.placeholder")} />
          </Form.Item>
        </div>
      </div>

      <div className="flex space-x-2">
        <div className="w-1/2">
          <Form.Item name="key" noStyle rules={[formRule]}>
            <Input.TextArea autoComplete="new-password" hidden placeholder={t("access.form.ssh_key.placeholder")} value={formInst.getFieldValue("key")} />
          </Form.Item>
          <Form.Item label={t("access.form.ssh_key.label")} tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.ssh_key.tooltip") }}></span>}>
            <Upload beforeUpload={() => false} fileList={keyFileList} maxCount={1} onChange={handleKeyFileChange}>
              <Button icon={<UploadOutlinedIcon />}>{t("access.form.ssh_key.upload")}</Button>
            </Upload>
          </Form.Item>
        </div>

        <div className="w-1/2">
          <Form.Item
            name="keyPassphrase"
            label={t("access.form.ssh_key_passphrase.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.ssh_key_passphrase.tooltip") }}></span>}
          >
            <Input.Password autoComplete="new-password" placeholder={t("access.form.ssh_key_passphrase.placeholder")} />
          </Form.Item>
        </div>
      </div>
    </Form>
  );
};

export default AccessEditFormSSHConfig;
