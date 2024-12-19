import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type QiniuAccessConfig } from "@/domain/access";

type AccessEditFormQiniuConfigModelType = Partial<QiniuAccessConfig>;

export type AccessEditFormQiniuConfigProps = {
  form: FormInstance;
  disabled?: boolean;
  loading?: boolean;
  model?: AccessEditFormQiniuConfigModelType;
  onModelChange?: (model: AccessEditFormQiniuConfigModelType) => void;
};

const initModel = () => {
  return {} as AccessEditFormQiniuConfigModelType;
};

const AccessEditFormQiniuConfig = ({ form, disabled, loading, model, onModelChange }: AccessEditFormQiniuConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKey: z
      .string()
      .trim()
      .min(1, t("access.form.qiniu_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretKey: z
      .string()
      .trim()
      .min(1, t("access.form.qiniu_secret_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormQiniuConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" name="configForm" onValuesChange={handleFormChange}>
      <Form.Item
        name="accessKey"
        label={t("access.form.qiniu_access_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.qiniu_access_key.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.qiniu_access_key.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretKey"
        label={t("access.form.qiniu_secret_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.qiniu_secret_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.qiniu_secret_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormQiniuConfig;
