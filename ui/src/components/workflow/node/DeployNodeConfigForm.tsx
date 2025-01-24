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

import DeployNodeConfigFormAliyunALBConfig from "./DeployNodeConfigFormAliyunALBConfig";
import DeployNodeConfigFormAliyunCDNConfig from "./DeployNodeConfigFormAliyunCDNConfig";
import DeployNodeConfigFormAliyunCLBConfig from "./DeployNodeConfigFormAliyunCLBConfig";
import DeployNodeConfigFormAliyunDCDNConfig from "./DeployNodeConfigFormAliyunDCDNConfig";
import DeployNodeConfigFormAliyunLiveConfig from "./DeployNodeConfigFormAliyunLiveConfig";
import DeployNodeConfigFormAliyunNLBConfig from "./DeployNodeConfigFormAliyunNLBConfig";
import DeployNodeConfigFormAliyunOSSConfig from "./DeployNodeConfigFormAliyunOSSConfig";
import DeployNodeConfigFormAliyunWAFConfig from "./DeployNodeConfigFormAliyunWAFConfig";
import DeployNodeConfigFormAWSCloudFrontConfig from "./DeployNodeConfigFormAWSCloudFrontConfig";
import DeployNodeConfigFormBaiduCloudCDNConfig from "./DeployNodeConfigFormBaiduCloudCDNConfig";
import DeployNodeConfigFormBytePlusCDNConfig from "./DeployNodeConfigFormBytePlusCDNConfig";
import DeployNodeConfigFormDogeCloudCDNConfig from "./DeployNodeConfigFormDogeCloudCDNConfig";
import DeployNodeConfigFormEdgioApplicationsConfig from "./DeployNodeConfigFormEdgioApplicationsConfig";
import DeployNodeConfigFormHuaweiCloudCDNConfig from "./DeployNodeConfigFormHuaweiCloudCDNConfig";
import DeployNodeConfigFormHuaweiCloudELBConfig from "./DeployNodeConfigFormHuaweiCloudELBConfig";
import DeployNodeConfigFormKubernetesSecretConfig from "./DeployNodeConfigFormKubernetesSecretConfig";
import DeployNodeConfigFormLocalConfig from "./DeployNodeConfigFormLocalConfig";
import DeployNodeConfigFormQiniuCDNConfig from "./DeployNodeConfigFormQiniuCDNConfig";
import DeployNodeConfigFormQiniuPiliConfig from "./DeployNodeConfigFormQiniuPiliConfig";
import DeployNodeConfigFormSSHConfig from "./DeployNodeConfigFormSSHConfig.tsx";
import DeployNodeConfigFormTencentCloudCDNConfig from "./DeployNodeConfigFormTencentCloudCDNConfig.tsx";
import DeployNodeConfigFormTencentCloudCLBConfig from "./DeployNodeConfigFormTencentCloudCLBConfig.tsx";
import DeployNodeConfigFormTencentCloudCOSConfig from "./DeployNodeConfigFormTencentCloudCOSConfig.tsx";
import DeployNodeConfigFormTencentCloudCSSConfig from "./DeployNodeConfigFormTencentCloudCSSConfig.tsx";
import DeployNodeConfigFormTencentCloudECDNConfig from "./DeployNodeConfigFormTencentCloudECDNConfig.tsx";
import DeployNodeConfigFormTencentCloudEOConfig from "./DeployNodeConfigFormTencentCloudEOConfig.tsx";
import DeployNodeConfigFormUCloudUCDNConfig from "./DeployNodeConfigFormUCloudUCDNConfig.tsx";
import DeployNodeConfigFormUCloudUS3Config from "./DeployNodeConfigFormUCloudUS3Config.tsx";
import DeployNodeConfigFormVolcEngineCDNConfig from "./DeployNodeConfigFormVolcEngineCDNConfig.tsx";
import DeployNodeConfigFormVolcEngineCLBConfig from "./DeployNodeConfigFormVolcEngineCLBConfig.tsx";
import DeployNodeConfigFormVolcEngineDCDNConfig from "./DeployNodeConfigFormVolcEngineDCDNConfig.tsx";
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
        .nonempty(t("workflow_node.deploy.form.provider_access.placeholder"))
        .refine(() => !!formInst.getFieldValue("provider"), t("workflow_node.deploy.form.provider.placeholder")),
      providerConfig: z.any(),
      skipOnLastSucceeded: z.boolean(),
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
        case DEPLOY_PROVIDERS.ALIYUN_ALB:
          return <DeployNodeConfigFormAliyunALBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_CLB:
          return <DeployNodeConfigFormAliyunCLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_CDN:
          return <DeployNodeConfigFormAliyunCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_DCDN:
          return <DeployNodeConfigFormAliyunDCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_LIVE:
          return <DeployNodeConfigFormAliyunLiveConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_NLB:
          return <DeployNodeConfigFormAliyunNLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_OSS:
          return <DeployNodeConfigFormAliyunOSSConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.ALIYUN_WAF:
          return <DeployNodeConfigFormAliyunWAFConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.AWS_CLOUDFRONT:
          return <DeployNodeConfigFormAWSCloudFrontConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BAIDUCLOUD_CDN:
          return <DeployNodeConfigFormBaiduCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.BYTEPLUS_CDN:
          return <DeployNodeConfigFormBytePlusCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.DOGECLOUD_CDN:
          return <DeployNodeConfigFormDogeCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.EDGIO_APPLICATIONS:
          return <DeployNodeConfigFormEdgioApplicationsConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.HUAWEICLOUD_CDN:
          return <DeployNodeConfigFormHuaweiCloudCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.HUAWEICLOUD_ELB:
          return <DeployNodeConfigFormHuaweiCloudELBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.KUBERNETES_SECRET:
          return <DeployNodeConfigFormKubernetesSecretConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.LOCAL:
          return <DeployNodeConfigFormLocalConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.QINIU_CDN:
          return <DeployNodeConfigFormQiniuCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.QINIU_PILI:
          return <DeployNodeConfigFormQiniuPiliConfig {...nestedFormProps} />;
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
        case DEPLOY_PROVIDERS.UCLOUD_UCDN:
          return <DeployNodeConfigFormUCloudUCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.UCLOUD_US3:
          return <DeployNodeConfigFormUCloudUS3Config {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_CDN:
          return <DeployNodeConfigFormVolcEngineCDNConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_CLB:
          return <DeployNodeConfigFormVolcEngineCLBConfig {...nestedFormProps} />;
        case DEPLOY_PROVIDERS.VOLCENGINE_DCDN:
          return <DeployNodeConfigFormVolcEngineDCDNConfig {...nestedFormProps} />;
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

      // TODO: 暂时不支持切换部署目标，需后端调整，否则之前若存在部署结果输出就不会再部署
      // 切换部署目标时重置表单，避免其他部署目标的配置字段影响当前部署目标
      if (initialValues?.provider === value) {
        formInst.resetFields();
      } else {
        const oldValues = formInst.getFieldsValue();
        const newValues: Record<string, unknown> = {};
        for (const key in oldValues) {
          if (key === "provider" || key === "providerAccessId" || key === "certificate") {
            newValues[key] = oldValues[key];
          } else {
            newValues[key] = undefined;
          }
        }
        (formInst as FormInstance).setFieldsValue(newValues);

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
            fallback={<DeployProviderPicker placeholder={t("workflow_node.deploy.search.provider.placeholder")} onSelect={handleProviderPick} />}
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
                        if (ACCESS_USAGES.ALL === provider?.usage || ACCESS_USAGES.DEPLOY === provider?.usage) {
                          formInst.setFieldValue("providerAccessId", record.id);
                        }
                      }}
                    />
                  </div>
                </div>
              </label>
              <Form.Item name="providerAccessId" rules={[formRule]}>
                <AccessSelect
                  placeholder={t("workflow_node.deploy.form.provider_access.placeholder")}
                  filter={(record) => {
                    if (fieldProvider) {
                      return deployProvidersMap.get(fieldProvider)?.provider === record.provider;
                    }

                    const provider = accessProvidersMap.get(record.provider);
                    return ACCESS_USAGES.ALL === provider?.usage || ACCESS_USAGES.APPLY === provider?.usage;
                  }}
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

        <Show when={!!fieldProvider}>
          <Divider className="my-1">
            <Typography.Text className="text-xs font-normal" type="secondary">
              {t("workflow_node.deploy.form.params_config.label")}
            </Typography.Text>
          </Divider>

          {nestedFormEl}

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
                    checkedChildren={t("workflow_node.deploy.form.skip_on_last_succeeded.enabled.on")}
                    unCheckedChildren={t("workflow_node.deploy.form.skip_on_last_succeeded.enabled.off")}
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
