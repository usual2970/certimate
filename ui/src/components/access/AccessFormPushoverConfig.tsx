import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForPushover } from "@/domain/access";

type AccessFormPushoverConfigFieldValues = Nullish<AccessConfigForPushover>;

export type AccessFormPushoverConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormPushoverConfigFieldValues;
  onValuesChange?: (values: AccessFormPushoverConfigFieldValues) => void;
};

const initFormModel = (): AccessFormPushoverConfigFieldValues => {
  return {
    token: "",
    user: "",
  };
};

const AccessFormPushoverConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormPushoverConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    token: z
      .string({ message: t("access.form.pushover_token.placeholder") })
      .nonempty(t("access.form.pushover_token.placeholder")),
    user: z
      .string({ message: t("access.form.pushover_user.placeholder") })
      .nonempty(t("access.form.pushover_user.placeholder")),
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
        name="token"
        label={t("access.form.pushover_token.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.pushover_token.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.pushover_token.placeholder")} />
      </Form.Item>
      <Form.Item
        name="user"
        label={t("access.form.pushover_user.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.pushover_user.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.pushover_user.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessFormPushoverConfig;
