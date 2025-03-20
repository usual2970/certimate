import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

type DeployNodeConfigFormAzureKeyVaultConfigFieldValues = Nullish<{
  keyvaultName: string;
}>;

export type DeployNodeConfigFormAzureKeyVaultConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormAzureKeyVaultConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormAzureKeyVaultConfigFieldValues) => void;
};

const initFormModel = (): DeployNodeConfigFormAzureKeyVaultConfigFieldValues => {
  return {};
};

const DeployNodeConfigFormAzureKeyVaultConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigFormAzureKeyVaultConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    keyvaultName: z
      .string({ message: t("workflow_node.deploy.form.azure_keyvault_name.placeholder") })
      .nonempty(t("workflow_node.deploy.form.azure_keyvault_name.placeholder"))
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
        name="keyvaultName"
        label={t("workflow_node.deploy.form.azure_keyvault_name.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.azure_keyvault_name.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.azure_keyvault_name.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigFormAzureKeyVaultConfig;
