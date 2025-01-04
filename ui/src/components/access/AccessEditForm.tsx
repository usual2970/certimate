import { forwardRef, useEffect, useImperativeHandle, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import AccessProviderSelect from "@/components/provider/AccessProviderSelect";
import { type AccessModel } from "@/domain/access";
import { ACCESS_PROVIDERS } from "@/domain/provider";
import { useAntdForm, useAntdFormName } from "@/hooks";

import AccessEditFormACMEHttpReqConfig from "./AccessEditFormACMEHttpReqConfig";
import AccessEditFormAWSConfig from "./AccessEditFormAWSConfig";
import AccessEditFormAliyunConfig from "./AccessEditFormAliyunConfig";
import AccessEditFormBaiduCloudConfig from "./AccessEditFormBaiduCloudConfig";
import AccessEditFormBytePlusConfig from "./AccessEditFormBytePlusConfig";
import AccessEditFormCloudflareConfig from "./AccessEditFormCloudflareConfig";
import AccessEditFormDogeCloudConfig from "./AccessEditFormDogeCloudConfig";
import AccessEditFormGoDaddyConfig from "./AccessEditFormGoDaddyConfig";
import AccessEditFormHuaweiCloudConfig from "./AccessEditFormHuaweiCloudConfig";
import AccessEditFormKubernetesConfig from "./AccessEditFormKubernetesConfig";
import AccessEditFormLocalConfig from "./AccessEditFormLocalConfig";
import AccessEditFormNameDotComConfig from "./AccessEditFormNameDotComConfig";
import AccessEditFormNameSiloConfig from "./AccessEditFormNameSiloConfig";
import AccessEditFormPowerDNSConfig from "./AccessEditFormPowerDNSConfig";
import AccessEditFormQiniuConfig from "./AccessEditFormQiniuConfig";
import AccessEditFormSSHConfig from "./AccessEditFormSSHConfig";
import AccessEditFormTencentCloudConfig from "./AccessEditFormTencentCloudConfig";
import AccessEditFormVolcEngineConfig from "./AccessEditFormVolcEngineConfig";
import AccessEditFormWebhookConfig from "./AccessEditFormWebhookConfig";

type AccessEditFormFieldValues = Partial<MaybeModelRecord<AccessModel>>;
type AccessEditFormPresets = "add" | "edit";

export type AccessEditFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: AccessEditFormFieldValues;
  preset: AccessEditFormPresets;
  onValuesChange?: (values: AccessEditFormFieldValues) => void;
};

export type AccessEditFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<AccessEditFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<AccessEditFormFieldValues>["resetFields"];
  validateFields: FormInstance<AccessEditFormFieldValues>["validateFields"];
};

const AccessEditForm = forwardRef<AccessEditFormInstance, AccessEditFormProps>(({ className, style, disabled, initialValues, preset, onValuesChange }, ref) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    name: z
      .string({ message: t("access.form.name.placeholder") })
      .min(1, t("access.form.name.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    provider: z.nativeEnum(ACCESS_PROVIDERS, { message: t("access.form.provider.placeholder") }),
    config: z.any(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    initialValues: initialValues,
  });

  const [fieldProvider, setFieldProvider] = useState(initialValues?.provider);
  useEffect(() => {
    setFieldProvider(initialValues?.provider);
  }, [initialValues?.provider]);

  const [configFormInst] = Form.useForm();
  const configFormName = useAntdFormName({ form: configFormInst, name: "accessEditConfigForm" });
  const configFormComponent = useMemo(() => {
    /*
      注意：如果追加新的子组件，请保持以 ASCII 排序。
      NOTICE: If you add new child component, please keep ASCII order.
     */
    const configFormProps = { form: configFormInst, formName: configFormName, disabled: disabled, initialValues: initialValues?.config };
    switch (fieldProvider) {
      case ACCESS_PROVIDERS.ACMEHTTPREQ:
        return <AccessEditFormACMEHttpReqConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.ALIYUN:
        return <AccessEditFormAliyunConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.AWS:
        return <AccessEditFormAWSConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.BAIDUCLOUD:
        return <AccessEditFormBaiduCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.BYTEPLUS:
        return <AccessEditFormBytePlusConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.CLOUDFLARE:
        return <AccessEditFormCloudflareConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.DOGECLOUD:
        return <AccessEditFormDogeCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.GODADDY:
        return <AccessEditFormGoDaddyConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.HUAWEICLOUD:
        return <AccessEditFormHuaweiCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.KUBERNETES:
        return <AccessEditFormKubernetesConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.LOCAL:
        return <AccessEditFormLocalConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.NAMEDOTCOM:
        return <AccessEditFormNameDotComConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.NAMESILO:
        return <AccessEditFormNameSiloConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.POWERDNS:
        return <AccessEditFormPowerDNSConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.QINIU:
        return <AccessEditFormQiniuConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.SSH:
        return <AccessEditFormSSHConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.TENCENTCLOUD:
        return <AccessEditFormTencentCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.VOLCENGINE:
        return <AccessEditFormVolcEngineConfig {...configFormProps} />;
      case ACCESS_PROVIDERS.WEBHOOK:
        return <AccessEditFormWebhookConfig {...configFormProps} />;
    }
  }, [disabled, initialValues, fieldProvider, configFormInst, configFormName]);

  const handleFormProviderChange = (name: string) => {
    if (name === configFormName) {
      formInst.setFieldValue("config", configFormInst.getFieldsValue());
      onValuesChange?.(formInst.getFieldsValue(true));
    }
  };

  const handleFormChange = (_: unknown, values: AccessEditFormFieldValues) => {
    if (values.provider !== fieldProvider) {
      setFieldProvider(values.provider);
    }

    onValuesChange?.(values);
  };

  useImperativeHandle(ref, () => {
    return {
      getFieldsValue: () => {
        return formInst.getFieldsValue(true);
      },
      resetFields: (fields) => {
        return formInst.resetFields(fields);
      },
      validateFields: (nameList, config) => {
        const t1 = formInst.validateFields(nameList, config);
        const t2 = configFormInst.validateFields(undefined, config);
        return Promise.all([t1, t2]).then(() => t1);
      },
    } as AccessEditFormInstance;
  });

  return (
    <Form.Provider onFormChange={handleFormProviderChange}>
      <div className={className} style={style}>
        <Form {...formProps} disabled={disabled} layout="vertical" onValuesChange={handleFormChange}>
          <Form.Item name="name" label={t("access.form.name.label")} rules={[formRule]}>
            <Input placeholder={t("access.form.name.placeholder")} />
          </Form.Item>

          <Form.Item
            name="provider"
            label={t("access.form.provider.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.provider.tooltip") }}></span>}
          >
            <AccessProviderSelect disabled={preset !== "add"} placeholder={t("access.form.provider.placeholder")} showSearch={!disabled} />
          </Form.Item>
        </Form>

        {configFormComponent}
      </div>
    </Form.Provider>
  );
});

export default AccessEditForm;
