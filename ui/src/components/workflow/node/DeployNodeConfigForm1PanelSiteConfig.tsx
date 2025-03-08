import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

type DeployNodeConfigForm1PanelSiteConfigFieldValues = Nullish<{
  websiteId: string | number;
}>;

export type DeployNodeConfigForm1PanelSiteConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigForm1PanelSiteConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigForm1PanelSiteConfigFieldValues) => void;
};

const initFormModel = (): DeployNodeConfigForm1PanelSiteConfigFieldValues => {
  return {};
};

const DeployNodeConfigForm1PanelSiteConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigForm1PanelSiteConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    websiteId: z.union([z.string(), z.number()]).refine((v) => {
      return /^\d+$/.test(v + "") && +v > 0;
    }, t("workflow_node.deploy.form.1panel_site_website_id.placeholder")),
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
        name="websiteId"
        label={t("workflow_node.deploy.form.1panel_site_website_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.1panel_site_website_id.tooltip") }}></span>}
      >
        <Input type="number" placeholder={t("workflow_node.deploy.form.1panel_site_website_id.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigForm1PanelSiteConfig;
