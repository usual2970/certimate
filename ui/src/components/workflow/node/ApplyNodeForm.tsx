import { memo, useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { AutoComplete, Button, Divider, Form, Input, Select, Space, Switch, Tooltip, Typography, type AutoCompleteProps } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { FormOutlined as FormOutlinedIcon, PlusOutlined as PlusOutlinedIcon, QuestionCircleOutlined as QuestionCircleOutlinedIcon } from "@ant-design/icons";
import { produce } from "immer";
import z from "zod";

import AccessEditModal from "@/components/access/AccessEditModal";
import AccessSelect from "@/components/access/AccessSelect";
import ModalForm from "@/components/core/ModalForm";
import MultipleInput from "@/components/core/MultipleInput";
import { usePanel } from "../PanelProvider";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { ACCESS_USAGES } from "@/domain/access";
import { accessProvidersMap } from "@/domain/provider";
import { type WorkflowNode, type WorkflowNodeConfig } from "@/domain/workflow";
import { useContactStore } from "@/stores/contact";
import { useWorkflowStore } from "@/stores/workflow";
import { validDomainName, validIPv4Address, validIPv6Address } from "@/utils/validators";

export type ApplyNodeFormProps = {
  data: WorkflowNode;
};

const MULTIPLE_INPUT_DELIMITER = ";";

const initFormModel = (): WorkflowNodeConfig => {
  return {
    domain: "",
    keyAlgorithm: "RSA2048",
    nameservers: "",
    propagationTimeout: 60,
    disableFollowCNAME: true,
  };
};

const ApplyNodeForm = ({ data }: ApplyNodeFormProps) => {
  const { t } = useTranslation();

  const { addEmail } = useContactStore();
  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));
  const { hidePanel } = usePanel();

  const formSchema = z.object({
    domain: z.string({ message: t("workflow_node.apply.form.domains.placeholder") }).refine((v) => {
      return String(v)
        .split(MULTIPLE_INPUT_DELIMITER)
        .every((e) => validDomainName(e, true));
    }, t("common.errmsg.domain_invalid")),
    email: z.string({ message: t("workflow_node.apply.form.email.placeholder") }).email("common.errmsg.email_invalid"),
    access: z.string({ message: t("workflow_node.apply.form.access.placeholder") }).min(1, t("workflow_node.apply.form.access.placeholder")),
    keyAlgorithm: z.string().nullish(),
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
  const {
    form: formInst,
    formPending,
    formProps,
  } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: data?.config ?? initFormModel(),
    onSubmit: async (values) => {
      await formInst.validateFields();
      await addEmail(values.email);
      await updateNode(
        produce(data, (draft) => {
          draft.config = { ...values };
          draft.validated = true;
        })
      );
      hidePanel();
    },
  });

  const [fieldDomains, setFieldDomains] = useState(data?.config?.domain as string);
  const [fieldNameservers, setFieldNameservers] = useState(data?.config?.nameservers as string);

  const handleFieldDomainsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setFieldDomains(value);
    formInst.setFieldValue("domain", value);
  };

  const handleFieldNameserversChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setFieldNameservers(value);
    formInst.setFieldValue("nameservers", value);
  };

  return (
    <Form {...formProps} form={formInst} disabled={formPending} layout="vertical">
      <Form.Item
        name="domain"
        label={t("workflow_node.apply.form.domains.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.domains.tooltip") }}></span>}
      >
        <Space.Compact style={{ width: "100%" }}>
          <Input
            disabled={formPending}
            value={fieldDomains}
            placeholder={t("workflow_node.apply.form.domains.placeholder")}
            onChange={handleFieldDomainsChange}
          />
          <FormFieldDomainsModalForm
            data={fieldDomains}
            trigger={
              <Button disabled={formPending}>
                <FormOutlinedIcon />
              </Button>
            }
            onFinish={(v) => {
              setFieldDomains(v);
              formInst.setFieldValue("domain", v);
            }}
          />
        </Space.Compact>
      </Form.Item>

      <Form.Item
        name="email"
        label={t("workflow_node.apply.form.email.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.email.tooltip") }}></span>}
      >
        <FormFieldEmailSelect placeholder={t("workflow_node.apply.form.email.placeholder")} />
      </Form.Item>

      <Form.Item className="mb-0">
        <label className="block mb-1">
          <div className="flex items-center justify-between gap-4 w-full">
            <div className="flex-grow max-w-full truncate">
              <span>{t("workflow_node.apply.form.access.label")}</span>
              <Tooltip title={t("workflow_node.apply.form.access.tooltip")}>
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
                    {t("workflow_node.apply.form.access.button")}
                  </Button>
                }
                onSubmit={(record) => {
                  const provider = accessProvidersMap.get(record.configType);
                  if (ACCESS_USAGES.ALL === provider?.usage || ACCESS_USAGES.APPLY === provider?.usage) {
                    formInst.setFieldValue("access", record.id);
                  }
                }}
              />
            </div>
          </div>
        </label>
        <Form.Item name="access" rules={[formRule]}>
          <AccessSelect
            placeholder={t("workflow_node.apply.form.access.placeholder")}
            filter={(record) => {
              const provider = accessProvidersMap.get(record.configType);
              return ACCESS_USAGES.ALL === provider?.usage || ACCESS_USAGES.APPLY === provider?.usage;
            }}
          />
        </Form.Item>
      </Form.Item>

      <Divider className="my-1">
        <Typography.Text className="font-normal text-xs" type="secondary">
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
        name="nameservers"
        label={t("workflow_node.apply.form.nameservers.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.apply.form.nameservers.tooltip") }}></span>}
      >
        <Space.Compact style={{ width: "100%" }}>
          <Input
            allowClear
            disabled={formPending}
            value={fieldNameservers}
            placeholder={t("workflow_node.apply.form.nameservers.placeholder")}
            onChange={handleFieldNameserversChange}
          />
          <FormFieldNameserversModalForm
            data={fieldNameservers}
            trigger={
              <Button disabled={formPending}>
                <FormOutlinedIcon />
              </Button>
            }
            onFinish={(v) => {
              setFieldNameservers(v);
              formInst.setFieldValue("nameservers", v);
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

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={formPending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

const FormFieldEmailSelect = ({
  className,
  style,
  disabled,
  placeholder,
  ...props
}: {
  className?: string;
  style?: React.CSSProperties;
  defaultValue?: string;
  disabled?: boolean;
  placeholder?: string;
  value?: string;
  onChange?: (value: string) => void;
}) => {
  const { emails, fetchEmails } = useContactStore();
  const emailsToOptions = useCallback(() => emails.map((email) => ({ label: email, value: email })), [emails]);
  useEffect(() => {
    fetchEmails();
  }, [fetchEmails]);

  const [value, setValue] = useControllableValue<string>(props, {
    valuePropName: "value",
    defaultValuePropName: "defaultValue",
    trigger: "onChange",
  });

  const [options, setOptions] = useState<AutoCompleteProps["options"]>([]);
  useEffect(() => {
    setOptions(emailsToOptions());
  }, [emails, emailsToOptions]);

  const handleChange = (value: string) => {
    setValue(value);
  };

  const handleSearch = (text: string) => {
    const temp = emailsToOptions();
    if (text) {
      if (temp.every((option) => option.label !== text)) {
        temp.unshift({ label: text, value: text });
      }
    }

    setOptions(temp);
  };

  return (
    <AutoComplete
      className={className}
      style={style}
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
};

const FormFieldDomainsModalForm = ({
  data,
  trigger,
  onFinish,
}: {
  data: string;
  disabled?: boolean;
  trigger?: React.ReactNode;
  onFinish?: (data: string) => void;
}) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    domains: z.array(z.string()).refine((v) => {
      return v.every((e) => !e?.trim() || validDomainName(e.trim(), true));
    }, t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const [formInst] = Form.useForm<z.infer<typeof formSchema>>();

  const [model, setModel] = useState<Partial<z.infer<typeof formSchema>>>({ domains: data?.split(MULTIPLE_INPUT_DELIMITER) });
  useEffect(() => {
    setModel({ domains: data?.split(MULTIPLE_INPUT_DELIMITER) });
  }, [data]);

  const handleFormFinish = (values: z.infer<typeof formSchema>) => {
    onFinish?.(
      values.domains
        .map((e) => e.trim())
        .filter((e) => !!e)
        .join(MULTIPLE_INPUT_DELIMITER)
    );
  };

  return (
    <ModalForm
      layout="vertical"
      form={formInst}
      initialValues={model}
      modalProps={{ destroyOnClose: true }}
      title={t("workflow_node.apply.form.domains.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
      onFinish={handleFormFinish}
    >
      <Form.Item name="domains" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.apply.form.domains.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
};

const FormFieldNameserversModalForm = ({ data, trigger, onFinish }: { data: string; trigger?: React.ReactNode; onFinish?: (data: string) => void }) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    nameservers: z.array(z.string()).refine((v) => {
      return v.every((e) => !e?.trim() || validIPv4Address(e) || validIPv6Address(e) || validDomainName(e));
    }, t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const [formInst] = Form.useForm<z.infer<typeof formSchema>>();

  const [model, setModel] = useState<Partial<z.infer<typeof formSchema>>>({ nameservers: data?.split(MULTIPLE_INPUT_DELIMITER) });
  useEffect(() => {
    setModel({ nameservers: data?.split(MULTIPLE_INPUT_DELIMITER) });
  }, [data]);

  const handleFormFinish = (values: z.infer<typeof formSchema>) => {
    onFinish?.(
      values.nameservers
        .map((e) => e.trim())
        .filter((e) => !!e)
        .join(MULTIPLE_INPUT_DELIMITER)
    );
  };

  return (
    <ModalForm
      layout="vertical"
      form={formInst}
      initialValues={model}
      modalProps={{ destroyOnClose: true }}
      title={t("workflow_node.apply.form.nameservers.multiple_input_modal.title")}
      trigger={trigger}
      validateTrigger="onSubmit"
      width={480}
      onFinish={handleFormFinish}
    >
      <Form.Item name="nameservers" rules={[formRule]}>
        <MultipleInput placeholder={t("workflow_node.apply.form.nameservers.multiple_input_modal.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
};

export default memo(ApplyNodeForm);
