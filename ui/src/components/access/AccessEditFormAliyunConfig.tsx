import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AliyunAccessConfig } from "@/domain/access";
import { useAntdForm } from "@/hooks";

type AccessEditFormAliyunConfigFieldValues = Partial<AliyunAccessConfig>;

export type AccessEditFormAliyunConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormAliyunConfigFieldValues;
  onValuesChange?: (values: AccessEditFormAliyunConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormAliyunConfigFieldValues => {
  return {
    accessKeyId: "",
    accessKeySecret: "",
  };
};

const AccessEditFormAliyunConfig = ({ form, formName, disabled, initialValues, onValuesChange: onValuesChange }: AccessEditFormAliyunConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKeyId: z
      .string()
      .min(1, t("access.form.aliyun_access_key_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    accessKeySecret: z
      .string()
      .min(1, t("access.form.aliyun_access_key_secret.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormAliyunConfigFieldValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="accessKeyId"
        label={t("access.form.aliyun_access_key_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aliyun_access_key_id.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.aliyun_access_key_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="accessKeySecret"
        label={t("access.form.aliyun_access_key_secret.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aliyun_access_key_secret.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.aliyun_access_key_secret.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormAliyunConfig;
