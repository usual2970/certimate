import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type NameDotComAccessConfig } from "@/domain/access";

type AccessEditFormNameDotComConfigModelType = Partial<NameDotComAccessConfig>;

export type AccessEditFormNameDotComConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  loading?: boolean;
  model?: AccessEditFormNameDotComConfigModelType;
  onModelChange?: (model: AccessEditFormNameDotComConfigModelType) => void;
};

const initModel = () => {
  return {
    username: "",
    apiToken: "",
  } as AccessEditFormNameDotComConfigModelType;
};

const AccessEditFormNameDotComConfig = ({ form, formName, disabled, loading, model, onModelChange }: AccessEditFormNameDotComConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    username: z
      .string()
      .trim()
      .min(1, t("access.form.namedotcom_username.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    apiToken: z
      .string()
      .trim()
      .min(1, t("access.form.namedotcom_api_token.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormNameDotComConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="username"
        label={t("access.form.namedotcom_username.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.namedotcom_username.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.namedotcom_username.placeholder")} />
      </Form.Item>

      <Form.Item
        name="apiToken"
        label={t("access.form.namedotcom_api_token.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.namedotcom_api_token.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.namedotcom_api_token.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormNameDotComConfig;
