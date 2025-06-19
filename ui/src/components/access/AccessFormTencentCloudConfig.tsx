import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForTencentCloud } from "@/domain/access";

type AccessFormTencentCloudConfigFieldValues = Nullish<AccessConfigForTencentCloud>;

export type AccessFormTencentCloudConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormTencentCloudConfigFieldValues;
  onValuesChange?: (values: AccessFormTencentCloudConfigFieldValues) => void;
};

const initFormModel = (): AccessFormTencentCloudConfigFieldValues => {
  return {
    secretId: "",
    secretKey: "",
  };
};

const AccessFormTencentCloudConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormTencentCloudConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    secretId: z
      .string()
      .min(1, t("access.form.tencentcloud_secret_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretKey: z
      .string()
      .min(1, t("access.form.tencentcloud_secret_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values);
  };

  return (
    <Form
      form={formInst}
      disabled={disabled}
      initialValues={initialValues ?? initFormModel()}
      layout="vertical"
      name={formName}
      onValuesChange={handleFormChange}
    >
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

export default AccessFormTencentCloudConfig;
