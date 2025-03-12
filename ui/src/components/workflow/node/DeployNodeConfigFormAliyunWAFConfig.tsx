import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validDomainName } from "@/utils/validators";

type DeployNodeConfigFormAliyunWAFConfigFieldValues = Nullish<{
  region: string;
  instanceId: string;
  domain?: string;
}>;

export type DeployNodeConfigFormAliyunWAFConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormAliyunWAFConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormAliyunWAFConfigFieldValues) => void;
};

const initFormModel = (): DeployNodeConfigFormAliyunWAFConfigFieldValues => {
  return {};
};

const DeployNodeConfigFormAliyunWAFConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigFormAliyunWAFConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    region: z
      .string({ message: t("workflow_node.deploy.form.aliyun_waf_region.placeholder") })
      .nonempty(t("workflow_node.deploy.form.aliyun_waf_region.placeholder"))
      .trim(),
    instanceId: z
      .string({ message: t("workflow_node.deploy.form.aliyun_waf_instance_id.placeholder") })
      .nonempty(t("workflow_node.deploy.form.aliyun_waf_instance_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    domain: z
      .string()
      .nullish()
      .refine((v) => {
        return !v || validDomainName(v!, { allowWildcard: true });
      }, t("common.errmsg.domain_invalid")),
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
        name="region"
        label={t("workflow_node.deploy.form.aliyun_waf_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_waf_region.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_waf_region.placeholder")} />
      </Form.Item>

      <Form.Item
        name="instanceId"
        label={t("workflow_node.deploy.form.aliyun_waf_instance_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_waf_instance_id.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_waf_instance_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="domain"
        label={t("workflow_node.deploy.form.aliyun_waf_domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_waf_domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_waf_domain.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default DeployNodeConfigFormAliyunWAFConfig;
