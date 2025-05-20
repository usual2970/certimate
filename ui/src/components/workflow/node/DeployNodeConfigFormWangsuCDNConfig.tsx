import { memo } from "react";
import { useTranslation } from "react-i18next";
import { FormOutlined as FormOutlinedIcon } from "@ant-design/icons";
import { Button, Form, type FormInstance, Input, Space } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import ModalForm from "@/components/ModalForm";
import MultipleInput from "@/components/MultipleInput";
import { useAntdForm } from "@/hooks";
import { validDomainName } from "@/utils/validators";

type DeployNodeConfigFormWangsuCDNConfigFieldValues = Nullish<{
  domains: string;
}>;

export type DeployNodeConfigFormWangsuCDNConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormWangsuCDNConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormWangsuCDNConfigFieldValues) => void;
};

const MULTIPLE_INPUT_DELIMITER = ";";

const initFormModel = (): DeployNodeConfigFormWangsuCDNConfigFieldValues => {
  return {
    domains: "",
  };
};

const DeployNodeConfigFormWangsuCDNConfig = ({
  form: formInst,
  formName,
  disabled,
  initialValues,
  onValuesChange,
}: DeployNodeConfigFormWangsuCDNConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    domains: z
      .string()
      .nullish()
      .refine((v) => {
        if (!v) return false;
        return String(v)
          .split(MULTIPLE_INPUT_DELIMITER)
          .every((e) => validDomainName(e));
      }, t("workflow_node.deploy.form.wangsu_cdn_domains.placeholder")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const fieldDomains = Form.useWatch<string>("domains", formInst);

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
        label={t("workflow_node.deploy.form.wangsu_cdn_domains.label")}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.wangsu_cdn_domains.tooltip") }}></span>}
      >
        <Space.Compact style={{ width: "100%" }}>
          <Form.Item name="domains" noStyle rules={[formRule]}>
            <Input
              allowClear
              disabled={disabled}
              value={fieldDomains}
              placeholder={t("workflow_node.deploy.form.wangsu_cdn_domains.placeholder")}
              onChange={(e) => {
                formInst.setFieldValue("domains", e.target.value);
              }}
              onClear={() => {
                formInst.setFieldValue("domains", "");
              }}
            />
          </Form.Item>
          <SiteNamesModalInput
            value={fieldDomains}
            trigger={
              <Button disabled={disabled}>
                <FormOutlinedIcon />
              </Button>
            }
            onChange={(value) => {
              formInst.setFieldValue("domains", value);
            }}
          />
        </Space.Compact>
      </Form.Item>
    </Form>
  );
};

const SiteNamesModalInput = memo(({ value, trigger, onChange }: { value?: string; trigger?: React.ReactNode; onChange?: (value: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    domains: z.array(z.string()).refine((v) => {
      return v.every((e) => validDomainName(e));
    }, t("workflow_node.deploy.form.wangsu_cdn_domains.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    name: "workflowNodeDeployConfigFormWangsuCDNNamesModalInput",
    initialValues: { domains: value?.split(MULTIPLE_INPUT_DELIMITER) },
    onSubmit: (values) => {
      onChange?.(
        values.domains
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
      title={t("workflow_node.deploy.form.wangsu_cdn_domains.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
    >
      <Form.Item name="domains" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.deploy.form.wangsu_cdn_domains.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
});

export default DeployNodeConfigFormWangsuCDNConfig;
