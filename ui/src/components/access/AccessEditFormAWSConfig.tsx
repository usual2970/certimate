import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AWSAccessConfig } from "@/domain/access";

type AccessEditFormAWSConfigModelType = Partial<AWSAccessConfig>;

export type AccessEditFormAWSConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormAWSConfigModelType;
  onModelChange?: (model: AccessEditFormAWSConfigModelType) => void;
};

const initModel = () => {
  return {
    accessKeyId: "",
    secretAccessKey: "",
    region: "us-east-1",
    hostedZoneId: "",
  } as AccessEditFormAWSConfigModelType;
};

const AccessEditFormAWSConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormAWSConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    accessKeyId: z
      .string()
      .trim()
      .min(1, t("access.form.aws_access_key_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    secretAccessKey: z
      .string()
      .trim()
      .min(1, t("access.form.aws_secret_access_key.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    // TODO: 该字段仅用于申请证书，后续迁移到工作流表单中
    region: z
      .string()
      .trim()
      .min(0, t("access.form.aws_region.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .nullish(),
    // TODO: 该字段仅用于申请证书，后续迁移到工作流表单中
    hostedZoneId: z
      .string()
      .trim()
      .min(0, t("access.form.aws_hosted_zone_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  const handleFormChange = (_: unknown, fields: AccessEditFormAWSConfigModelType) => {
    onModelChange?.(fields);
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="accessKeyId"
        label={t("access.form.aws_access_key_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aws_access_key_id.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("access.form.aws_access_key_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretAccessKey"
        label={t("access.form.aws_secret_access_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aws_secret_access_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.aws_secret_access_key.placeholder")} />
      </Form.Item>

      <Form.Item
        name="region"
        label={t("access.form.aws_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aws_region.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.aws_region.placeholder")} />
      </Form.Item>

      <Form.Item
        name="hostedZoneId"
        label={t("access.form.aws_hosted_zone_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.aws_hosted_zone_id.tooltip") }}></span>}
      >
        <Input placeholder={t("access.form.aws_hosted_zone_id.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormAWSConfig;
