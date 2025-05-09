import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input, Switch } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigFor1Panel } from "@/domain/access";

type AccessForm1PanelConfigFieldValues = Nullish<AccessConfigFor1Panel>;

export type AccessForm1PanelConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessForm1PanelConfigFieldValues;
  onValuesChange?: (values: AccessForm1PanelConfigFieldValues) => void;
};

const initFormModel = (): AccessForm1PanelConfigFieldValues => {
  return {
    apiUrl: "http://<your-host-addr>:20410/",
    apiKey: "",
  };
};

const AccessForm1PanelConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessForm1PanelConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    apiUrl: z.string().url(t("common.errmsg.url_invalid")),
    apiKey: z
      .string()
      .min(1, t("access.form.1panel_api_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    allowInsecureConnections: z.boolean().nullish(),
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
      <Form.Item name="apiUrl" label={t("access.form.1panel_api_url.label")} rules={[formRule]}>
        <Input placeholder={t("access.form.1panel_api_url.placeholder")} />
      </Form.Item>

      <Form.Item
        name="apiKey"
        label={t("access.form.1panel_api_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.1panel_api_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.1panel_api_key.placeholder")} />
      </Form.Item>

      <Form.Item name="allowInsecureConnections" label={t("access.form.1panel_allow_insecure_conns.label")} rules={[formRule]}>
        <Switch
          checkedChildren={t("access.form.1panel_allow_insecure_conns.switch.on")}
          unCheckedChildren={t("access.form.1panel_allow_insecure_conns.switch.off")}
        />
      </Form.Item>
    </Form>
  );
};

export default AccessForm1PanelConfig;
