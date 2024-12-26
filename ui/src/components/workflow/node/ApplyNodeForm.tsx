import { memo, useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { AutoComplete, Button, Divider, Form, Input, InputNumber, Select, Switch, Typography, type AutoCompleteProps } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { PlusOutlined as PlusOutlinedIcon } from "@ant-design/icons";
import z from "zod";

import AccessEditModal from "@/components/access/AccessEditModal";
import AccessSelect from "@/components/access/AccessSelect";
import { usePanel } from "../PanelProvider";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { ACCESS_PROVIDER_USAGES, accessProvidersMap } from "@/domain/access";
import { type WorkflowNode, type WorkflowNodeConfig } from "@/domain/workflow";
import { useContactStore } from "@/stores/contact";
import { useWorkflowStore } from "@/stores/workflow";
import { validDomainName, validIPv4Address, validIPv6Address } from "@/utils/validators";

export type ApplyNodeFormProps = {
  data: WorkflowNode;
};

const initFormModel = (): WorkflowNodeConfig => {
  return {
    domain: "",
    keyAlgorithm: "RSA2048",
    timeout: 60,
    disableFollowCNAME: true,
  };
};

const ApplyNodeForm = ({ data }: ApplyNodeFormProps) => {
  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));
  const { hidePanel } = usePanel();

  const formSchema = z.object({
    domain: z.string({ message: t("workflow.nodes.apply.form.domain.placeholder") }).refine(
      (str) => {
        return String(str)
          .split(";")
          .every((e) => validDomainName(e, true));
      },
      { message: t("common.errmsg.domain_invalid") }
    ),
    email: z.string({ message: t("workflow.nodes.apply.form.email.placeholder") }).email("common.errmsg.email_invalid"),
    access: z.string({ message: t("workflow.nodes.apply.form.access.placeholder") }).min(1, t("workflow.nodes.apply.form.access.placeholder")),
    keyAlgorithm: z.string().nullish(),
    nameservers: z
      .string()
      .refine(
        (str) => {
          if (!str) return true;
          return String(str)
            .split(";")
            .every((e) => validDomainName(e) || validIPv4Address(e) || validIPv6Address(e));
        },
        { message: t("common.errmsg.host_invalid") }
      )
      .nullish(),
    timeout: z.number().gte(1, t("workflow.nodes.apply.form.timeout.placeholder")).nullish(),
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
      await updateNode({ ...data, config: { ...values }, validated: true });
      hidePanel();
    },
  });

  return (
    <Form {...formProps} form={formInst} disabled={formPending} layout="vertical">
      <Form.Item
        name="domain"
        label={t("workflow.nodes.apply.form.domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.apply.form.domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow.nodes.apply.form.domain.placeholder")} />
      </Form.Item>

      <Form.Item
        name="email"
        label={t("workflow.nodes.apply.form.email.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.apply.form.email.tooltip") }}></span>}
      >
        <ContactEmailSelect placeholder={t("workflow.nodes.apply.form.email.placeholder")} />
      </Form.Item>

      <Form.Item>
        <label className="block mb-[2px]">
          <div className="flex items-center justify-between gap-4 w-full overflow-hidden">
            <div className="flex-grow max-w-full truncate">{t("workflow.nodes.apply.form.access.label")}</div>
            <div className="text-right">
              <AccessEditModal
                preset="add"
                trigger={
                  <Button className="p-0" type="link">
                    <PlusOutlinedIcon />
                    {t("workflow.nodes.apply.form.access.button")}
                  </Button>
                }
              />
            </div>
          </div>
        </label>
        <Form.Item name="access" rules={[formRule]}>
          <AccessSelect
            placeholder={t("workflow.nodes.apply.form.access.placeholder")}
            filter={(record) => {
              const provider = accessProvidersMap.get(record.configType);
              return ACCESS_PROVIDER_USAGES.ALL === provider?.usage || ACCESS_PROVIDER_USAGES.APPLY === provider?.usage;
            }}
          />
        </Form.Item>
      </Form.Item>

      <Divider className="my-1">
        <Typography.Text className="text-xs" type="secondary">
          {t("workflow.nodes.apply.form.advanced_settings.label")}
        </Typography.Text>
      </Divider>

      <Form.Item name="keyAlgorithm" label={t("workflow.nodes.apply.form.key_algorithm.label")} rules={[formRule]}>
        <Select
          options={["RSA2048", "RSA3072", "RSA4096", "RSA8192", "EC256", "EC384"].map((e) => ({
            label: e,
            value: e,
          }))}
          placeholder={t("workflow.nodes.apply.form.key_algorithm.placeholder")}
        />
      </Form.Item>

      <Form.Item
        name="nameservers"
        label={t("workflow.nodes.apply.form.nameservers.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.apply.form.nameservers.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow.nodes.apply.form.nameservers.placeholder")} />
      </Form.Item>

      <Form.Item
        name="timeout"
        label={t("workflow.nodes.apply.form.timeout.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.apply.form.timeout.tooltip") }}></span>}
      >
        <InputNumber
          className="w-full"
          min={0}
          max={3600}
          placeholder={t("workflow.nodes.apply.form.timeout.placeholder")}
          addonAfter={t("workflow.nodes.apply.form.timeout.suffix")}
        />
      </Form.Item>

      <Form.Item
        name="disableFollowCNAME"
        label={t("workflow.nodes.apply.form.disable_follow_cname.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.apply.form.disable_follow_cname.tooltip") }}></span>}
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

const ContactEmailSelect = ({
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

export default memo(ApplyNodeForm);
