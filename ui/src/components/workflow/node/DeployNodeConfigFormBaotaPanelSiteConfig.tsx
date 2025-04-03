import { memo } from "react";
import { useTranslation } from "react-i18next";
import { FormOutlined as FormOutlinedIcon } from "@ant-design/icons";
import { Button, Form, type FormInstance, Input, Select, Space } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import ModalForm from "@/components/ModalForm";
import MultipleInput from "@/components/MultipleInput";
import Show from "@/components/Show";
import { useAntdForm } from "@/hooks";

type DeployNodeConfigFormBaotaPanelSiteConfigFieldValues = Nullish<{
  siteType: string;
  siteName?: string;
  siteNames?: string;
}>;

export type DeployNodeConfigFormBaotaPanelSiteConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormBaotaPanelSiteConfigFieldValues;
  onValuesChange?: (values: DeployNodeConfigFormBaotaPanelSiteConfigFieldValues) => void;
};

const SITE_TYPE_PHP = "php";
const SITE_TYPE_OTHER = "other";

const MULTIPLE_INPUT_DELIMITER = ";";

const initFormModel = (): DeployNodeConfigFormBaotaPanelSiteConfigFieldValues => {
  return {
    siteType: SITE_TYPE_OTHER,
    siteName: "",
    siteNames: "",
  };
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
    siteType: z.union([z.literal(SITE_TYPE_PHP), z.literal(SITE_TYPE_OTHER)], {
      message: t("workflow_node.deploy.form.baotapanel_site_type.placeholder"),
    }),
    siteName: z
      .string()
      .nullish()
      .refine((v) => {
        if (fieldSiteType !== SITE_TYPE_PHP) return true;
        return !!v?.trim();
      }, t("workflow_node.deploy.form.baotapanel_site_name.placeholder")),
    siteNames: z
      .string()
      .nullish()
      .refine((v) => {
        if (fieldSiteType !== SITE_TYPE_OTHER) return true;
        if (!v) return false;
        return String(v)
          .split(MULTIPLE_INPUT_DELIMITER)
          .every((e) => !!e.trim());
      }, t("workflow_node.deploy.form.baotapanel_site_names.placeholder")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const fieldSiteType = Form.useWatch<string>("siteType", formInst);
  const fieldSiteNames = Form.useWatch<string>("siteNames", formInst);

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
      <Form.Item name="siteType" label={t("workflow_node.deploy.form.baotapanel_site_type.label")} rules={[formRule]}>
        <Select placeholder={t("workflow_node.deploy.form.baotapanel_site_type.placeholder")}>
          <Select.Option key={SITE_TYPE_PHP} value={SITE_TYPE_PHP}>
            {t("workflow_node.deploy.form.baotapanel_site_type.option.php.label")}
          </Select.Option>
          <Select.Option key={SITE_TYPE_OTHER} value={SITE_TYPE_OTHER}>
            {t("workflow_node.deploy.form.baotapanel_site_type.option.other.label")}
          </Select.Option>
        </Select>
      </Form.Item>

      <Show when={fieldSiteType === SITE_TYPE_PHP}>
        <Form.Item
          name="siteName"
          label={t("workflow_node.deploy.form.baotapanel_site_name.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.baotapanel_site_name.tooltip") }}></span>}
        >
          <Input placeholder={t("workflow_node.deploy.form.baotapanel_site_name.placeholder")} />
        </Form.Item>
      </Show>

      <Show when={fieldSiteType === SITE_TYPE_OTHER}>
        <Form.Item
          label={t("workflow_node.deploy.form.baotapanel_site_names.label")}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.baotapanel_site_names.tooltip") }}></span>}
        >
          <Space.Compact style={{ width: "100%" }}>
            <Form.Item name="siteNames" noStyle rules={[formRule]}>
              <Input
                allowClear
                disabled={disabled}
                value={fieldSiteNames}
                placeholder={t("workflow_node.deploy.form.baotapanel_site_names.placeholder")}
                onChange={(e) => {
                  formInst.setFieldValue("siteNames", e.target.value);
                }}
                onClear={() => {
                  formInst.setFieldValue("siteNames", "");
                }}
              />
            </Form.Item>
            <SiteNamesModalInput
              value={fieldSiteNames}
              trigger={
                <Button disabled={disabled}>
                  <FormOutlinedIcon />
                </Button>
              }
              onChange={(value) => {
                formInst.setFieldValue("siteNames", value);
              }}
            />
          </Space.Compact>
        </Form.Item>
      </Show>
    </Form>
  );
};

const SiteNamesModalInput = memo(({ value, trigger, onChange }: { value?: string; trigger?: React.ReactNode; onChange?: (value: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    siteNames: z.array(z.string()).refine((v) => {
      return v.every((e) => !!e?.trim());
    }, t("workflow_node.deploy.form.baotapanel_site_names.errmsg.invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    name: "workflowNodeDeployConfigFormBaotaPanelSiteNamesModalInput",
    initialValues: { siteNames: value?.split(MULTIPLE_INPUT_DELIMITER) },
    onSubmit: (values) => {
      onChange?.(
        values.siteNames
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
      title={t("workflow_node.deploy.form.baotapanel_site_names.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
    >
      <Form.Item name="siteNames" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.deploy.form.baotapanel_site_names.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
});

export default DeployNodeConfigFormBaotaPanelSiteConfig;
