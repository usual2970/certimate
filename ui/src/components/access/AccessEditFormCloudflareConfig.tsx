import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type CloudflareAccessConfig } from "@/domain/access";

type AccessEditFormCloudflareConfigModelType = Partial<CloudflareAccessConfig>;

export type AccessEditFormCloudflareConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  loading?: boolean;
  model?: AccessEditFormCloudflareConfigModelType;
  onModelChange?: (model: AccessEditFormCloudflareConfigModelType) => void;
};

const initModel = () => {
  return {} as AccessEditFormCloudflareConfigModelType;
};

const AccessEditFormCloudflareConfig = ({ form, formName, disabled, loading, model, onModelChange }: AccessEditFormCloudflareConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    dnsApiToken: z
      .string()
      .trim()
      .min(1, t("access.form.cloudflare_dns_api_token.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormCloudflareConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="dnsApiToken"
        label={t("access.form.cloudflare_dns_api_token.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.cloudflare_dns_api_token.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.cloudflare_dns_api_token.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormCloudflareConfig;
