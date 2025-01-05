import { memo, useCallback, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { PlusOutlined as PlusOutlinedIcon, QuestionCircleOutlined as QuestionCircleOutlinedIcon } from "@ant-design/icons";
import { Button, Divider, Form, type FormInstance, Select, Tooltip, Typography } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { init } from "i18next";
import { z } from "zod";

import Show from "@/components/Show";
import AccessEditModal from "@/components/access/AccessEditModal";
import AccessSelect from "@/components/access/AccessSelect";
import DeployProviderPicker from "@/components/provider/DeployProviderPicker";
import DeployProviderSelect from "@/components/provider/DeployProviderSelect";
import { ACCESS_USAGES, DEPLOY_PROVIDERS, accessProvidersMap, deployProvidersMap } from "@/domain/provider";
import { type WorkflowNode, type WorkflowNodeConfigForDeploy } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import DeployNodeFormAliyunALBFields from "./DeployNodeFormAliyunALBFields";
import DeployNodeFormAliyunCDNFields from "./DeployNodeFormAliyunCDNFields";
import DeployNodeFormAliyunCLBFields from "./DeployNodeFormAliyunCLBFields";
import DeployNodeFormAliyunDCDNFields from "./DeployNodeFormAliyunDCDNFields";
import DeployNodeFormAliyunNLBFields from "./DeployNodeFormAliyunNLBFields";
import DeployNodeFormAliyunOSSFields from "./DeployNodeFormAliyunOSSFields";
import DeployNodeFormBaiduCloudCDNFields from "./DeployNodeFormBaiduCloudCDNFields";
import DeployNodeFormBytePlusCDNFields from "./DeployNodeFormBytePlusCDNFields";
import DeployNodeFormDogeCloudCDNFields from "./DeployNodeFormDogeCloudCDNFields";
import DeployNodeFormHuaweiCloudCDNFields from "./DeployNodeFormHuaweiCloudCDNFields";
import DeployNodeFormHuaweiCloudELBFields from "./DeployNodeFormHuaweiCloudELBFields";
import DeployNodeFormKubernetesSecretFields from "./DeployNodeFormKubernetesSecretFields";
import DeployNodeFormLocalFields from "./DeployNodeFormLocalFields";
import DeployNodeFormQiniuCDNFields from "./DeployNodeFormQiniuCDNFields";
import DeployNodeFormSSHFields from "./DeployNodeFormSSHFields";
import DeployNodeFormTencentCloudCDNFields from "./DeployNodeFormTencentCloudCDNFields";
import DeployNodeFormTencentCloudCLBFields from "./DeployNodeFormTencentCloudCLBFields";
import DeployNodeFormTencentCloudCOSFields from "./DeployNodeFormTencentCloudCOSFields";
import DeployNodeFormTencentCloudECDNFields from "./DeployNodeFormTencentCloudECDNFields";
import DeployNodeFormTencentCloudEOFields from "./DeployNodeFormTencentCloudEOFields";
import DeployNodeFormVolcEngineCDNFields from "./DeployNodeFormVolcEngineCDNFields";
import DeployNodeFormVolcEngineLiveFields from "./DeployNodeFormVolcEngineLiveFields";
import DeployNodeFormWebhookFields from "./DeployNodeFormWebhookFields";

type DeployNodeFormFieldValues = Partial<WorkflowNodeConfigForDeploy>;

export type DeployFormProps = {
  form: FormInstance;
  formName?: string;
  disabled?: boolean;
  workflowNode: WorkflowNode;
  onValuesChange?: (values: DeployNodeFormFieldValues) => void;
};

const initFormModel = (): DeployNodeFormFieldValues => {
  return {};
};

const DeployNodeForm = ({ form, formName, disabled, workflowNode, onValuesChange }: DeployFormProps) => {
  const { t } = useTranslation();

  const { getWorkflowOuptutBeforeId } = useWorkflowStore(useZustandShallowSelector(["updateNode", "getWorkflowOuptutBeforeId"]));

  const [previousOutput, setPreviousOutput] = useState<WorkflowNode[]>([]);
  useEffect(() => {
    const rs = getWorkflowOuptutBeforeId(workflowNode.id, "certificate");
    setPreviousOutput(rs);
  }, [workflowNode.id, getWorkflowOuptutBeforeId]);

  const formSchema = z.object({
    provider: z.string({ message: t("workflow_node.deploy.form.provider.placeholder") }).nonempty(t("workflow_node.deploy.form.provider.placeholder")),
    providerAccessId: z
      .string({ message: t("workflow_node.deploy.form.provider_access.placeholder") })
      .nonempty(t("workflow_node.deploy.form.provider_access.placeholder")),
    certificate: z.string({ message: t("workflow_node.deploy.form.certificate.placeholder") }).nonempty(t("workflow_node.deploy.form.certificate.placeholder")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const initialValues: DeployNodeFormFieldValues = (workflowNode.config as WorkflowNodeConfigForDeploy) ?? initFormModel();

  const fieldProvider = Form.useWatch("provider", { form: form, preserve: true });

  const formFieldsComponent = useMemo(() => {
    /*
      注意：如果追加新的子组件，请保持以 ASCII 排序。
      NOTICE: If you add new child component, please keep ASCII order.
     */
    switch (fieldProvider) {
      case DEPLOY_PROVIDERS.ALIYUN_ALB:
        return <DeployNodeFormAliyunALBFields />;
      case DEPLOY_PROVIDERS.ALIYUN_CLB:
        return <DeployNodeFormAliyunCLBFields />;
      case DEPLOY_PROVIDERS.ALIYUN_CDN:
        return <DeployNodeFormAliyunCDNFields />;
      case DEPLOY_PROVIDERS.ALIYUN_DCDN:
        return <DeployNodeFormAliyunDCDNFields />;
      case DEPLOY_PROVIDERS.ALIYUN_NLB:
        return <DeployNodeFormAliyunNLBFields />;
      case DEPLOY_PROVIDERS.ALIYUN_OSS:
        return <DeployNodeFormAliyunOSSFields />;
      case DEPLOY_PROVIDERS.BAIDUCLOUD_CDN:
        return <DeployNodeFormBaiduCloudCDNFields />;
      case DEPLOY_PROVIDERS.BYTEPLUS_CDN:
        return <DeployNodeFormBytePlusCDNFields />;
      case DEPLOY_PROVIDERS.DOGECLOUD_CDN:
        return <DeployNodeFormDogeCloudCDNFields />;
      case DEPLOY_PROVIDERS.HUAWEICLOUD_CDN:
        return <DeployNodeFormHuaweiCloudCDNFields />;
      case DEPLOY_PROVIDERS.HUAWEICLOUD_ELB:
        return <DeployNodeFormHuaweiCloudELBFields />;
      case DEPLOY_PROVIDERS.KUBERNETES_SECRET:
        return <DeployNodeFormKubernetesSecretFields />;
      case DEPLOY_PROVIDERS.LOCAL:
        return <DeployNodeFormLocalFields />;
      case DEPLOY_PROVIDERS.QINIU_CDN:
        return <DeployNodeFormQiniuCDNFields />;
      case DEPLOY_PROVIDERS.SSH:
        return <DeployNodeFormSSHFields />;
      case DEPLOY_PROVIDERS.TENCENTCLOUD_CDN:
        return <DeployNodeFormTencentCloudCDNFields />;
      case DEPLOY_PROVIDERS.TENCENTCLOUD_CLB:
        return <DeployNodeFormTencentCloudCLBFields />;
      case DEPLOY_PROVIDERS.TENCENTCLOUD_COS:
        return <DeployNodeFormTencentCloudCOSFields />;
      case DEPLOY_PROVIDERS.TENCENTCLOUD_ECDN:
        return <DeployNodeFormTencentCloudECDNFields />;
      case DEPLOY_PROVIDERS.TENCENTCLOUD_EO:
        return <DeployNodeFormTencentCloudEOFields />;
      case DEPLOY_PROVIDERS.VOLCENGINE_CDN:
        return <DeployNodeFormVolcEngineCDNFields />;
      case DEPLOY_PROVIDERS.VOLCENGINE_LIVE:
        return <DeployNodeFormVolcEngineLiveFields />;
      case DEPLOY_PROVIDERS.WEBHOOK:
        return <DeployNodeFormWebhookFields />;
    }
  }, [fieldProvider]);

  const handleProviderPick = useCallback(
    (value: string) => {
      form.setFieldValue("provider", value);
    },
    [form]
  );

  const handleProviderSelect = (value: string) => {
    if (fieldProvider === value) return;

    // TODO: 暂时不支持切换部署目标，需后端调整，否则之前若存在部署结果输出就不会再部署
    // 切换部署目标时重置表单，避免其他部署目标的配置字段影响当前部署目标
    if (initialValues?.provider === value) {
      form.resetFields();
    } else {
      const oldValues = form.getFieldsValue();
      const newValues: Record<string, unknown> = {};
      for (const key in oldValues) {
        if (key === "provider" || key === "providerAccessId" || key === "certificate") {
          newValues[key] = oldValues[key];
        } else {
          newValues[key] = undefined;
        }
      }
      form.setFieldsValue(newValues);

      if (deployProvidersMap.get(fieldProvider)?.provider !== deployProvidersMap.get(value)?.provider) {
        form.setFieldValue("providerAccessId", undefined);
      }
    }
  };

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as DeployNodeFormFieldValues);
  };

  return (
    <Form
      form={form}
      disabled={disabled}
      initialValues={initialValues}
      layout="vertical"
      name={formName}
      preserve={false}
      scrollToFirstError
      onValuesChange={handleFormChange}
    >
      <Show when={!!fieldProvider} fallback={<DeployProviderPicker onSelect={handleProviderPick} />}>
        <Form.Item name="provider" label={t("workflow_node.deploy.form.provider.label")} rules={[formRule]}>
          <DeployProviderSelect
            allowClear
            disabled
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
                  onSubmit={(record) => {
                    const provider = accessProvidersMap.get(record.provider);
                    if (ACCESS_USAGES.ALL === provider?.usage || ACCESS_USAGES.DEPLOY === provider?.usage) {
                      form.setFieldValue("providerAccessId", record.id);
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

        <Form.Item
          name="certificate"
          label={t("workflow_node.deploy.form.certificate.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.certificate.tooltip") }}></span>}
        >
          <Select
            options={previousOutput.map((item) => {
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

        <Divider className="my-1">
          <Typography.Text className="text-xs font-normal" type="secondary">
            {t("workflow_node.deploy.form.params_config.label")}
          </Typography.Text>
        </Divider>

        {formFieldsComponent}
      </Show>
    </Form>
  );
};

export default memo(DeployNodeForm);
