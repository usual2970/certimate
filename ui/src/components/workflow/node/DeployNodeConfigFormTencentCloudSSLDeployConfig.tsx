import { memo } from "react";
import { useTranslation } from "react-i18next";
import { FormOutlined as FormOutlinedIcon } from "@ant-design/icons";
import { Alert, AutoComplete, Button, Form, type FormInstance, Input, Space } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import ModalForm from "@/components/ModalForm";
import MultipleInput from "@/components/MultipleInput";
import { useAntdForm } from "@/hooks";

type DeployNodeConfigFormTencentCloudSSLDeployConfigFieldValues = Nullish<{
  region: string;
  resourceType: string;
  resourceIds: string;
}>;

export type DeployNodeConfigFormTencentCloudSSLDeployConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormTencentCloudSSLDeployConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormTencentCloudSSLDeployConfigFieldValues) => void;
};

const MULTIPLE_INPUT_DELIMITER = ";";

const initFormModel = (): DeployNodeConfigFormTencentCloudSSLDeployConfigFieldValues => {
  return {};
};

const DeployNodeConfigFormTencentCloudSSLDeployConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigFormTencentCloudSSLDeployConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    region: z
      .string({ message: t("workflow_node.deploy.form.tencentcloud_ssl_deploy_region.placeholder") })
      .nonempty(t("workflow_node.deploy.form.tencentcloud_ssl_deploy_region.placeholder"))
      .trim(),
    resourceType: z
      .string({ message: t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_type.placeholder") })
      .nonempty(t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_type.placeholder"))
      .trim(),
    resourceIds: z.string({ message: t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.placeholder") }).refine((v) => {
      if (!v) return false;
      return String(v)
        .split(MULTIPLE_INPUT_DELIMITER)
        .every((e) => /^[A-Za-z0-9*._-]+$/.test(e));
    }, t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const fieldResourceIds = Form.useWatch<string>("resourceIds", formInst);

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
        label={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_ssl_deploy_region.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_region.placeholder")} />
      </Form.Item>

      <Form.Item
        name="resourceType"
        label={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_type.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_type.tooltip") }}></span>}
      >
        <AutoComplete
          options={["apigateway", "cdn", "clb", "cos", "ddos", "lighthouse", "live", "tcb", "teo", "tke", "tse", "vod", "waf"].map((value) => ({ value }))}
          placeholder={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_type.placeholder")}
          filterOption={(inputValue, option) => option!.value.toLowerCase().includes(inputValue.toLowerCase())}
        />
      </Form.Item>

      <Form.Item
        label={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.label")}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.tooltip") }}></span>}
      >
        <Space.Compact style={{ width: "100%" }}>
          <Form.Item name="resourceIds" noStyle rules={[formRule]}>
            <Input
              allowClear
              disabled={disabled}
              value={fieldResourceIds}
              placeholder={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.placeholder")}
              onChange={(e) => {
                formInst.setFieldValue("resourceIds", e.target.value);
              }}
              onClear={() => {
                formInst.setFieldValue("resourceIds", "");
              }}
            />
          </Form.Item>
          <ResourceIdsModalInput
            value={fieldResourceIds}
            trigger={
              <Button disabled={disabled}>
                <FormOutlinedIcon />
              </Button>
            }
            onChange={(value) => {
              formInst.setFieldValue("resourceIds", value);
            }}
          />
        </Space.Compact>
      </Form.Item>

      <Form.Item>
        <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_ssl_deploy.guide") }}></span>} />
      </Form.Item>
    </Form>
  );
};

const ResourceIdsModalInput = memo(({ value, trigger, onChange }: { value?: string; trigger?: React.ReactNode; onChange?: (value: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    resourceIds: z.array(z.string()).refine((v) => {
      return v.every((e) => !e?.trim() || /^[A-Za-z0-9*._-]+$/.test(e));
    }, t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    name: "workflowNodeDeployConfigFormTencentCloudSSLDeployResourceIdsModalInput",
    initialValues: { resourceIds: value?.split(MULTIPLE_INPUT_DELIMITER) },
    onSubmit: (values) => {
      onChange?.(
        values.resourceIds
          .map((e) => e.trim())
          .filter((e) => !!e)
          .join(MULTIPLE_INPUT_DELIMITER)
      );
    },
  });

  return (
    <ModalForm
      {...formProps}
      layout="vertical"
      form={formInst}
      modalProps={{ destroyOnClose: true }}
      title={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
    >
      <Form.Item name="resourceIds" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.deploy.form.tencentcloud_ssl_deploy_resource_ids.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
});

export default DeployNodeConfigFormTencentCloudSSLDeployConfig;
