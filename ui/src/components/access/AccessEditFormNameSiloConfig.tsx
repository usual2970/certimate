import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type NameSiloAccessConfig } from "@/domain/access";

type AccessEditFormNameSiloConfigModelType = Partial<NameSiloAccessConfig>;

export type AccessEditFormNameSiloConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormNameSiloConfigModelType;
  onModelChange?: (model: AccessEditFormNameSiloConfigModelType) => void;
};

const initModel = () => {
  return {
    apiKey: "",
  } as AccessEditFormNameSiloConfigModelType;
};

const AccessEditFormNameSiloConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormNameSiloConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    apiKey: z
      .string()
      .trim()
      .min(1, t("access.form.namesilo_api_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormNameSiloConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="apiKey"
        label={t("access.form.namesilo_api_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.namesilo_api_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.namesilo_api_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormNameSiloConfig;
