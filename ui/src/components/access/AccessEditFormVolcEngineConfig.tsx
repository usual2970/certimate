import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type VolcEngineAccessConfig } from "@/domain/access";

type AccessEditFormVolcEngineConfigModelValues = Partial<VolcEngineAccessConfig>;

export type AccessEditFormVolcEngineConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormVolcEngineConfigModelValues;
  onModelChange?: (model: AccessEditFormVolcEngineConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormVolcEngineConfigModelValues => {
  return {
    accessKeyId: "",
    secretAccessKey: "",
  };
};

const AccessEditFormVolcEngineConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormVolcEngineConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKeyId: z
      .string()
      .trim()
      .min(1, t("access.form.volcengine_access_key_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretAccessKey: z
      .string()
      .trim()
      .min(1, t("access.form.volcengine_secret_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onModelChange?.(values as AccessEditFormVolcEngineConfigModelValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="accessKeyId"
        label={t("access.form.volcengine_access_key_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.volcengine_access_key_id.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.volcengine_access_key_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretAccessKey"
        label={t("access.form.volcengine_secret_access_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.volcengine_secret_access_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.volcengine_secret_access_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormVolcEngineConfig;
