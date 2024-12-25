import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type GoDaddyAccessConfig } from "@/domain/access";

type AccessEditFormGoDaddyConfigModelValues = Partial<GoDaddyAccessConfig>;

export type AccessEditFormGoDaddyConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormGoDaddyConfigModelValues;
  onModelChange?: (model: AccessEditFormGoDaddyConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormGoDaddyConfigModelValues => {
  return {
    apiKey: "",
    apiSecret: "",
  };
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
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onModelChange?.(values as AccessEditFormGoDaddyConfigModelValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
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
