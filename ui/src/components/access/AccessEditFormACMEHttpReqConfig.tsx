import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForACMEHttpReq } from "@/domain/access";
import { useAntdForm } from "@/hooks";

type AccessEditFormACMEHttpReqConfigFieldValues = Partial<AccessConfigForACMEHttpReq>;

export type AccessEditFormACMEHttpReqConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormACMEHttpReqConfigFieldValues;
  onValuesChange?: (values: AccessEditFormACMEHttpReqConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormACMEHttpReqConfigFieldValues => {
  return {
    endpoint: "https://example.com/api/",
    mode: "",
  };
};

const AccessEditFormACMEHttpReqConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormACMEHttpReqConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    endpoint: z.string().url(t("common.errmsg.url_invalid")),
    mode: z.string().nullish(),
    username: z
      .string()
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish(),
    password: z
      .string()
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormACMEHttpReqConfigFieldValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="endpoint"
        label={t("access.form.acmehttpreq_endpoint.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmehttpreq_endpoint.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.acmehttpreq_endpoint.placeholder")} />
      </Form.Item>

      <Form.Item
        name="mode"
        label={t("access.form.acmehttpreq_mode.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmehttpreq_mode.tooltip") }}></span>}
      >
        <Select
          options={[
            { value: "", label: "(default)" },
            { value: "RAW", label: "RAW" },
          ]}
          placeholder={t("access.form.acmehttpreq_mode.placeholder")}
        />
      </Form.Item>

      <Form.Item
        name="username"
        label={t("access.form.acmehttpreq_username.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmehttpreq_username.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.acmehttpreq_username.placeholder")} />
      </Form.Item>

      <Form.Item
        name="password"
        label={t("access.form.acmehttpreq_password.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.acmehttpreq_password.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.acmehttpreq_password.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormACMEHttpReqConfig;
