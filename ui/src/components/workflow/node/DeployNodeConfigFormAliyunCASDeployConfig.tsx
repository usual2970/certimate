import { memo } from "react";
import { useTranslation } from "react-i18next";
import { FormOutlined as FormOutlinedIcon } from "@ant-design/icons";
import { Alert, Button, Form, type FormInstance, Input, Space } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import ModalForm from "@/components/ModalForm";
import MultipleInput from "@/components/MultipleInput";
import { useAntdForm } from "@/hooks";

type DeployNodeConfigFormAliyunCASDeployConfigFieldValues = Nullish<{
  region: string;
  resourceIds: string;
  contactIds?: string;
}>;

export type DeployNodeConfigFormAliyunCASDeployConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormAliyunCASDeployConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormAliyunCASDeployConfigFieldValues) => void;
};

const MULTIPLE_INPUT_DELIMITER = ";";

const initFormModel = (): DeployNodeConfigFormAliyunCASDeployConfigFieldValues => {
  return {};
};

const DeployNodeConfigFormAliyunCASDeployConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigFormAliyunCASDeployConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    region: z
      .string({ message: t("workflow_node.deploy.form.aliyun_cas_deploy_region.placeholder") })
      .nonempty(t("workflow_node.deploy.form.aliyun_cas_deploy_region.placeholder"))
      .trim(),
    resourceIds: z.string({ message: t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.placeholder") }).refine((v) => {
      if (!v) return false;
      return String(v)
        .split(MULTIPLE_INPUT_DELIMITER)
        .every((e) => /^[1-9]\d*$/.test(e));
    }, t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.errmsg.invalid")),
    contactIds: z
      .string({ message: t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.placeholder") })
      .nullish()
      .refine((v) => {
        if (!v) return true;
        return String(v)
          .split(MULTIPLE_INPUT_DELIMITER)
          .every((e) => /^[1-9]\d*$/.test(e));
      }, t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const fieldResourceIds = Form.useWatch<string>("resourceIds", formInst);
  const fieldContactIds = Form.useWatch<string>("contactIds", formInst);

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
        label={t("workflow_node.deploy.form.aliyun_cas_deploy_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_cas_deploy_region.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_cas_deploy_region.placeholder")} />
      </Form.Item>

      <Form.Item
        label={t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.label")}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.tooltip") }}></span>}
      >
        <Space.Compact style={{ width: "100%" }}>
          <Form.Item name="resourceIds" noStyle rules={[formRule]}>
            <Input
              allowClear
              disabled={disabled}
              value={fieldResourceIds}
              placeholder={t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.placeholder")}
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

      <Form.Item
        label={t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.label")}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.tooltip") }}></span>}
      >
        <Space.Compact style={{ width: "100%" }}>
          <Form.Item name="contactIds" noStyle rules={[formRule]}>
            <Input
              allowClear
              disabled={disabled}
              value={fieldContactIds}
              placeholder={t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.placeholder")}
              onChange={(e) => {
                formInst.setFieldValue("contactIds", e.target.value);
              }}
              onClear={() => {
                formInst.setFieldValue("contactIds", "");
              }}
            />
          </Form.Item>
          <ContactIdsModalInput
            value={fieldContactIds}
            trigger={
              <Button disabled={disabled}>
                <FormOutlinedIcon />
              </Button>
            }
            onChange={(value) => {
              formInst.setFieldValue("contactIds", value);
            }}
          />
        </Space.Compact>
      </Form.Item>

      <Form.Item>
        <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_cas_deploy.guide") }}></span>} />
      </Form.Item>
    </Form>
  );
};

const ResourceIdsModalInput = memo(({ value, trigger, onChange }: { value?: string; trigger?: React.ReactNode; onChange?: (value: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    resourceIds: z.array(z.string()).refine((v) => {
      return v.every((e) => !e?.trim() || /^[1-9]\d*$/.test(e));
    }, t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    name: "workflowNodeDeployConfigFormAliyunCASResourceIdsModalInput",
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
      modalProps={{ destroyOnHidden: true }}
      title={t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
    >
      <Form.Item name="resourceIds" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.deploy.form.aliyun_cas_deploy_resource_ids.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
});

const ContactIdsModalInput = memo(({ value, trigger, onChange }: { value?: string; trigger?: React.ReactNode; onChange?: (value: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    contactIds: z.array(z.string()).refine((v) => {
      return v.every((e) => !e?.trim() || /^[1-9]\d*$/.test(e));
    }, t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    name: "workflowNodeDeployConfigFormAliyunCASDeploymentJobContactIdsModalInput",
    initialValues: { contactIds: value?.split(MULTIPLE_INPUT_DELIMITER) },
    onSubmit: (values) => {
      onChange?.(
        values.contactIds
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
      modalProps={{ destroyOnHidden: true }}
      title={t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
    >
      <Form.Item name="contactIds" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.deploy.form.aliyun_cas_deploy_contact_ids.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
});

export default DeployNodeConfigFormAliyunCASDeployConfig;
