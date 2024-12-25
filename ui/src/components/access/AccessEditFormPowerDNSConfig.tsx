import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type PowerDNSAccessConfig } from "@/domain/access";

type AccessEditFormPowerDNSConfigModelValues = Partial<PowerDNSAccessConfig>;

export type AccessEditFormPowerDNSConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormPowerDNSConfigModelValues;
  onModelChange?: (model: AccessEditFormPowerDNSConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormPowerDNSConfigModelValues => {
  return {
    apiUrl: "",
    apiKey: "",
  };
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
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onModelChange?.(values as AccessEditFormPowerDNSConfigModelValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
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
