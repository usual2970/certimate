import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type NameDotComAccessConfig } from "@/domain/access";
import { useAntdForm } from "@/hooks";

type AccessEditFormNameDotComConfigFieldValues = Partial<NameDotComAccessConfig>;

export type AccessEditFormNameDotComConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormNameDotComConfigFieldValues;
  onValuesChange?: (values: AccessEditFormNameDotComConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormNameDotComConfigFieldValues => {
  return {
    username: "",
    apiToken: "",
  };
};

const AccessEditFormNameDotComConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormNameDotComConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    username: z
      .string()
      .min(1, t("access.form.namedotcom_username.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    apiToken: z
      .string()
      .min(1, t("access.form.namedotcom_api_token.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormNameDotComConfigFieldValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
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
