import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type QiniuAccessConfig } from "@/domain/access";

type AccessEditFormQiniuConfigModelValues = Partial<QiniuAccessConfig>;

export type AccessEditFormQiniuConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormQiniuConfigModelValues;
  onModelChange?: (model: AccessEditFormQiniuConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormQiniuConfigModelValues => {
  return {
    accessKey: "",
    secretKey: "",
  };
};

const AccessEditFormQiniuConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormQiniuConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKey: z
      .string()
      .trim()
      .min(1, t("access.form.qiniu_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretKey: z
      .string()
      .trim()
      .min(1, t("access.form.qiniu_secret_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onModelChange?.(values as AccessEditFormQiniuConfigModelValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="accessKey"
        label={t("access.form.qiniu_access_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.qiniu_access_key.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.qiniu_access_key.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretKey"
        label={t("access.form.qiniu_secret_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.qiniu_secret_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.qiniu_secret_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormQiniuConfig;
