import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type NameSiloAccessConfig } from "@/domain/access";

type AccessEditFormNameSiloConfigFieldValues = Partial<NameSiloAccessConfig>;

export type AccessEditFormNameSiloConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormNameSiloConfigFieldValues;
  onValuesChange?: (values: AccessEditFormNameSiloConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormNameSiloConfigFieldValues => {
  return {
    apiKey: "",
  };
};

const AccessEditFormNameSiloConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormNameSiloConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    apiKey: z
      .string()
      .min(1, t("access.form.namesilo_api_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormNameSiloConfigFieldValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
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
