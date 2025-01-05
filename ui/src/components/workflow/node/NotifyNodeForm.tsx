import { memo, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import { RightOutlined as RightOutlinedIcon } from "@ant-design/icons";
import { Button, Form, type FormInstance, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { notifyChannelsMap } from "@/domain/settings";
import { type WorkflowNode, type WorkflowNodeConfigForNotify } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useNotifyChannelsStore } from "@/stores/notify";

type NotifyNodeFormFieldValues = Partial<WorkflowNodeConfigForNotify>;

export type NotifyNodeFormProps = {
  form: FormInstance;
  formName?: string;
  disabled?: boolean;
  workflowNode: WorkflowNode;
  onValuesChange?: (values: NotifyNodeFormFieldValues) => void;
};

const initFormModel = (): NotifyNodeFormFieldValues => {
  return {
    subject: "Completed!",
    message: "Your workflow has been completed on Certimate.",
  };
};

const NotifyNodeForm = ({ form, formName, disabled, workflowNode, onValuesChange }: NotifyNodeFormProps) => {
  const { t } = useTranslation();

  const {
    channels,
    loadedAtOnce: channelsLoadedAtOnce,
    fetchChannels,
  } = useNotifyChannelsStore(useZustandShallowSelector(["channels", "loadedAtOnce", "fetchChannels"]));
  useEffect(() => {
    fetchChannels();
  }, [fetchChannels]);

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

  const initialValues: NotifyNodeFormFieldValues = (workflowNode.config as WorkflowNodeConfigForNotify) ?? initFormModel();

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as NotifyNodeFormFieldValues);
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
      <Form.Item name="subject" label={t("workflow_node.notify.form.subject.label")} rules={[formRule]}>
        <Input placeholder={t("workflow_node.notify.form.subject.placeholder")} />
      </Form.Item>

      <Form.Item name="message" label={t("workflow_node.notify.form.message.label")} rules={[formRule]}>
        <Input.TextArea autoSize={{ minRows: 3, maxRows: 10 }} placeholder={t("workflow_node.notify.form.message.placeholder")} />
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
};

export default memo(NotifyNodeForm);
