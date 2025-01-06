import { forwardRef, memo, useEffect, useImperativeHandle, useState } from "react";
import { useTranslation } from "react-i18next";
import { FormOutlined as FormOutlinedIcon, PlusOutlined as PlusOutlinedIcon, QuestionCircleOutlined as QuestionCircleOutlinedIcon } from "@ant-design/icons";
import { useControllableValue } from "ahooks";
import { AutoComplete, type AutoCompleteProps, Button, Divider, Form, type FormInstance, Input, Select, Space, Switch, Tooltip, Typography } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import ModalForm from "@/components/ModalForm";
import MultipleInput from "@/components/MultipleInput";
import AccessEditModal from "@/components/access/AccessEditModal";
import AccessSelect from "@/components/access/AccessSelect";
import { ACCESS_USAGES, accessProvidersMap } from "@/domain/provider";
import { type WorkflowNodeConfigForApply } from "@/domain/workflow";
import { useAntdForm } from "@/hooks";
import { useContactEmailsStore } from "@/stores/contact";
import { validDomainName, validIPv4Address, validIPv6Address } from "@/utils/validators";

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
    keyAlgorithm: "RSA2048",
    propagationTimeout: 60,
    disableFollowCNAME: true,
  };
};

const ApplyNodeConfigForm = forwardRef<ApplyNodeConfigFormInstance, ApplyNodeConfigFormProps>(
  ({ className, style, disabled, initialValues, onValuesChange }, ref) => {
    const { t } = useTranslation();

    const formSchema = z.object({
      domains: z.string({ message: t("workflow_node.apply.form.domains.placeholder") }).refine((v) => {
        return String(v)
          .split(MULTIPLE_INPUT_DELIMITER)
          .every((e) => validDomainName(e, true));
      }, t("common.errmsg.domain_invalid")),
      contactEmail: z.string({ message: t("workflow_node.apply.form.contact_email.placeholder") }).email(t("common.errmsg.email_invalid")),
      providerAccessId: z
        .string({ message: t("workflow_node.apply.form.provider_access.placeholder") })
        .min(1, t("workflow_node.apply.form.provider_access.placeholder")),
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
      propagationTimeout: z
        .union([
          z.number().int().gte(1, t("workflow_node.apply.form.propagation_timeout.placeholder")),
          z.string().refine((v) => !v || /^[1-9]\d*$/.test(v), t("workflow_node.apply.form.propagation_timeout.placeholder")),
        ])
        .nullish(),
      disableFollowCNAME: z.boolean().nullish(),
    });
    const formRule = createSchemaFieldRule(formSchema);
    const { form: formInst, formProps } = useAntdForm({
      name: "workflowNodeApplyConfigForm",
      initialValues: initialValues ?? initFormModel(),
    });

    const fieldDomains = Form.useWatch<string>("domains", formInst);
    const fieldNameservers = Form.useWatch<string>("nameservers", formInst);

    const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
      onValuesChange?.(values as ApplyNodeConfigFormFieldValues);
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
          return formInst.validateFields(nameList, config);
        },
      } as ApplyNodeConfigFormInstance;
    });

    return (
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
                    if (ACCESS_USAGES.ALL === provider?.usage || ACCESS_USAGES.APPLY === provider?.usage) {
                      formInst.setFieldValue("providerAccessId", record.id);
                    }
                  }}
                />
              </div>
            </div>
          </label>
          <Form.Item name="providerAccessId" rules={[formRule]}>
            <AccessSelect
              placeholder={t("workflow_node.apply.form.provider_access.placeholder")}
              filter={(record) => {
                const provider = accessProvidersMap.get(record.provider);
                return ACCESS_USAGES.ALL === provider?.usage || ACCESS_USAGES.APPLY === provider?.usage;
              }}
            />
          </Form.Item>
        </Form.Item>

        <Divider className="my-1">
          <Typography.Text className="text-xs font-normal" type="secondary">
            {t("workflow_node.apply.form.advanced_config.label")}
          </Typography.Text>
        </Divider>

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
          name="propagationTimeout"
          label={t("workflow_node.apply.form.propagation_timeout.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.propagation_timeout.tooltip") }}></span>}
        >
          <Input
            type="number"
            allowClear
            min={0}
            max={3600}
            placeholder={t("workflow_node.apply.form.propagation_timeout.placeholder")}
            addonAfter={t("workflow_node.apply.form.propagation_timeout.suffix")}
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
      </Form>
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
      return v.every((e) => !e?.trim() || validDomainName(e.trim(), true));
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
