import { memo, useEffect } from "react";
import { Link } from "react-router";
import { useTranslation } from "react-i18next";
import { Button, Form, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { RightOutlined as RightOutlinedIcon } from "@ant-design/icons";
import { z } from "zod";

import { usePanel } from "../PanelProvider";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { notifyChannelsMap } from "@/domain/settings";
import { type WorkflowNode, type WorkflowNodeConfig } from "@/domain/workflow";
import { useNotifyChannelStore } from "@/stores/notify";
import { useWorkflowStore } from "@/stores/workflow";

export type NotifyNodeFormProps = {
  data: WorkflowNode;
};

const initFormModel = (): WorkflowNodeConfig => {
  return {
    subject: "Completed!",
    message: "Your workflow has been completed on Certimate.",
  };
};

const NotifyNodeForm = ({ data }: NotifyNodeFormProps) => {
  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));
  const { hidePanel } = usePanel();

  const { channels, loadedAtOnce: channelsLoadedAtOnce, fetchChannels } = useNotifyChannelStore();
  useEffect(() => {
    fetchChannels();
  }, [fetchChannels]);

  const formSchema = z.object({
    subject: z
      .string({ message: t("workflow.nodes.notify.form.subject.placeholder") })
      .min(1, t("workflow.nodes.notify.form.subject.placeholder"))
      .max(1000, t("common.errmsg.string_max", { max: 1000 })),
    message: z
      .string({ message: t("workflow.nodes.notify.form.message.placeholder") })
      .min(1, t("workflow.nodes.notify.form.message.placeholder"))
      .max(1000, t("common.errmsg.string_max", { max: 1000 })),
    channel: z.string({ message: t("workflow.nodes.notify.form.channel.placeholder") }).min(1, t("workflow.nodes.notify.form.channel.placeholder")),
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
      <Form.Item name="subject" label={t("workflow.nodes.notify.form.subject.label")} rules={[formRule]}>
        <Input placeholder={t("workflow.nodes.notify.form.subject.placeholder")} />
      </Form.Item>

      <Form.Item name="message" label={t("workflow.nodes.notify.form.message.label")} rules={[formRule]}>
        <Input.TextArea autoSize={{ minRows: 3, maxRows: 10 }} placeholder={t("workflow.nodes.notify.form.message.placeholder")} />
      </Form.Item>

      <Form.Item>
        <label className="block mb-1">
          <div className="flex items-center justify-between gap-4 w-full">
            <div className="flex-grow max-w-full truncate">{t("workflow.nodes.notify.form.channel.label")}</div>
            <div className="text-right">
              <Link className="ant-typography" to="/settings/notification" target="_blank">
                <Button size="small" type="link">
                  {t("workflow.nodes.notify.form.channel.button")}
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
            placeholder={t("workflow.nodes.notify.form.channel.placeholder")}
          />
        </Form.Item>
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={formPending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

export default memo(NotifyNodeForm);
