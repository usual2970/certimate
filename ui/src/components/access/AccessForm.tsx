import { forwardRef, useImperativeHandle, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import AccessProviderSelect from "@/components/provider/AccessProviderSelect";
import { type AccessModel } from "@/domain/access";
import { ACCESS_PROVIDERS, ACCESS_USAGES } from "@/domain/provider";
import { useAntdForm, useAntdFormName } from "@/hooks";

import AccessForm1PanelConfig from "./AccessForm1PanelConfig";
import AccessFormACMEHttpReqConfig from "./AccessFormACMEHttpReqConfig";
import AccessFormAliyunConfig from "./AccessFormAliyunConfig";
import AccessFormAWSConfig from "./AccessFormAWSConfig";
import AccessFormAzureConfig from "./AccessFormAzureConfig";
import AccessFormBaiduCloudConfig from "./AccessFormBaiduCloudConfig";
import AccessFormBaishanConfig from "./AccessFormBaishanConfig";
import AccessFormBaotaPanelConfig from "./AccessFormBaotaPanelConfig";
import AccessFormBunnyConfig from "./AccessFormBunnyConfig";
import AccessFormBytePlusConfig from "./AccessFormBytePlusConfig";
import AccessFormCacheFlyConfig from "./AccessFormCacheFlyConfig";
import AccessFormCdnflyConfig from "./AccessFormCdnflyConfig";
import AccessFormCloudflareConfig from "./AccessFormCloudflareConfig";
import AccessFormClouDNSConfig from "./AccessFormClouDNSConfig";
import AccessFormCMCCCloudConfig from "./AccessFormCMCCCloudConfig";
import AccessFormDeSECConfig from "./AccessFormDeSECConfig";
import AccessFormDingTalkBotConfig from "./AccessFormDingTalkBotConfig";
import AccessFormDNSLAConfig from "./AccessFormDNSLAConfig";
import AccessFormDogeCloudConfig from "./AccessFormDogeCloudConfig";
import AccessFormDynv6Config from "./AccessFormDynv6Config";
import AccessFormEdgioConfig from "./AccessFormEdgioConfig";
import AccessFormEmailConfig from "./AccessFormEmailConfig";
import AccessFormGcoreConfig from "./AccessFormGcoreConfig";
import AccessFormGnameConfig from "./AccessFormGnameConfig";
import AccessFormGoDaddyConfig from "./AccessFormGoDaddyConfig";
import AccessFormGoEdgeConfig from "./AccessFormGoEdgeConfig";
import AccessFormGoogleTrustServicesConfig from "./AccessFormGoogleTrustServicesConfig";
import AccessFormHuaweiCloudConfig from "./AccessFormHuaweiCloudConfig";
import AccessFormJDCloudConfig from "./AccessFormJDCloudConfig";
import AccessFormKubernetesConfig from "./AccessFormKubernetesConfig";
import AccessFormLarkBotConfig from "./AccessFormLarkBotConfig";
import AccessFormMattermostConfig from "./AccessFormMattermostConfig";
import AccessFormNamecheapConfig from "./AccessFormNamecheapConfig";
import AccessFormNameDotComConfig from "./AccessFormNameDotComConfig";
import AccessFormNameSiloConfig from "./AccessFormNameSiloConfig";
import AccessFormNS1Config from "./AccessFormNS1Config";
import AccessFormPorkbunConfig from "./AccessFormPorkbunConfig";
import AccessFormPowerDNSConfig from "./AccessFormPowerDNSConfig";
import AccessFormProxmoxVEConfig from "./AccessFormProxmoxVEConfig";
import AccessFormQiniuConfig from "./AccessFormQiniuConfig";
import AccessFormRainYunConfig from "./AccessFormRainYunConfig";
import AccessFormSafeLineConfig from "./AccessFormSafeLineConfig";
import AccessFormSSHConfig from "./AccessFormSSHConfig";
import AccessFormSSLComConfig from "./AccessFormSSLComConfig";
import AccessFormTelegramConfig from "./AccessFormTelegramConfig";
import AccessFormTencentCloudConfig from "./AccessFormTencentCloudConfig";
import AccessFormUCloudConfig from "./AccessFormUCloudConfig";
import AccessFormUpyunConfig from "./AccessFormUpyunConfig";
import AccessFormVercelConfig from "./AccessFormVercelConfig";
import AccessFormVolcEngineConfig from "./AccessFormVolcEngineConfig";
import AccessFormWangsuConfig from "./AccessFormWangsuConfig";
import AccessFormWebhookConfig from "./AccessFormWebhookConfig";
import AccessFormWeComBotConfig from "./AccessFormWeComBotConfig";
import AccessFormWestcnConfig from "./AccessFormWestcnConfig";
import AccessFormZeroSSLConfig from "./AccessFormZeroSSLConfig";

type AccessFormFieldValues = Partial<MaybeModelRecord<AccessModel>>;
type AccessFormScenes = "add" | "edit";
type AccessFormUsages = "both-dns-hosting" | "ca-only" | "notification-only";

export type AccessFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: AccessFormFieldValues;
  scene: AccessFormScenes;
  usage?: AccessFormUsages;
  onValuesChange?: (values: AccessFormFieldValues) => void;
};

export type AccessFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<AccessFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<AccessFormFieldValues>["resetFields"];
  validateFields: FormInstance<AccessFormFieldValues>["validateFields"];
};

const AccessForm = forwardRef<AccessFormInstance, AccessFormProps>(({ className, style, disabled, initialValues, usage, scene, onValuesChange }, ref) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    name: z
      .string({ message: t("access.form.name.placeholder") })
      .min(1, t("access.form.name.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    provider: z.nativeEnum(ACCESS_PROVIDERS, {
      message:
        usage === "ca-only"
          ? t("access.form.certificate_authority.placeholder")
          : usage === "notification-only"
            ? t("access.form.notification_channel.placeholder")
            : t("access.form.provider.placeholder"),
    }),
    config: z.any(),
    reserve: z.string().nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    initialValues: initialValues,
  });

  const providerLabel = useMemo(() => {
    switch (usage) {
      case "ca-only":
        return t("access.form.certificate_authority.label");
      case "notification-only":
        return t("access.form.notification_channel.label");
    }

    return t("access.form.provider.label");
  }, [usage]);
  const providerPlaceholder = useMemo(() => {
    switch (usage) {
      case "ca-only":
        return t("access.form.certificate_authority.placeholder");
      case "notification-only":
        return t("access.form.notification_channel.placeholder");
    }

    return t("access.form.provider.placeholder");
  }, [usage]);
  const providerTooltip = useMemo(() => {
    switch (usage) {
      case "both-dns-hosting":
        return <span dangerouslySetInnerHTML={{ __html: t("access.form.provider.tooltip") }}></span>;
    }

    return undefined;
  }, [usage]);

  const fieldProvider = Form.useWatch("provider", formInst);

  const [nestedFormInst] = Form.useForm();
  const nestedFormName = useAntdFormName({ form: nestedFormInst, name: "accessEditFormConfigForm" });
  const nestedFormEl = useMemo(() => {
    const nestedFormProps = {
      form: nestedFormInst,
      formName: nestedFormName,
      disabled: disabled,
      initialValues: initialValues?.config,
    };

    /*
      注意：如果追加新的子组件，请保持以 ASCII 排序。
      NOTICE: If you add new child component, please keep ASCII order.
     */
    switch (fieldProvider) {
      case ACCESS_PROVIDERS["1PANEL"]:
        return <AccessForm1PanelConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.ACMEHTTPREQ:
        return <AccessFormACMEHttpReqConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.ALIYUN:
        return <AccessFormAliyunConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.AWS:
        return <AccessFormAWSConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.AZURE:
        return <AccessFormAzureConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.BAIDUCLOUD:
        return <AccessFormBaiduCloudConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.BAISHAN:
        return <AccessFormBaishanConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.BAOTAPANEL:
        return <AccessFormBaotaPanelConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.BUNNY:
        return <AccessFormBunnyConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.BYTEPLUS:
        return <AccessFormBytePlusConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.CACHEFLY:
        return <AccessFormCacheFlyConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.CDNFLY:
        return <AccessFormCdnflyConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.CLOUDFLARE:
        return <AccessFormCloudflareConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.CLOUDNS:
        return <AccessFormClouDNSConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.CMCCCLOUD:
        return <AccessFormCMCCCloudConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.DESEC:
        return <AccessFormDeSECConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.DINGTALKBOT:
        return <AccessFormDingTalkBotConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.DNSLA:
        return <AccessFormDNSLAConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.DOGECLOUD:
        return <AccessFormDogeCloudConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.DYNV6:
        return <AccessFormDynv6Config {...nestedFormProps} />;
      case ACCESS_PROVIDERS.GCORE:
        return <AccessFormGcoreConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.GNAME:
        return <AccessFormGnameConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.GODADDY:
        return <AccessFormGoDaddyConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.GOEDGE:
        return <AccessFormGoEdgeConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.GOOGLETRUSTSERVICES:
        return <AccessFormGoogleTrustServicesConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.EDGIO:
        return <AccessFormEdgioConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.EMAIL:
        return <AccessFormEmailConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.HUAWEICLOUD:
        return <AccessFormHuaweiCloudConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.JDCLOUD:
        return <AccessFormJDCloudConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.KUBERNETES:
        return <AccessFormKubernetesConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.LARKBOT:
        return <AccessFormLarkBotConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.MATTERMOST:
        return <AccessFormMattermostConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.NAMECHEAP:
        return <AccessFormNamecheapConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.NAMEDOTCOM:
        return <AccessFormNameDotComConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.NAMESILO:
        return <AccessFormNameSiloConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.NS1:
        return <AccessFormNS1Config {...nestedFormProps} />;
      case ACCESS_PROVIDERS.PORKBUN:
        return <AccessFormPorkbunConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.POWERDNS:
        return <AccessFormPowerDNSConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.PROXMOXVE:
        return <AccessFormProxmoxVEConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.QINIU:
        return <AccessFormQiniuConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.RAINYUN:
        return <AccessFormRainYunConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.SAFELINE:
        return <AccessFormSafeLineConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.SSH:
        return <AccessFormSSHConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.TELEGRAM:
        return <AccessFormTelegramConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.SSLCOM:
        return <AccessFormSSLComConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.TENCENTCLOUD:
        return <AccessFormTencentCloudConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.UCLOUD:
        return <AccessFormUCloudConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.UPYUN:
        return <AccessFormUpyunConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.VERCEL:
        return <AccessFormVercelConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.VOLCENGINE:
        return <AccessFormVolcEngineConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.WANGSU:
        return <AccessFormWangsuConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.WEBHOOK:
        return (
          <AccessFormWebhookConfig
            usage={usage === "notification-only" ? "notification" : usage === "both-dns-hosting" ? "deployment" : "none"}
            {...nestedFormProps}
          />
        );
      case ACCESS_PROVIDERS.WECOMBOT:
        return <AccessFormWeComBotConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.WESTCN:
        return <AccessFormWestcnConfig {...nestedFormProps} />;
      case ACCESS_PROVIDERS.ZEROSSL:
        return <AccessFormZeroSSLConfig {...nestedFormProps} />;
    }
  }, [disabled, initialValues?.config, fieldProvider, nestedFormInst, nestedFormName]);

  const handleFormProviderChange = (name: string) => {
    if (name === nestedFormName) {
      formInst.setFieldValue("config", nestedFormInst.getFieldsValue());
      onValuesChange?.(formInst.getFieldsValue(true));
    }
  };

  const handleFormChange = (_: unknown, values: AccessFormFieldValues) => {
    onValuesChange?.(values);
  };

  useImperativeHandle(ref, () => {
    return {
      getFieldsValue: () => {
        const values = formInst.getFieldsValue(true);
        values.config = nestedFormInst.getFieldsValue();
        values.reserve = usage === "ca-only" ? "ca" : usage === "notification-only" ? "notification" : undefined;
        return values;
      },
      resetFields: (fields) => {
        formInst.resetFields(fields);

        if (!!fields && fields.includes("config")) {
          nestedFormInst.resetFields(fields);
        }
      },
      validateFields: (nameList, config) => {
        const t1 = formInst.validateFields(nameList, config);
        const t2 = nestedFormInst.validateFields(undefined, config);
        return Promise.all([t1, t2]).then(() => t1);
      },
    } as AccessFormInstance;
  });

  return (
    <Form.Provider onFormChange={handleFormProviderChange}>
      <div className={className} style={style}>
        <Form {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
          <Form.Item name="name" label={t("access.form.name.label")} rules={[formRule]}>
            <Input placeholder={t("access.form.name.placeholder")} />
          </Form.Item>

          <Form.Item name="provider" label={providerLabel} rules={[formRule]} tooltip={providerTooltip}>
            <AccessProviderSelect
              filter={(record) => {
                if (usage == null) return true;

                switch (usage) {
                  case "both-dns-hosting":
                    return record.usages.includes(ACCESS_USAGES.DNS) || record.usages.includes(ACCESS_USAGES.HOSTING);
                  case "ca-only":
                    return record.usages.includes(ACCESS_USAGES.CA);
                  case "notification-only":
                    return record.usages.includes(ACCESS_USAGES.NOTIFICATION);
                }
              }}
              disabled={scene !== "add"}
              placeholder={providerPlaceholder}
              showOptionTags={usage == null || (usage === "both-dns-hosting" ? { [ACCESS_USAGES.DNS]: true, [ACCESS_USAGES.HOSTING]: true } : false)}
              showSearch={!disabled}
            />
          </Form.Item>
        </Form>

        {nestedFormEl}
      </div>
    </Form.Provider>
  );
});

export default AccessForm;
