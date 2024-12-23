import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type GoDaddyAccessConfig } from "@/domain/access";

type AccessEditFormGoDaddyConfigModelType = Partial<GoDaddyAccessConfig>;

export type AccessEditFormGoDaddyConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormGoDaddyConfigModelType;
  onModelChange?: (model: AccessEditFormGoDaddyConfigModelType) => void;
};

const initModel = () => {
  return {
    apiKey: "",
    apiSecret: "",
  } as AccessEditFormGoDaddyConfigModelType;
};

const AccessEditFormGoDaddyConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormGoDaddyConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    apiKey: z
      .string()
      .trim()
      .min(1, t("access.form.godaddy_api_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    apiSecret: z
      .string()
      .trim()
      .min(1, t("access.form.godaddy_api_secret.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormGoDaddyConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="apiKey"
        label={t("access.form.godaddy_api_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.godaddy_api_key.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.godaddy_api_key.placeholder")} />
      </Form.Item>

      <Form.Item
        name="apiSecret"
        label={t("access.form.godaddy_api_secret.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.godaddy_api_secret.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.godaddy_api_secret.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormGoDaddyConfig;
