import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

type DeployNodeConfigFormBaotaPanelSiteConfigFieldValues = Nullish<{
  siteName: string;
}>;

export type DeployNodeConfigFormBaotaPanelSiteConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormBaotaPanelSiteConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormBaotaPanelSiteConfigFieldValues) => void;
};

const initFormModel = (): DeployNodeConfigFormBaotaPanelSiteConfigFieldValues => {
  return {};
};

const DeployNodeConfigFormBaotaPanelSiteConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigFormBaotaPanelSiteConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    siteName: z
      .string({ message: t("workflow_node.deploy.form.baotapanel_site_name.placeholder") })
      .nonempty(t("workflow_node.deploy.form.baotapanel_site_name.placeholder"))
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
        name="siteName"
        label={t("workflow_node.deploy.form.baotapanel_site_name.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.baotapanel_site_name.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.baotapanel_site_name.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigFormBaotaPanelSiteConfig;
