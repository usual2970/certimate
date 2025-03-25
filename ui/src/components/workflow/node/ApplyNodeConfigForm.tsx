import { forwardRef, memo, useEffect, useImperativeHandle, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { FormOutlined as FormOutlinedIcon, PlusOutlined as PlusOutlinedIcon, QuestionCircleOutlined as QuestionCircleOutlinedIcon } from "@ant-design/icons";
import { useControllableValue } from "ahooks";
import {
  AutoComplete,
  type AutoCompleteProps,
  Button,
  Divider,
  Flex,
  Form,
  type FormInstance,
  Input,
  InputNumber,
  Select,
  Space,
  Switch,
  Tooltip,
  Typography,
} from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import AccessEditModal from "@/components/access/AccessEditModal";
import AccessSelect from "@/components/access/AccessSelect";
import ModalForm from "@/components/ModalForm";
import MultipleInput from "@/components/MultipleInput";
import ApplyDNSProviderSelect from "@/components/provider/ApplyDNSProviderSelect";
import { ACCESS_USAGES, APPLY_DNS_PROVIDERS, accessProvidersMap, applyDNSProvidersMap } from "@/domain/provider";
import { type WorkflowNodeConfigForApply } from "@/domain/workflow";
import { useAntdForm, useAntdFormName, useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { useContactEmailsStore } from "@/stores/contact";
import { validDomainName, validIPv4Address, validIPv6Address } from "@/utils/validators";

import ApplyNodeConfigFormAWSRoute53Config from "./ApplyNodeConfigFormAWSRoute53Config";
import ApplyNodeConfigFormHuaweiCloudDNSConfig from "./ApplyNodeConfigFormHuaweiCloudDNSConfig";
import ApplyNodeConfigFormJDCloudDNSConfig from "./ApplyNodeConfigFormJDCloudDNSConfig";
import ApplyNodeConfigFormTencentCloudEOConfig from "./ApplyNodeConfigFormTencentCloudEOConfig";

type ApplyNodeConfigFormFieldValues = Partial<WorkflowNodeConfigForApply>;

export type ApplyNodeConfigFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: ApplyNodeConfigFormFieldValues;
  onValuesChange?: (values: ApplyNodeConfigFormFieldValues) => void;
};

export type ApplyNodeConfigFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<ApplyNodeConfigFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<ApplyNodeConfigFormFieldValues>["resetFields"];
  validateFields: FormInstance<ApplyNodeConfigFormFieldValues>["validateFields"];
};

const MULTIPLE_INPUT_DELIMITER = ";";

const initFormModel = (): ApplyNodeConfigFormFieldValues => {
  return {
    challengeType: "dns-01",
    keyAlgorithm: "RSA2048",
    skipBeforeExpiryDays: 20,
  };
};

const ApplyNodeConfigForm = forwardRef<ApplyNodeConfigFormInstance, ApplyNodeConfigFormProps>(
  ({ className, style, disabled, initialValues, onValuesChange }, ref) => {
    const { t } = useTranslation();

    const { accesses } = useAccessesStore(useZustandShallowSelector("accesses"));

    const formSchema = z.object({
      domains: z.string({ message: t("workflow_node.apply.form.domains.placeholder") }).refine((v) => {
        if (!v) return false;
        return String(v)
          .split(MULTIPLE_INPUT_DELIMITER)
          .every((e) => validDomainName(e, { allowWildcard: true }));
      }, t("common.errmsg.domain_invalid")),
      contactEmail: z.string({ message: t("workflow_node.apply.form.contact_email.placeholder") }).email(t("common.errmsg.email_invalid")),
      challengeType: z.string().nullish(),
      provider: z.string({ message: t("workflow_node.apply.form.provider.placeholder") }).nonempty(t("workflow_node.apply.form.provider.placeholder")),
      providerAccessId: z
        .string({ message: t("workflow_node.apply.form.provider_access.placeholder") })
        .min(1, t("workflow_node.apply.form.provider_access.placeholder")),
      providerConfig: z.any(),
      keyAlgorithm: z
        .string({ message: t("workflow_node.apply.form.key_algorithm.placeholder") })
        .nonempty(t("workflow_node.apply.form.key_algorithm.placeholder")),
      nameservers: z
        .string()
        .nullish()
        .refine((v) => {
          if (!v) return true;
          return String(v)
            .split(MULTIPLE_INPUT_DELIMITER)
            .every((e) => validIPv4Address(e) || validIPv6Address(e) || validDomainName(e));
        }, t("common.errmsg.host_invalid")),
      dnsPropagationTimeout: z
        .union([
          z.number().int().gte(1, t("workflow_node.apply.form.dns_propagation_timeout.placeholder")),
          z.string().refine((v) => !v || /^[1-9]\d*$/.test(v), t("workflow_node.apply.form.dns_propagation_timeout.placeholder")),
        ])
        .nullish(),
      dnsTTL: z
        .union([
          z.number().int().gte(1, t("workflow_node.apply.form.dns_ttl.placeholder")),
          z.string().refine((v) => !v || /^[1-9]\d*$/.test(v), t("workflow_node.apply.form.dns_ttl.placeholder")),
        ])
        .nullish(),
      disableFollowCNAME: z.boolean().nullish(),
      disableARI: z.boolean().nullish(),
      skipBeforeExpiryDays: z
        .number({ message: t("workflow_node.apply.form.skip_before_expiry_days.placeholder") })
        .int(t("workflow_node.apply.form.skip_before_expiry_days.placeholder"))
        .gte(1, t("workflow_node.apply.form.skip_before_expiry_days.placeholder")),
    });
    const formRule = createSchemaFieldRule(formSchema);
    const { form: formInst, formProps } = useAntdForm({
      name: "workflowNodeApplyConfigForm",
      initialValues: initialValues ?? initFormModel(),
    });

    const fieldProvider = Form.useWatch<string>("provider", { form: formInst, preserve: true });
    const fieldProviderAccessId = Form.useWatch<string>("providerAccessId", formInst);
    const fieldDomains = Form.useWatch<string>("domains", formInst);
    const fieldNameservers = Form.useWatch<string>("nameservers", formInst);

    const [showProvider, setShowProvider] = useState(false);
    useEffect(() => {
      // 通常情况下每个授权信息只对应一个 DNS 提供商，此时无需显示 DNS 提供商字段；
      // 如果对应多个（如 AWS 的 Route53、Lightsail，腾讯云的 DNS、EdgeOne 等），则显示。
      if (fieldProviderAccessId) {
        const access = accesses.find((e) => e.id === fieldProviderAccessId);
        const providers = Array.from(applyDNSProvidersMap.values()).filter((e) => e.provider === access?.provider);
        setShowProvider(providers.length > 1);
      } else {
        setShowProvider(false);
      }
    }, [accesses, fieldProviderAccessId]);

    const [nestedFormInst] = Form.useForm();
    const nestedFormName = useAntdFormName({ form: nestedFormInst, name: "workflowNodeApplyConfigFormProviderConfigForm" });
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
        case APPLY_DNS_PROVIDERS.AWS:
        case APPLY_DNS_PROVIDERS.AWS_ROUTE53:
          return <ApplyNodeConfigFormAWSRoute53Config {...nestedFormProps} />;
        case APPLY_DNS_PROVIDERS.HUAWEICLOUD:
        case APPLY_DNS_PROVIDERS.HUAWEICLOUD_DNS:
          return <ApplyNodeConfigFormHuaweiCloudDNSConfig {...nestedFormProps} />;
        case APPLY_DNS_PROVIDERS.JDCLOUD:
        case APPLY_DNS_PROVIDERS.JDCLOUD_DNS:
          return <ApplyNodeConfigFormJDCloudDNSConfig {...nestedFormProps} />;
        case APPLY_DNS_PROVIDERS.TENCENTCLOUD_EO:
          return <ApplyNodeConfigFormTencentCloudEOConfig {...nestedFormProps} />;
      }
    }, [disabled, initialValues?.providerConfig, fieldProvider, nestedFormInst, nestedFormName]);

    const handleProviderSelect = (value: string) => {
      if (fieldProvider === value) return;

      // 切换 DNS 提供商时联动授权信息
      if (initialValues?.provider === value) {
        formInst.setFieldValue("providerAccessId", initialValues?.providerAccessId);
        onValuesChange?.(formInst.getFieldsValue(true));
      } else {
        if (applyDNSProvidersMap.get(fieldProvider)?.provider !== applyDNSProvidersMap.get(value)?.provider) {
          formInst.setFieldValue("providerAccessId", undefined);
          onValuesChange?.(formInst.getFieldsValue(true));
        }
      }
    };

    const handleProviderAccessSelect = (value: string) => {
      if (fieldProviderAccessId === value) return;

      // 切换授权信息时联动 DNS 提供商
      const access = accesses.find((access) => access.id === value);
      const provider = Array.from(applyDNSProvidersMap.values()).find((provider) => provider.provider === access?.provider);
      if (fieldProvider !== provider?.type) {
        formInst.setFieldValue("provider", provider?.type);
        onValuesChange?.(formInst.getFieldsValue(true));
      }
    };

    const handleFormProviderChange = (name: string) => {
      if (name === nestedFormName) {
        formInst.setFieldValue("providerConfig", nestedFormInst.getFieldsValue());
        onValuesChange?.(formInst.getFieldsValue(true));
      }
    };

    const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
      onValuesChange?.(values as ApplyNodeConfigFormFieldValues);
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
      } as ApplyNodeConfigFormInstance;
    });

    return (
      <Form.Provider onFormChange={handleFormProviderChange}>
        <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
          <Form.Item
            label={t("workflow_node.apply.form.domains.label")}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.domains.tooltip") }}></span>}
          >
            <Space.Compact style={{ width: "100%" }}>
              <Form.Item name="domains" noStyle rules={[formRule]}>
                <Input placeholder={t("workflow_node.apply.form.domains.placeholder")} />
              </Form.Item>
              <DomainsModalInput
                value={fieldDomains}
                trigger={
                  <Button disabled={disabled}>
                    <FormOutlinedIcon />
                  </Button>
                }
                onChange={(v) => {
                  formInst.setFieldValue("domains", v);
                }}
              />
            </Space.Compact>
          </Form.Item>

          <Form.Item
            name="contactEmail"
            label={t("workflow_node.apply.form.contact_email.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.contact_email.tooltip") }}></span>}
          >
            <EmailInput placeholder={t("workflow_node.apply.form.contact_email.placeholder")} />
          </Form.Item>

          <Form.Item name="challengeType" label={t("workflow_node.apply.form.challenge_type.label")} rules={[formRule]} hidden>
            <Select
              options={["DNS-01"].map((e) => ({
                label: e,
                value: e.toLowerCase(),
              }))}
              placeholder={t("workflow_node.apply.form.challenge_type.placeholder")}
            />
          </Form.Item>

          <Form.Item name="provider" label={t("workflow_node.apply.form.provider.label")} hidden={!showProvider} rules={[formRule]}>
            <ApplyDNSProviderSelect
              disabled={!showProvider}
              filter={(record) => {
                if (fieldProviderAccessId) {
                  return accesses.find((e) => e.id === fieldProviderAccessId)?.provider === record.provider;
                }

                return true;
              }}
              placeholder={t("workflow_node.apply.form.provider.placeholder")}
              showSearch
              onSelect={handleProviderSelect}
            />
          </Form.Item>

          <Form.Item className="mb-0">
            <label className="mb-1 block">
              <div className="flex w-full items-center justify-between gap-4">
                <div className="max-w-full grow truncate">
                  <span>{t("workflow_node.apply.form.provider_access.label")}</span>
                  <Tooltip title={t("workflow_node.apply.form.provider_access.tooltip")}>
                    <Typography.Text className="ms-1" type="secondary">
                      <QuestionCircleOutlinedIcon />
                    </Typography.Text>
                  </Tooltip>
                </div>
                <div className="text-right">
                  <AccessEditModal
                    preset="add"
                    trigger={
                      <Button size="small" type="link">
                        <PlusOutlinedIcon />
                        {t("workflow_node.apply.form.provider_access.button")}
                      </Button>
                    }
                    afterSubmit={(record) => {
                      const provider = accessProvidersMap.get(record.provider);
                      if (provider?.usages?.includes(ACCESS_USAGES.APPLY)) {
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
                  const provider = accessProvidersMap.get(record.provider);
                  return !!provider?.usages?.includes(ACCESS_USAGES.APPLY);
                }}
                placeholder={t("workflow_node.apply.form.provider_access.placeholder")}
                onChange={handleProviderAccessSelect}
              />
            </Form.Item>
          </Form.Item>
        </Form>

        {nestedFormEl}

        <Divider className="my-1">
          <Typography.Text className="text-xs font-normal" type="secondary">
            {t("workflow_node.apply.form.advanced_config.label")}
          </Typography.Text>
        </Divider>

        <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
          <Form.Item name="keyAlgorithm" label={t("workflow_node.apply.form.key_algorithm.label")} rules={[formRule]}>
            <Select
              options={["RSA2048", "RSA3072", "RSA4096", "RSA8192", "EC256", "EC384"].map((e) => ({
                label: e,
                value: e,
              }))}
              placeholder={t("workflow_node.apply.form.key_algorithm.placeholder")}
            />
          </Form.Item>

          <Form.Item
            label={t("workflow_node.apply.form.nameservers.label")}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.nameservers.tooltip") }}></span>}
          >
            <Space.Compact style={{ width: "100%" }}>
              <Form.Item name="nameservers" noStyle rules={[formRule]}>
                <Input
                  allowClear
                  disabled={disabled}
                  value={fieldNameservers}
                  placeholder={t("workflow_node.apply.form.nameservers.placeholder")}
                  onChange={(e) => {
                    formInst.setFieldValue("nameservers", e.target.value);
                  }}
                />
              </Form.Item>
              <NameserversModalInput
                value={fieldNameservers}
                trigger={
                  <Button disabled={disabled}>
                    <FormOutlinedIcon />
                  </Button>
                }
                onChange={(value) => {
                  formInst.setFieldValue("nameservers", value);
                }}
              />
            </Space.Compact>
          </Form.Item>

          <Form.Item
            name="dnsPropagationTimeout"
            label={t("workflow_node.apply.form.dns_propagation_timeout.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.dns_propagation_timeout.tooltip") }}></span>}
          >
            <Input
              type="number"
              allowClear
              min={0}
              max={3600}
              placeholder={t("workflow_node.apply.form.dns_propagation_timeout.placeholder")}
              addonAfter={t("workflow_node.apply.form.dns_propagation_timeout.unit")}
            />
          </Form.Item>

          <Form.Item
            name="dnsTTL"
            label={t("workflow_node.apply.form.dns_ttl.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.dns_ttl.tooltip") }}></span>}
          >
            <Input
              type="number"
              allowClear
              min={0}
              max={86400}
              placeholder={t("workflow_node.apply.form.dns_ttl.placeholder")}
              addonAfter={t("workflow_node.apply.form.dns_ttl.unit")}
            />
          </Form.Item>

          <Form.Item
            name="disableFollowCNAME"
            label={t("workflow_node.apply.form.disable_follow_cname.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.disable_follow_cname.tooltip") }}></span>}
          >
            <Switch />
          </Form.Item>

          <Form.Item
            name="disableARI"
            label={t("workflow_node.apply.form.disable_ari.label")}
            rules={[formRule]}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.disable_ari.tooltip") }}></span>}
          >
            <Switch />
          </Form.Item>
        </Form>

        <Divider className="my-1">
          <Typography.Text className="text-xs font-normal" type="secondary">
            {t("workflow_node.apply.form.strategy_config.label")}
          </Typography.Text>
        </Divider>

        <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
          <Form.Item
            label={t("workflow_node.apply.form.skip_before_expiry_days.label")}
            tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.skip_before_expiry_days.tooltip") }}></span>}
          >
            <Flex align="center" gap={8} wrap="wrap">
              <div>{t("workflow_node.apply.form.skip_before_expiry_days.prefix")}</div>
              <Form.Item name="skipBeforeExpiryDays" noStyle rules={[formRule]}>
                <InputNumber
                  className="w-36"
                  min={1}
                  max={90}
                  placeholder={t("workflow_node.apply.form.skip_before_expiry_days.placeholder")}
                  addonAfter={t("workflow_node.apply.form.skip_before_expiry_days.unit")}
                />
              </Form.Item>
              <div>{t("workflow_node.apply.form.skip_before_expiry_days.suffix")}</div>
            </Flex>
          </Form.Item>
        </Form>
      </Form.Provider>
    );
  }
);

const EmailInput = memo(
  ({ disabled, placeholder, ...props }: { disabled?: boolean; placeholder?: string; value?: string; onChange?: (value: string) => void }) => {
    const { emails, fetchEmails } = useContactEmailsStore();
    const emailsToOptions = () => emails.map((email) => ({ label: email, value: email }));
    useEffect(() => {
      fetchEmails();
    }, []);

    const [value, setValue] = useControllableValue<string>(props, {
      valuePropName: "value",
      defaultValuePropName: "defaultValue",
      trigger: "onChange",
    });

    const [options, setOptions] = useState<AutoCompleteProps["options"]>([]);
    useEffect(() => {
      setOptions(emailsToOptions());
    }, [emails]);

    const handleChange = (value: string) => {
      setValue(value);
    };

    const handleSearch = (text: string) => {
      const temp = emailsToOptions();
      if (text?.trim()) {
        if (temp.every((option) => option.label !== text)) {
          temp.unshift({ label: text, value: text });
        }
      }

      setOptions(temp);
    };

    return (
      <AutoComplete
        backfill
        defaultValue={value}
        disabled={disabled}
        filterOption
        options={options}
        placeholder={placeholder}
        showSearch
        value={value}
        onChange={handleChange}
        onSearch={handleSearch}
      />
    );
  }
);

const DomainsModalInput = memo(({ value, trigger, onChange }: { value?: string; trigger?: React.ReactNode; onChange?: (value: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    domains: z.array(z.string()).refine((v) => {
      return v.every((e) => !e?.trim() || validDomainName(e.trim(), { allowWildcard: true }));
    }, t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    name: "workflowNodeApplyConfigFormDomainsModalInput",
    initialValues: { domains: value?.split(MULTIPLE_INPUT_DELIMITER) },
    onSubmit: (values) => {
      onChange?.(
        values.domains
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
      title={t("workflow_node.apply.form.domains.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
    >
      <Form.Item name="domains" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.apply.form.domains.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
});

const NameserversModalInput = memo(({ trigger, value, onChange }: { trigger?: React.ReactNode; value?: string; onChange?: (value: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    nameservers: z.array(z.string()).refine((v) => {
      return v.every((e) => !e?.trim() || validIPv4Address(e) || validIPv6Address(e) || validDomainName(e));
    }, t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    name: "workflowNodeApplyConfigFormNameserversModalInput",
    initialValues: { nameservers: value?.split(MULTIPLE_INPUT_DELIMITER) },
    onSubmit: (values) => {
      onChange?.(
        values.nameservers
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
      title={t("workflow_node.apply.form.nameservers.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
    >
      <Form.Item name="nameservers" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.apply.form.nameservers.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
});

export default memo(ApplyNodeConfigForm);
