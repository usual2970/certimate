import { useTranslation } from "react-i18next";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { type HuaweiCloudAccessConfig } from "@/domain/access";

type AccessEditFormHuaweiCloudConfigModelValues = Partial<HuaweiCloudAccessConfig>;

export type AccessEditFormHuaweiCloudConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormHuaweiCloudConfigModelValues;
  onModelChange?: (model: AccessEditFormHuaweiCloudConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormHuaweiCloudConfigModelValues => {
  return {
    accessKeyId: "",
    secretAccessKey: "",
    region: "cn-north-1",
  };
};

const AccessEditFormHuaweiCloudConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormHuaweiCloudConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKeyId: z
      .string()
      .trim()
      .min(1, t("access.form.huaweicloud_access_key_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretAccessKey: z
      .string()
      .trim()
      .min(1, t("access.form.huaweicloud_secret_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    // TODO: 该字段仅用于申请证书，后续迁移到工作流表单中
    region: z
      .string()
      .trim()
      .min(0, t("access.form.huaweicloud_region.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onModelChange?.(values as AccessEditFormHuaweiCloudConfigModelValues);
  };

  return (
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="accessKeyId"
        label={t("access.form.huaweicloud_access_key_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.huaweicloud_access_key_id.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.huaweicloud_access_key_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretAccessKey"
        label={t("access.form.huaweicloud_secret_access_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.huaweicloud_secret_access_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.huaweicloud_secret_access_key.placeholder")} />
      </Form.Item>

      <Form.Item
        name="region"
        label={t("access.form.huaweicloud_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.huaweicloud_region.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.huaweicloud_region.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormHuaweiCloudConfig;
