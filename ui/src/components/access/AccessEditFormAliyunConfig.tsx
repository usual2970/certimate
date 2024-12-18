import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AliyunAccessConfig } from "@/domain/access";

type AccessEditFormAliyunConfigModelType = Partial<AliyunAccessConfig>;

export type AccessEditFormAliyunConfigProps = {
  form: FormInstance;
  disabled?: boolean;
  loading?: boolean;
  model?: AccessEditFormAliyunConfigModelType;
  onModelChange?: (model: AccessEditFormAliyunConfigModelType) => void;
};

const initModel = () => {
  return {} as AccessEditFormAliyunConfigModelType;
};

const AccessEditFormAliyunConfig = ({ form, disabled, loading, model, onModelChange }: AccessEditFormAliyunConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKeyId: z
      .string()
      .trim()
      .min(1, t("access.form.aliyun_access_key_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    accessKeySecret: z
      .string()
      .trim()
      .min(1, t("access.form.aliyun_access_key_secret.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormAliyunConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" name="configForm" onValuesChange={handleFormChange}>
      <Form.Item
        name="accessKeyId"
        label={t("access.form.aliyun_access_key_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aliyun_access_key_id.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.aliyun_access_key_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="accessKeySecret"
        label={t("access.form.aliyun_access_key_secret.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aliyun_access_key_secret.tooltip") }}></span>}
      >
        <Input.Password placeholder={t("access.form.aliyun_access_key_secret.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormAliyunConfig;
