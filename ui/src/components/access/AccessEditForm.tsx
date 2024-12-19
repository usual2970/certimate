import { forwardRef, useEffect, useImperativeHandle, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import {
  ACCESS_PROVIDER_TYPES,
  type AccessModel,
  type ACMEHttpReqAccessConfig,
  type AliyunAccessConfig,
  type AWSAccessConfig,
  type BaiduCloudAccessConfig,
  type BytePlusAccessConfig,
  type CloudflareAccessConfig,
  type DogeCloudAccessConfig,
  type GoDaddyAccessConfig,
  type HuaweiCloudAccessConfig,
  type KubernetesAccessConfig,
  type LocalAccessConfig,
  type NameSiloAccessConfig,
  type PowerDNSAccessConfig,
  type QiniuAccessConfig,
  type SSHAccessConfig,
  type TencentCloudAccessConfig,
  type VolcEngineAccessConfig,
  type WebhookAccessConfig,
} from "@/domain/access";
import AccessTypeSelect from "./AccessTypeSelect";
import AccessEditFormACMEHttpReqConfig from "./AccessEditFormACMEHttpReqConfig";
import AccessEditFormAliyunConfig from "./AccessEditFormAliyunConfig";
import AccessEditFormAWSConfig from "./AccessEditFormAWSConfig";
import AccessEditFormBaiduCloudConfig from "./AccessEditFormBaiduCloudConfig";
import AccessEditFormBytePlusConfig from "./AccessEditFormBytePlusConfig";
import AccessEditFormCloudflareConfig from "./AccessEditFormCloudflareConfig";
import AccessEditFormDogeCloudConfig from "./AccessEditFormDogeCloudConfig";
import AccessEditFormGoDaddyConfig from "./AccessEditFormGoDaddyConfig";
import AccessEditFormHuaweiCloudConfig from "./AccessEditFormHuaweiCloudConfig";
import AccessEditFormKubernetesConfig from "./AccessEditFormKubernetesConfig";
import AccessEditFormLocalConfig from "./AccessEditFormLocalConfig";
import AccessEditFormNameSiloConfig from "./AccessEditFormNameSiloConfig";
import AccessEditFormPowerDNSConfig from "./AccessEditFormPowerDNSConfig";
import AccessEditFormQiniuConfig from "./AccessEditFormQiniuConfig";
import AccessEditFormSSHConfig from "./AccessEditFormSSHConfig";
import AccessEditFormTencentCloudConfig from "./AccessEditFormTencentCloudConfig";
import AccessEditFormVolcEngineConfig from "./AccessEditFormVolcEngineConfig";
import AccessEditFormWebhookConfig from "./AccessEditFormWebhookConfig";

type AccessEditFormModelType = Partial<MaybeModelRecord<AccessModel>>;

export type AccessEditFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  loading?: boolean;
  mode: "add" | "edit";
  model?: AccessEditFormModelType;
  onModelChange?: (model: AccessEditFormModelType) => void;
};

export type AccessEditFormInstance = {
  getFieldsValue: () => AccessEditFormModelType;
  resetFields: () => void;
  validateFields: () => Promise<AccessEditFormModelType>;
};

const AccessEditForm = forwardRef<AccessEditFormInstance, AccessEditFormProps>(({ className, style, disabled, loading, mode, model, onModelChange }, ref) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    name: z
      .string()
      .trim()
      .min(1, t("access.form.name.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    configType: z.nativeEnum(ACCESS_PROVIDER_TYPES, { message: t("access.form.type.placeholder") }),
    config: z.any(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const [form] = Form.useForm<z.infer<typeof formSchema>>();

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model as Partial<z.infer<typeof formSchema>>);
  useDeepCompareEffect(() => {
    setInitialValues(model as Partial<z.infer<typeof formSchema>>);
  }, [model]);

  const [configType, setConfigType] = useState(model?.configType);
  useEffect(() => {
    setConfigType(model?.configType);
  }, [model?.configType]);

  const [configFormInst] = Form.useForm();
  const configFormComponent = useMemo(() => {
    /*
      注意：如果追加新的子组件，请保持以 ASCII 排序。
      NOTICE: If you add new child component, please keep ASCII order.
     */
    switch (configType) {
      case ACCESS_PROVIDER_TYPES.ACMEHTTPREQ:
        return <AccessEditFormACMEHttpReqConfig form={configFormInst} model={model?.config as ACMEHttpReqAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.ALIYUN:
        return <AccessEditFormAliyunConfig form={configFormInst} model={model?.config as AliyunAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.AWS:
        return <AccessEditFormAWSConfig form={configFormInst} model={model?.config as AWSAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.BAIDUCLOUD:
        return <AccessEditFormBaiduCloudConfig form={configFormInst} model={model?.config as BaiduCloudAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.BYTEPLUS:
        return <AccessEditFormBytePlusConfig form={configFormInst} model={model?.config as BytePlusAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.CLOUDFLARE:
        return <AccessEditFormCloudflareConfig form={configFormInst} model={model?.config as CloudflareAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.DOGECLOUD:
        return <AccessEditFormDogeCloudConfig form={configFormInst} model={model?.config as DogeCloudAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.GODADDY:
        return <AccessEditFormGoDaddyConfig form={configFormInst} model={model?.config as GoDaddyAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.HUAWEICLOUD:
        return <AccessEditFormHuaweiCloudConfig form={configFormInst} model={model?.config as HuaweiCloudAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.KUBERNETES:
        return <AccessEditFormKubernetesConfig form={configFormInst} model={model?.config as KubernetesAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.LOCAL:
        return <AccessEditFormLocalConfig form={configFormInst} model={model?.config as LocalAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.NAMESILO:
        return <AccessEditFormNameSiloConfig form={configFormInst} model={model?.config as NameSiloAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.POWERDNS:
        return <AccessEditFormPowerDNSConfig form={configFormInst} model={model?.config as PowerDNSAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.QINIU:
        return <AccessEditFormQiniuConfig form={configFormInst} model={model?.config as QiniuAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.SSH:
        return <AccessEditFormSSHConfig form={configFormInst} model={model?.config as SSHAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.TENCENTCLOUD:
        return <AccessEditFormTencentCloudConfig form={configFormInst} model={model?.config as TencentCloudAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.VOLCENGINE:
        return <AccessEditFormVolcEngineConfig form={configFormInst} model={model?.config as VolcEngineAccessConfig} />;
      case ACCESS_PROVIDER_TYPES.WEBHOOK:
        return <AccessEditFormWebhookConfig form={configFormInst} model={model?.config as WebhookAccessConfig} />;
    }
  }, [model, configType, configFormInst]);

  const handleFormProviderChange = (name: string) => {
    if (name === "configForm") {
      form.setFieldValue("config", configFormInst.getFieldsValue());
      onModelChange?.(form.getFieldsValue(true));
    }
  };

  const handleFormChange = (_: unknown, fields: AccessEditFormModelType) => {
    if (fields.configType !== configType) {
      setConfigType(fields.configType);
    }

    onModelChange?.(fields);
  };

  useImperativeHandle(ref, () => ({
    getFieldsValue: () => {
      return form.getFieldsValue(true);
    },
    resetFields: () => {
      return form.resetFields();
    },
    validateFields: () => {
      const t1 = form.validateFields();
      const t2 = configFormInst.validateFields();
      return Promise.all([t1, t2]).then(() => t1);
    },
  }));

  return (
    <Form.Provider onFormChange={handleFormProviderChange}>
      <div className={className} style={style}>
        <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" onValuesChange={handleFormChange}>
          <Form.Item name="name" label={t("access.form.name.label")} rules={[formRule]}>
            <Input placeholder={t("access.form.name.placeholder")} />
          </Form.Item>

          <Form.Item
            name="configType"
            label={t("access.form.type.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.type.tooltip") }}></span>}
          >
            <AccessTypeSelect disabled={mode !== "add"} placeholder={t("access.form.type.placeholder")} showSearch={!disabled} />
          </Form.Item>
        </Form>

        {configFormComponent}
      </div>
    </Form.Provider>
  );
});

export default AccessEditForm;
