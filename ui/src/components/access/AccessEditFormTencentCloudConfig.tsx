import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type TencentCloudAccessConfig } from "@/domain/access";

type AccessEditFormTencentCloudConfigModelType = Partial<TencentCloudAccessConfig>;

export type AccessEditFormTencentCloudConfigProps = {
  form: FormInstance;
  disabled?: boolean;
  loading?: boolean;
  model?: AccessEditFormTencentCloudConfigModelType;
  onModelChange?: (model: AccessEditFormTencentCloudConfigModelType) => void;
};

const initModel = () => {
  return {} as AccessEditFormTencentCloudConfigModelType;
};

const AccessEditFormTencentCloudConfig = ({ form, disabled, loading, model, onModelChange }: AccessEditFormTencentCloudConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    secretId: z
      .string()
      .trim()
      .min(1, t("access.form.tencentcloud_secret_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretKey: z
      .string()
      .trim()
      .min(1, t("access.form.tencentcloud_secret_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormTencentCloudConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" name="configForm" onValuesChange={handleFormChange}>
      <Form.Item
        name="secretId"
        label={t("access.form.tencentcloud_secret_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.tencentcloud_secret_id.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.tencentcloud_secret_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretKey"
        label={t("access.form.tencentcloud_secret_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.tencentcloud_secret_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.tencentcloud_secret_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormTencentCloudConfig;
