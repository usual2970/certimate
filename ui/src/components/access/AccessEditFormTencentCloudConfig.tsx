import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type TencentCloudAccessConfig } from "@/domain/access";
import { useAntdForm } from "@/hooks";

type AccessEditFormTencentCloudConfigFieldValues = Partial<TencentCloudAccessConfig>;

export type AccessEditFormTencentCloudConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormTencentCloudConfigFieldValues;
  onValuesChange?: (values: AccessEditFormTencentCloudConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormTencentCloudConfigFieldValues => {
  return {
    secretId: "",
    secretKey: "",
  };
};

const AccessEditFormTencentCloudConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormTencentCloudConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    secretId: z
      .string()
      .min(1, t("access.form.tencentcloud_secret_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    secretKey: z
      .string()
      .min(1, t("access.form.tencentcloud_secret_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormTencentCloudConfigFieldValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="secretId"
        label={t("access.form.tencentcloud_secret_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.tencentcloud_secret_id.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.tencentcloud_secret_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretKey"
        label={t("access.form.tencentcloud_secret_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.tencentcloud_secret_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.tencentcloud_secret_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormTencentCloudConfig;
