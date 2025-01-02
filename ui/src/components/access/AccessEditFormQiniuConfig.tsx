import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type QiniuAccessConfig } from "@/domain/access";
import { useAntdForm } from "@/hooks";

type AccessEditFormQiniuConfigFieldValues = Partial<QiniuAccessConfig>;

export type AccessEditFormQiniuConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormQiniuConfigFieldValues;
  onValuesChange?: (values: AccessEditFormQiniuConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormQiniuConfigFieldValues => {
  return {
    accessKey: "",
    secretKey: "",
  };
};

const AccessEditFormQiniuConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormQiniuConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKey: z
      .string()
      .min(1, t("access.form.qiniu_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    secretKey: z
      .string()
      .min(1, t("access.form.qiniu_secret_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormQiniuConfigFieldValues);
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
