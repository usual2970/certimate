import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForDogeCloud } from "@/domain/access";

type AccessFormDogeCloudConfigFieldValues = Nullish<AccessConfigForDogeCloud>;

export type AccessFormDogeCloudConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormDogeCloudConfigFieldValues;
  onValuesChange?: (values: AccessFormDogeCloudConfigFieldValues) => void;
};

const initFormModel = (): AccessFormDogeCloudConfigFieldValues => {
  return {
    accessKey: "",
    secretKey: "",
  };
};

const AccessFormDogeCloudConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormDogeCloudConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKey: z
      .string()
      .min(1, t("access.form.dogecloud_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretKey: z
      .string()
      .min(1, t("access.form.dogecloud_secret_key.placeholder"))
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
        name="accessKey"
        label={t("access.form.dogecloud_access_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.dogecloud_access_key.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.dogecloud_access_key.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretKey"
        label={t("access.form.dogecloud_secret_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.dogecloud_secret_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.dogecloud_secret_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessFormDogeCloudConfig;
