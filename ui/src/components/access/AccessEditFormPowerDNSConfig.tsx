import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type PowerDNSAccessConfig } from "@/domain/access";

type AccessEditFormPowerDNSConfigModelType = Partial<PowerDNSAccessConfig>;

export type AccessEditFormPowerDNSConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormPowerDNSConfigModelType;
  onModelChange?: (model: AccessEditFormPowerDNSConfigModelType) => void;
};

const initModel = () => {
  return {
    apiUrl: "",
    apiKey: "",
  } as AccessEditFormPowerDNSConfigModelType;
};

const AccessEditFormPowerDNSConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormPowerDNSConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    apiUrl: z.string().url(t("common.errmsg.url_invalid")),
    apiKey: z
      .string()
      .trim()
      .min(1, t("access.form.powerdns_api_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormPowerDNSConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="apiUrl"
        label={t("access.form.powerdns_api_url.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.powerdns_api_url.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.powerdns_api_url.placeholder")} />
      </Form.Item>

      <Form.Item
        name="apiKey"
        label={t("access.form.powerdns_api_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.powerdns_api_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.powerdns_api_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormPowerDNSConfig;
