import { forwardRef, useEffect, useImperativeHandle, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useCreation } from "ahooks";
import { Form, Input, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { useAntdForm } from "@/hooks";
import { ACCESS_PROVIDER_TYPES, type AccessModel } from "@/domain/access";
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

export type AccessEditFormInstance = FormInstance<AccessEditFormFieldValues>;

const AccessEditForm = forwardRef<AccessEditFormInstance, AccessEditFormProps>(({ className, style, disabled, initialValues, preset, onValuesChange }, ref) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    name: z
      .string({ message: t("access.form.name.placeholder") })
      .trim()
      .min(1, t("access.form.name.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    configType: z.nativeEnum(ACCESS_PROVIDER_TYPES, { message: t("access.form.type.placeholder") }),
    config: z.any(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    initialValues: initialValues,
  });

  const [configType, setConfigType] = useState(initialValues?.configType);
  useEffect(() => {
    setConfigType(initialValues?.configType);
  }, [initialValues?.configType]);

  const [configFormInst] = Form.useForm();
  const configFormName = useCreation(() => `accessEditForm_config${Math.random().toString(36).substring(2, 10)}${new Date().getTime()}`, []);
  const configFormComponent = useMemo(() => {
    /*
      注意：如果追加新的子组件，请保持以 ASCII 排序。
      NOTICE: If you add new child component, please keep ASCII order.
     */
    const configFormProps = { form: configFormInst, formName: configFormName, disabled: disabled, initialValues: initialValues?.config };
    switch (configType) {
      case ACCESS_PROVIDER_TYPES.ACMEHTTPREQ:
        return <AccessEditFormACMEHttpReqConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.ALIYUN:
        return <AccessEditFormAliyunConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.AWS:
        return <AccessEditFormAWSConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.BAIDUCLOUD:
        return <AccessEditFormBaiduCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.BYTEPLUS:
        return <AccessEditFormBytePlusConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.CLOUDFLARE:
        return <AccessEditFormCloudflareConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.DOGECLOUD:
        return <AccessEditFormDogeCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.GODADDY:
        return <AccessEditFormGoDaddyConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.HUAWEICLOUD:
        return <AccessEditFormHuaweiCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.KUBERNETES:
        return <AccessEditFormKubernetesConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.LOCAL:
        return <AccessEditFormLocalConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.NAMEDOTCOM:
        return <AccessEditFormNameDotComConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.NAMESILO:
        return <AccessEditFormNameSiloConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.POWERDNS:
        return <AccessEditFormPowerDNSConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.QINIU:
        return <AccessEditFormQiniuConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.SSH:
        return <AccessEditFormSSHConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.TENCENTCLOUD:
        return <AccessEditFormTencentCloudConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.VOLCENGINE:
        return <AccessEditFormVolcEngineConfig {...configFormProps} />;
      case ACCESS_PROVIDER_TYPES.WEBHOOK:
        return <AccessEditFormWebhookConfig {...configFormProps} />;
    }
  }, [disabled, initialValues, configType, configFormInst, configFormName]);

  const handleFormProviderChange = (name: string) => {
    if (name === configFormName) {
      formInst.setFieldValue("config", configFormInst.getFieldsValue());
      onValuesChange?.(formInst.getFieldsValue(true));
    }
  };

  const handleFormChange = (_: unknown, values: AccessEditFormFieldValues) => {
    if (values.configType !== configType) {
      setConfigType(values.configType);
    }

    onValuesChange?.(values);
  };

  useImperativeHandle(ref, () => {
    return {
      getFieldValue: (name) => formInst.getFieldValue(name),
      getFieldsValue: (...args) => {
        if (Array.from(args).length === 0) {
          return formInst.getFieldsValue(true);
        }

        return formInst.getFieldsValue(args[0] as any, args[1] as any);
      },
      getFieldError: (name) => formInst.getFieldError(name),
      getFieldsError: (nameList) => formInst.getFieldsError(nameList),
      getFieldWarning: (name) => formInst.getFieldWarning(name),
      isFieldsTouched: (nameList, allFieldsTouched) => formInst.isFieldsTouched(nameList, allFieldsTouched),
      isFieldTouched: (name) => formInst.isFieldTouched(name),
      isFieldValidating: (name) => formInst.isFieldValidating(name),
      isFieldsValidating: (nameList) => formInst.isFieldsValidating(nameList),
      resetFields: (fields) => formInst.resetFields(fields),
      setFields: (fields) => formInst.setFields(fields),
      setFieldValue: (name, value) => formInst.setFieldValue(name, value),
      setFieldsValue: (values) => formInst.setFieldsValue(values),
      submit: () => formInst.submit(),
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
            name="configType"
            label={t("access.form.type.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.type.tooltip") }}></span>}
          >
            <AccessTypeSelect disabled={preset !== "add"} placeholder={t("access.form.type.placeholder")} showSearch={!disabled} />
          </Form.Item>
        </Form>

        {configFormComponent}
      </div>
    </Form.Provider>
  );
});

export default AccessEditForm;
