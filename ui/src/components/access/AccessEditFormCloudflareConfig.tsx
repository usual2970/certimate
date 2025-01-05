import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForCloudflare } from "@/domain/access";

type AccessEditFormCloudflareConfigFieldValues = Partial<AccessConfigForCloudflare>;

export type AccessEditFormCloudflareConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormCloudflareConfigFieldValues;
  onValuesChange?: (values: AccessEditFormCloudflareConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormCloudflareConfigFieldValues => {
  return {
    dnsApiToken: "",
  };
};

const AccessEditFormCloudflareConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormCloudflareConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    dnsApiToken: z
      .string()
      .min(1, t("access.form.cloudflare_dns_api_token.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormCloudflareConfigFieldValues);
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues ?? initFormModel()} layout="vertical" name={formName} onValuesChange={handleFormChange}>
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
