import { forwardRef, memo, useEffect, useImperativeHandle } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import { RightOutlined as RightOutlinedIcon } from "@ant-design/icons";
import { Button, Form, type FormInstance, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { notifyChannelsMap } from "@/domain/settings";
import { type WorkflowNodeConfigForNotify } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useNotifyChannelsStore } from "@/stores/notify";

type NotifyNodeConfigFormFieldValues = Partial<WorkflowNodeConfigForNotify>;

export type NotifyNodeConfigFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: NotifyNodeConfigFormFieldValues;
  onValuesChange?: (values: NotifyNodeConfigFormFieldValues) => void;
};

export type NotifyNodeConfigFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<NotifyNodeConfigFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<NotifyNodeConfigFormFieldValues>["resetFields"];
  validateFields: FormInstance<NotifyNodeConfigFormFieldValues>["validateFields"];
};

const initFormModel = (): NotifyNodeConfigFormFieldValues => {
  return {};
};

const NotifyNodeConfigForm = forwardRef<NotifyNodeConfigFormInstance, NotifyNodeConfigFormProps>(
  ({ className, style, disabled, initialValues, onValuesChange }, ref) => {
    const { t } = useTranslation();

    const {
      channels,
      loadedAtOnce: channelsLoadedAtOnce,
      fetchChannels,
    } = useNotifyChannelsStore(useZustandShallowSelector(["channels", "loadedAtOnce", "fetchChannels"]));
    useEffect(() => {
      fetchChannels();
    }, []);

    const formSchema = z.object({
      subject: z
        .string({ message: t("workflow_node.notify.form.subject.placeholder") })
        .min(1, t("workflow_node.notify.form.subject.placeholder"))
        .max(1000, t("common.errmsg.string_max", { max: 1000 })),
      message: z
        .string({ message: t("workflow_node.notify.form.message.placeholder") })
        .min(1, t("workflow_node.notify.form.message.placeholder"))
        .max(1000, t("common.errmsg.string_max", { max: 1000 })),
      channel: z.string({ message: t("workflow_node.notify.form.channel.placeholder") }).min(1, t("workflow_node.notify.form.channel.placeholder")),
    });
    const formRule = createSchemaFieldRule(formSchema);
    const { form: formInst, formProps } = useAntdForm({
      name: "workflowNodeNotifyConfigForm",
      initialValues: initialValues ?? initFormModel(),
    });

    const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
      onValuesChange?.(values as NotifyNodeConfigFormFieldValues);
    };

    useImperativeHandle(ref, () => {
      return {
        getFieldsValue: () => {
          return formInst.getFieldsValue(true);
        },
        resetFields: (fields) => {
          return formInst.resetFields(fields as (keyof NotifyNodeConfigFormFieldValues)[]);
        },
        validateFields: (nameList, config) => {
          return formInst.validateFields(nameList, config);
        },
      } as NotifyNodeConfigFormInstance;
    });

    return (
      <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
        <Form.Item name="subject" label={t("workflow_node.notify.form.subject.label")} rules={[formRule]}>
          <Input placeholder={t("workflow_node.notify.form.subject.placeholder")} />
        </Form.Item>

        <Form.Item name="message" label={t("workflow_node.notify.form.message.label")} rules={[formRule]}>
          <Input.TextArea autoSize={{ minRows: 3, maxRows: 5 }} placeholder={t("workflow_node.notify.form.message.placeholder")} />
        </Form.Item>

        <Form.Item className="mb-0">
          <label className="mb-1 block">
            <div className="flex w-full items-center justify-between gap-4">
              <div className="max-w-full grow truncate">{t("workflow_node.notify.form.channel.label")}</div>
              <div className="text-right">
                <Link className="ant-typography" to="/settings/notification" target="_blank">
                  <Button size="small" type="link">
                    {t("workflow_node.notify.form.channel.button")}
                    <RightOutlinedIcon className="text-xs" />
                  </Button>
                </Link>
              </div>
            </div>
          </label>
          <Form.Item name="channel" rules={[formRule]}>
            <Select
              loading={!channelsLoadedAtOnce}
              options={Object.entries(channels)
                .filter(([_, v]) => v?.enabled)
                .map(([k, _]) => ({
                  label: t(notifyChannelsMap.get(k)?.name ?? k),
                  value: k,
                }))}
              placeholder={t("workflow_node.notify.form.channel.placeholder")}
            />
          </Form.Item>
        </Form.Item>
      </Form>
    );
  }
);

export default memo(NotifyNodeConfigForm);
