import { forwardRef, memo, useEffect, useImperativeHandle, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { PlusOutlined as PlusOutlinedIcon, QuestionCircleOutlined as QuestionCircleOutlinedIcon } from "@ant-design/icons";
import { Alert, Button, Divider, Flex, Form, type FormInstance, Select, Switch, Tooltip, Typography } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import AccessEditModal from "@/components/access/AccessEditModal";
import AccessSelect from "@/components/access/AccessSelect";
import DeployProviderPicker from "@/components/provider/DeployProviderPicker";
import DeployProviderSelect from "@/components/provider/DeployProviderSelect";
import Show from "@/components/Show";
import { ACCESS_USAGES, DEPLOY_PROVIDERS, accessProvidersMap, deployProvidersMap } from "@/domain/provider";
import { type WorkflowNode, type WorkflowNodeConfigForDeploy } from "@/domain/workflow";
import { useAntdForm, useAntdFormName, useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import DeployNodeConfigForm1PanelConsoleConfig from "./DeployNodeConfigForm1PanelConsoleConfig";
import DeployNodeConfigForm1PanelSiteConfig from "./DeployNodeConfigForm1PanelSiteConfig";
import DeployNodeConfigFormAliyunALBConfig from "./DeployNodeConfigFormAliyunALBConfig";
import DeployNodeConfigFormAliyunCASConfig from "./DeployNodeConfigFormAliyunCASConfig";
import DeployNodeConfigFormAliyunCASDeployConfig from "./DeployNodeConfigFormAliyunCASDeployConfig";
import DeployNodeConfigFormAliyunCDNConfig from "./DeployNodeConfigFormAliyunCDNConfig";
import DeployNodeConfigFormAliyunCLBConfig from "./DeployNodeConfigFormAliyunCLBConfig";
import DeployNodeConfigFormAliyunDCDNConfig from "./DeployNodeConfigFormAliyunDCDNConfig";
import DeployNodeConfigFormAliyunESAConfig from "./DeployNodeConfigFormAliyunESAConfig";
import DeployNodeConfigFormAliyunFCConfig from "./DeployNodeConfigFormAliyunFCConfig";
import DeployNodeConfigFormAliyunLiveConfig from "./DeployNodeConfigFormAliyunLiveConfig";
import DeployNodeConfigFormAliyunNLBConfig from "./DeployNodeConfigFormAliyunNLBConfig";
import DeployNodeConfigFormAliyunOSSConfig from "./DeployNodeConfigFormAliyunOSSConfig";
import DeployNodeConfigFormAliyunVODConfig from "./DeployNodeConfigFormAliyunVODConfig";
import DeployNodeConfigFormAliyunWAFConfig from "./DeployNodeConfigFormAliyunWAFConfig";
import DeployNodeConfigFormAWSACMConfig from "./DeployNodeConfigFormAWSACMConfig";
import DeployNodeConfigFormAWSCloudFrontConfig from "./DeployNodeConfigFormAWSCloudFrontConfig";
import DeployNodeConfigFormAzureKeyVaultConfig from "./DeployNodeConfigFormAzureKeyVaultConfig";
import DeployNodeConfigFormBaiduCloudAppBLBConfig from "./DeployNodeConfigFormBaiduCloudAppBLBConfig";
import DeployNodeConfigFormBaiduCloudBLBConfig from "./DeployNodeConfigFormBaiduCloudBLBConfig";
import DeployNodeConfigFormBaiduCloudCDNConfig from "./DeployNodeConfigFormBaiduCloudCDNConfig";
import DeployNodeConfigFormBaishanCDNConfig from "./DeployNodeConfigFormBaishanCDNConfig";
import DeployNodeConfigFormBaotaPanelConsoleConfig from "./DeployNodeConfigFormBaotaPanelConsoleConfig";
import DeployNodeConfigFormBaotaPanelSiteConfig from "./DeployNodeConfigFormBaotaPanelSiteConfig";
import DeployNodeConfigFormBytePlusCDNConfig from "./DeployNodeConfigFormBytePlusCDNConfig";
import DeployNodeConfigFormCdnflyConfig from "./DeployNodeConfigFormCdnflyConfig";
import DeployNodeConfigFormDogeCloudCDNConfig from "./DeployNodeConfigFormDogeCloudCDNConfig";
import DeployNodeConfigFormEdgioApplicationsConfig from "./DeployNodeConfigFormEdgioApplicationsConfig";
import DeployNodeConfigFormGcoreCDNConfig from "./DeployNodeConfigFormGcoreCDNConfig";
import DeployNodeConfigFormHuaweiCloudCDNConfig from "./DeployNodeConfigFormHuaweiCloudCDNConfig";
import DeployNodeConfigFormHuaweiCloudELBConfig from "./DeployNodeConfigFormHuaweiCloudELBConfig";
import DeployNodeConfigFormHuaweiCloudWAFConfig from "./DeployNodeConfigFormHuaweiCloudWAFConfig";
import DeployNodeConfigFormJDCloudALBConfig from "./DeployNodeConfigFormJDCloudALBConfig";
import DeployNodeConfigFormJDCloudCDNConfig from "./DeployNodeConfigFormJDCloudCDNConfig";
import DeployNodeConfigFormJDCloudLiveConfig from "./DeployNodeConfigFormJDCloudLiveConfig";
import DeployNodeConfigFormJDCloudVODConfig from "./DeployNodeConfigFormJDCloudVODConfig";
import DeployNodeConfigFormKubernetesSecretConfig from "./DeployNodeConfigFormKubernetesSecretConfig";
import DeployNodeConfigFormLocalConfig from "./DeployNodeConfigFormLocalConfig";
import DeployNodeConfigFormQiniuCDNConfig from "./DeployNodeConfigFormQiniuCDNConfig";
import DeployNodeConfigFormQiniuKodoConfig from "./DeployNodeConfigFormQiniuKodoConfig";
import DeployNodeConfigFormQiniuPiliConfig from "./DeployNodeConfigFormQiniuPiliConfig";
import DeployNodeConfigFormSafeLineConfig from "./DeployNodeConfigFormSafeLineConfig";
import DeployNodeConfigFormSSHConfig from "./DeployNodeConfigFormSSHConfig.tsx";
import DeployNodeConfigFormTencentCloudCDNConfig from "./DeployNodeConfigFormTencentCloudCDNConfig.tsx";
import DeployNodeConfigFormTencentCloudCLBConfig from "./DeployNodeConfigFormTencentCloudCLBConfig.tsx";
import DeployNodeConfigFormTencentCloudCOSConfig from "./DeployNodeConfigFormTencentCloudCOSConfig.tsx";
import DeployNodeConfigFormTencentCloudCSSConfig from "./DeployNodeConfigFormTencentCloudCSSConfig.tsx";
import DeployNodeConfigFormTencentCloudECDNConfig from "./DeployNodeConfigFormTencentCloudECDNConfig.tsx";
import DeployNodeConfigFormTencentCloudEOConfig from "./DeployNodeConfigFormTencentCloudEOConfig.tsx";
import DeployNodeConfigFormTencentCloudSCFConfig from "./DeployNodeConfigFormTencentCloudSCFConfig";
import DeployNodeConfigFormTencentCloudSSLDeployConfig from "./DeployNodeConfigFormTencentCloudSSLDeployConfig";
import DeployNodeConfigFormTencentCloudVODConfig from "./DeployNodeConfigFormTencentCloudVODConfig";
import DeployNodeConfigFormTencentCloudWAFConfig from "./DeployNodeConfigFormTencentCloudWAFConfig";
import DeployNodeConfigFormUCloudUCDNConfig from "./DeployNodeConfigFormUCloudUCDNConfig.tsx";
import DeployNodeConfigFormUCloudUS3Config from "./DeployNodeConfigFormUCloudUS3Config.tsx";
import DeployNodeConfigFormUpyunCDNConfig from "./DeployNodeConfigFormUpyunCDNConfig.tsx";
import DeployNodeConfigFormUpyunFileConfig from "./DeployNodeConfigFormUpyunFileConfig.tsx";
import DeployNodeConfigFormVolcEngineALBConfig from "./DeployNodeConfigFormVolcEngineALBConfig.tsx";
import DeployNodeConfigFormVolcEngineCDNConfig from "./DeployNodeConfigFormVolcEngineCDNConfig.tsx";
import DeployNodeConfigFormVolcEngineCertCenterConfig from "./DeployNodeConfigFormVolcEngineCertCenterConfig.tsx";
import DeployNodeConfigFormVolcEngineCLBConfig from "./DeployNodeConfigFormVolcEngineCLBConfig.tsx";
import DeployNodeConfigFormVolcEngineDCDNConfig from "./DeployNodeConfigFormVolcEngineDCDNConfig.tsx";
import DeployNodeConfigFormVolcEngineImageXConfig from "./DeployNodeConfigFormVolcEngineImageXConfig.tsx";
import DeployNodeConfigFormVolcEngineLiveConfig from "./DeployNodeConfigFormVolcEngineLiveConfig.tsx";
import DeployNodeConfigFormVolcEngineTOSConfig from "./DeployNodeConfigFormVolcEngineTOSConfig.tsx";
import DeployNodeConfigFormWebhookConfig from "./DeployNodeConfigFormWebhookConfig.tsx";

type DeployNodeConfigFormFieldValues = Partial<WorkflowNodeConfigForDeploy>;

export type DeployNodeConfigFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: DeployNodeConfigFormFieldValues;
  nodeId: string;
  onValuesChange?: (values: DeployNodeConfigFormFieldValues) => void;
};

export type DeployNodeConfigFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<DeployNodeConfigFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<DeployNodeConfigFormFieldValues>["resetFields"];
  validateFields: FormInstance<DeployNodeConfigFormFieldValues>["validateFields"];
};

const initFormModel = (): DeployNodeConfigFormFieldValues => {
  return {
    skipOnLastSucceeded: true,
  };
};

const DeployNodeConfigForm = forwardRef<DeployNodeConfigFormInstance, DeployNodeConfigFormProps>(
  ({ className, style, disabled, initialValues, nodeId, onValuesChange }, ref) => {
    const { t } = useTranslation();

    const { getWorkflowOuptutBeforeId } = useWorkflowStore(useZustandShallowSelector(["updateNode", "getWorkflowOuptutBeforeId"]));

    // TODO: 优化此处逻辑
    const [previousNodes, setPreviousNodes] = useState<WorkflowNode[]>([]);
    useEffect(() => {
      const previousNodes = getWorkflowOuptutBeforeId(nodeId, "certificate");
      setPreviousNodes(previousNodes);
    }, [nodeId]);

    const formSchema = z.object({
      certificate: z
        .string({ message: t("workflow_node.deploy.form.certificate.placeholder") })
        .nonempty(t("workflow_node.deploy.form.certificate.placeholder")),
      provider: z.string({ message: t("workflow_node.deploy.form.provider.placeholder") }).nonempty(t("workflow_node.deploy.form.provider.placeholder")),
      providerAccessId: z
        .string({ message: t("workflow_node.deploy.form.provider_access.placeholder") })
        .nonempty(t("workflow_node.deploy.form.provider_access.placeholder")),
      providerConfig: z.any(),
      skipOnLastSucceeded: z.boolean().nullish(),
    });
    const formRule = createSchemaFieldRule(formSchema);
    const { form: formInst, formProps } = useAntdForm({
      name: "workflowNodeDeployConfigForm",
      initialValues: initialValues ?? initFormModel(),
    });

    const fieldProvider = Form.useWatch("provider", { form: formInst, preserve: true });

    const [nestedFormInst] = Form.useForm();
    const nestedFormName = useAntdFormName({ form: nestedFormInst, name: "workflowNodeDeployConfigFormProviderConfigForm" });
    const nestedFormEl = useMemo(() => {
      const nestedFormProps = {
        form: nestedFormInst,
        formName: nestedFormName,
        disabled: disabled,
        initialValues: initialValues?.providerConfig,
      };

      /*
        注意：如果追加新的子组件，请保持以 ASCII 排序。
        NOTICE: If you add new child component, please keep ASCII order.
       */
      switch (fieldProvider) {
        case DEPLOY_PROVIDERS["1PANEL_CONSOLE"]:
          return <DeployNodeConfigForm1PanelConsoleConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS["1PANEL_SITE"]:
          return <DeployNodeConfigForm1PanelSiteConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_ALB:
          return <DeployNodeConfigFormAliyunALBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_CAS:
          return <DeployNodeConfigFormAliyunCASConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_CAS_DEPLOY:
          return <DeployNodeConfigFormAliyunCASDeployConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_CLB:
          return <DeployNodeConfigFormAliyunCLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_CDN:
          return <DeployNodeConfigFormAliyunCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_DCDN:
          return <DeployNodeConfigFormAliyunDCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_ESA:
          return <DeployNodeConfigFormAliyunESAConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_FC:
          return <DeployNodeConfigFormAliyunFCConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_LIVE:
          return <DeployNodeConfigFormAliyunLiveConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_NLB:
          return <DeployNodeConfigFormAliyunNLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_OSS:
          return <DeployNodeConfigFormAliyunOSSConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_VOD:
          return <DeployNodeConfigFormAliyunVODConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_WAF:
          return <DeployNodeConfigFormAliyunWAFConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.AWS_ACM:
          return <DeployNodeConfigFormAWSACMConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.AWS_CLOUDFRONT:
          return <DeployNodeConfigFormAWSCloudFrontConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.AZURE_KEYVAULT:
          return <DeployNodeConfigFormAzureKeyVaultConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BAIDUCLOUD_APPBLB:
          return <DeployNodeConfigFormBaiduCloudAppBLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BAIDUCLOUD_BLB:
          return <DeployNodeConfigFormBaiduCloudBLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BAIDUCLOUD_CDN:
          return <DeployNodeConfigFormBaiduCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BAISHAN_CDN:
          return <DeployNodeConfigFormBaishanCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BAOTAPANEL_CONSOLE:
          return <DeployNodeConfigFormBaotaPanelConsoleConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BAOTAPANEL_SITE:
          return <DeployNodeConfigFormBaotaPanelSiteConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BYTEPLUS_CDN:
          return <DeployNodeConfigFormBytePlusCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.CDNFLY:
          return <DeployNodeConfigFormCdnflyConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.DOGECLOUD_CDN:
          return <DeployNodeConfigFormDogeCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.EDGIO_APPLICATIONS:
          return <DeployNodeConfigFormEdgioApplicationsConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.GCORE_CDN:
          return <DeployNodeConfigFormGcoreCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.HUAWEICLOUD_CDN:
          return <DeployNodeConfigFormHuaweiCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.HUAWEICLOUD_ELB:
          return <DeployNodeConfigFormHuaweiCloudELBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.HUAWEICLOUD_WAF:
          return <DeployNodeConfigFormHuaweiCloudWAFConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.JDCLOUD_ALB:
          return <DeployNodeConfigFormJDCloudALBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.JDCLOUD_CDN:
          return <DeployNodeConfigFormJDCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.JDCLOUD_LIVE:
          return <DeployNodeConfigFormJDCloudLiveConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.JDCLOUD_VOD:
          return <DeployNodeConfigFormJDCloudVODConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.KUBERNETES_SECRET:
          return <DeployNodeConfigFormKubernetesSecretConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.LOCAL:
          return <DeployNodeConfigFormLocalConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.QINIU_CDN:
          return <DeployNodeConfigFormQiniuCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.QINIU_KODO:
          return <DeployNodeConfigFormQiniuKodoConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.QINIU_PILI:
          return <DeployNodeConfigFormQiniuPiliConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.SAFELINE:
          return <DeployNodeConfigFormSafeLineConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.SSH:
          return <DeployNodeConfigFormSSHConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_CDN:
          return <DeployNodeConfigFormTencentCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_CLB:
          return <DeployNodeConfigFormTencentCloudCLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_COS:
          return <DeployNodeConfigFormTencentCloudCOSConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_CSS:
          return <DeployNodeConfigFormTencentCloudCSSConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_ECDN:
          return <DeployNodeConfigFormTencentCloudECDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_EO:
          return <DeployNodeConfigFormTencentCloudEOConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_SCF:
          return <DeployNodeConfigFormTencentCloudSCFConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_SSL_DEPLOY:
          return <DeployNodeConfigFormTencentCloudSSLDeployConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_VOD:
          return <DeployNodeConfigFormTencentCloudVODConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.TENCENTCLOUD_WAF:
          return <DeployNodeConfigFormTencentCloudWAFConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.UCLOUD_UCDN:
          return <DeployNodeConfigFormUCloudUCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.UCLOUD_US3:
          return <DeployNodeConfigFormUCloudUS3Config {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.UPYUN_CDN:
          return <DeployNodeConfigFormUpyunCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.UPYUN_FILE:
          return <DeployNodeConfigFormUpyunFileConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_ALB:
          return <DeployNodeConfigFormVolcEngineALBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_CDN:
          return <DeployNodeConfigFormVolcEngineCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_CERTCENTER:
          return <DeployNodeConfigFormVolcEngineCertCenterConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_CLB:
          return <DeployNodeConfigFormVolcEngineCLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_DCDN:
          return <DeployNodeConfigFormVolcEngineDCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_IMAGEX:
          return <DeployNodeConfigFormVolcEngineImageXConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_LIVE:
          return <DeployNodeConfigFormVolcEngineLiveConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_TOS:
          return <DeployNodeConfigFormVolcEngineTOSConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.WEBHOOK:
          return <DeployNodeConfigFormWebhookConfig {...nestedFormProps} />;
      }
    }, [disabled, initialValues?.providerConfig, fieldProvider, nestedFormInst, nestedFormName]);

    const handleProviderPick = (value: string) => {
      formInst.setFieldValue("provider", value);
      onValuesChange?.(formInst.getFieldsValue(true));
    };

    const handleProviderSelect = (value: string) => {
      if (fieldProvider === value) return;

      // 切换部署目标时重置表单，避免其他部署目标的配置字段影响当前部署目标
      if (initialValues?.provider === value) {
        formInst.resetFields();
      } else {
        const oldValues = formInst.getFieldsValue();
        const newValues: Record<string, unknown> = {};
        for (const key in oldValues) {
          if (key === "provider" || key === "providerAccessId" || key === "certificate" || key === "skipOnLastSucceeded") {
            newValues[key] = oldValues[key];
          } else {
            newValues[key] = undefined;
          }
        }
        formInst.setFieldsValue(newValues);

        if (deployProvidersMap.get(fieldProvider)?.provider !== deployProvidersMap.get(value)?.provider) {
          formInst.setFieldValue("providerAccessId", undefined);
          onValuesChange?.(formInst.getFieldsValue(true));
        }
      }
    };

    const handleFormProviderChange = (name: string) => {
      if (name === nestedFormName) {
        formInst.setFieldValue("providerConfig", nestedFormInst.getFieldsValue());
        onValuesChange?.(formInst.getFieldsValue(true));
      }
    };

    const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
      onValuesChange?.(values as DeployNodeConfigFormFieldValues);
    };

    useImperativeHandle(ref, () => {
      return {
        getFieldsValue: () => {
          const values = formInst.getFieldsValue(true);
          values.providerConfig = nestedFormInst.getFieldsValue();
          return values;
        },
        resetFields: (fields) => {
          formInst.resetFields(fields);

          if (!!fields && fields.includes("providerConfig")) {
            nestedFormInst.resetFields(fields);
          }
        },
        validateFields: (nameList, config) => {
          const t1 = formInst.validateFields(nameList, config);
          const t2 = nestedFormInst.validateFields(undefined, config);
          return Promise.all([t1, t2]).then(() => t1);
        },
      } as DeployNodeConfigFormInstance;
    });

    return (
      <Form.Provider onFormChange={handleFormProviderChange}>
        <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
          <Show
            when={!!fieldProvider}
            fallback={<DeployProviderPicker autoFocus placeholder={t("workflow_node.deploy.search.provider.placeholder")} onSelect={handleProviderPick} />}
          >
            <Form.Item name="provider" label={t("workflow_node.deploy.form.provider.label")} rules={[formRule]}>
              <DeployProviderSelect
                allowClear
                disabled={!!initialValues?.provider}
                placeholder={t("workflow_node.deploy.form.provider.placeholder")}
                showSearch
                onSelect={handleProviderSelect}
              />
            </Form.Item>

            <Form.Item className="mb-0">
              <label className="mb-1 block">
                <div className="flex w-full items-center justify-between gap-4">
                  <div className="max-w-full grow truncate">
                    <span>{t("workflow_node.deploy.form.provider_access.label")}</span>
                    <Tooltip title={t("workflow_node.deploy.form.provider_access.tooltip")}>
                      <Typography.Text className="ms-1" type="secondary">
                        <QuestionCircleOutlinedIcon />
                      </Typography.Text>
                    </Tooltip>
                  </div>
                  <div className="text-right">
                    <AccessEditModal
                      data={{ provider: deployProvidersMap.get(fieldProvider!)?.provider }}
                      preset="add"
                      trigger={
                        <Button size="small" type="link">
                          <PlusOutlinedIcon />
                          {t("workflow_node.deploy.form.provider_access.button")}
                        </Button>
                      }
                      afterSubmit={(record) => {
                        const provider = accessProvidersMap.get(record.provider);
                        if (provider?.usages?.includes(ACCESS_USAGES.DEPLOY)) {
                          formInst.setFieldValue("providerAccessId", record.id);
                        }
                      }}
                    />
                  </div>
                </div>
              </label>
              <Form.Item name="providerAccessId" rules={[formRule]}>
                <AccessSelect
                  filter={(record) => {
                    if (fieldProvider) {
                      return deployProvidersMap.get(fieldProvider)?.provider === record.provider;
                    }

                    const provider = accessProvidersMap.get(record.provider);
                    return !!provider?.usages?.includes(ACCESS_USAGES.DEPLOY);
                  }}
                  placeholder={t("workflow_node.deploy.form.provider_access.placeholder")}
                />
              </Form.Item>
            </Form.Item>

            <Show when={fieldProvider === DEPLOY_PROVIDERS.LOCAL}>
              <Form.Item>
                <Alert
                  type="info"
                  message={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.provider_access.guide_for_local") }}></span>}
                />
              </Form.Item>
            </Show>

            <Form.Item
              name="certificate"
              label={t("workflow_node.deploy.form.certificate.label")}
              rules={[formRule]}
              tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.certificate.tooltip") }}></span>}
            >
              <Select
                options={previousNodes.map((item) => {
                  return {
                    label: item.name,
                    options: item.outputs?.map((output) => {
                      return {
                        label: `${item.name} - ${output.label}`,
                        value: `${item.id}#${output.name}`,
                      };
                    }),
                  };
                })}
                placeholder={t("workflow_node.deploy.form.certificate.placeholder")}
              />
            </Form.Item>
          </Show>
        </Form>

        <Show when={!!nestedFormEl}>
          <Divider className="my-1">
            <Typography.Text className="text-xs font-normal" type="secondary">
              {t("workflow_node.deploy.form.params_config.label")}
            </Typography.Text>
          </Divider>

          {nestedFormEl}
        </Show>

        <Show when={!!fieldProvider}>
          <Divider className="my-1">
            <Typography.Text className="text-xs font-normal" type="secondary">
              {t("workflow_node.deploy.form.strategy_config.label")}
            </Typography.Text>
          </Divider>

          <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
            <Form.Item label={t("workflow_node.deploy.form.skip_on_last_succeeded.label")}>
              <Flex align="center" gap={8} wrap="wrap">
                <div>{t("workflow_node.deploy.form.skip_on_last_succeeded.prefix")}</div>
                <Form.Item name="skipOnLastSucceeded" noStyle rules={[formRule]}>
                  <Switch
                    checkedChildren={t("workflow_node.deploy.form.skip_on_last_succeeded.switch.on")}
                    unCheckedChildren={t("workflow_node.deploy.form.skip_on_last_succeeded.switch.off")}
                  />
                </Form.Item>
                <div>{t("workflow_node.deploy.form.skip_on_last_succeeded.suffix")}</div>
              </Flex>
            </Form.Item>
          </Form>
        </Show>
      </Form.Provider>
    );
  }
);

export default memo(DeployNodeConfigForm);
