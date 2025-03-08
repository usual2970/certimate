import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validDomainName } from "@/utils/validators";

type DeployNodeConfigFormAliyunFCConfigFieldValues = Nullish<{
  region: string;
  serviceVersion: string;
  domain: string;
}>;

export type DeployNodeConfigFormAliyunFCConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormAliyunFCConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormAliyunFCConfigFieldValues) => void;
};

const initFormModel = (): DeployNodeConfigFormAliyunFCConfigFieldValues => {
  return {
    serviceVersion: "3.0",
  };
};

const DeployNodeConfigFormAliyunFCConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: DeployNodeConfigFormAliyunFCConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    serviceVersion: z.union([z.literal("2.0"), z.literal("3.0")], {
      message: t("workflow_node.deploy.form.aliyun_fc_service_version.placeholder"),
    }),
    region: z
      .string({ message: t("workflow_node.deploy.form.aliyun_fc_region.placeholder") })
      .nonempty(t("workflow_node.deploy.form.aliyun_fc_region.placeholder"))
      .trim(),
    domain: z
      .string({ message: t("workflow_node.deploy.form.aliyun_fc_domain.placeholder") })
      .refine((v) => validDomainName(v), t("common.errmsg.domain_invalid")),
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
      <Form.Item name="serviceVersion" label={t("workflow_node.deploy.form.aliyun_fc_service_version.label")} rules={[formRule]}>
        <Select placeholder={t("workflow_node.deploy.form.aliyun_fc_service_version.placeholder")}>
          <Select.Option key="2.0" value="2.0">
            2.0
          </Select.Option>
          <Select.Option key="3.0" value="3.0">
            3.0
          </Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        name="region"
        label={t("workflow_node.deploy.form.aliyun_fc_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_fc_region.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_fc_region.placeholder")} />
      </Form.Item>

      <Form.Item
        name="domain"
        label={t("workflow_node.deploy.form.aliyun_fc_domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_fc_domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_fc_domain.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigFormAliyunFCConfig;
