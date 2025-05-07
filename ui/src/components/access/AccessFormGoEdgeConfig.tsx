import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForGoEdge } from "@/domain/access";

type AccessFormGoEdgeConfigFieldValues = Nullish<AccessConfigForGoEdge>;

export type AccessFormGoEdgeConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormGoEdgeConfigFieldValues;
  onValuesChange?: (values: AccessFormGoEdgeConfigFieldValues) => void;
};

const initFormModel = (): AccessFormGoEdgeConfigFieldValues => {
  return {
    apiUrl: "http://<your-host-addr>:7788/",
    accessKeyId: "",
    accessKey: "",
  };
};

const AccessFormGoEdgeConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormGoEdgeConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    apiUrl: z.string().url(t("common.errmsg.url_invalid")),
    accessKeyId: z
      .string()
      .min(1, t("access.form.goedge_access_key_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    accessKey: z
      .string()
      .min(1, t("access.form.goedge_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
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
        name="apiUrl"
        label={t("access.form.goedge_api_url.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.goedge_api_url.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.goedge_api_url.placeholder")} />
      </Form.Item>

      <Form.Item
        name="accessKeyId"
        label={t("access.form.goedge_access_key_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.goedge_access_key_id.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.goedge_access_key_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="accessKey"
        label={t("access.form.goedge_access_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.goedge_access_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.goedge_access_key.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessFormGoEdgeConfig;
