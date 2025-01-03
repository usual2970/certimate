import { memo, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import { RightOutlined as RightOutlinedIcon } from "@ant-design/icons";
import { Button, Form, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { produce } from "immer";
import { z } from "zod";

import { notifyChannelsMap } from "@/domain/settings";
import { type WorkflowNode } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useNotifyChannelsStore } from "@/stores/notify";
import { useWorkflowStore } from "@/stores/workflow";
import { usePanel } from "../PanelProvider";

export type NotifyNodeFormProps = {
  node: WorkflowNode;
};

const initFormModel = () => {
  return {
    subject: "Completed!",
    message: "Your workflow has been completed on Certimate.",
  };
};

const NotifyNodeForm = ({ node }: NotifyNodeFormProps) => {
  const { t } = useTranslation();

  const {
    channels,
    loadedAtOnce: channelsLoadedAtOnce,
    fetchChannels,
  } = useNotifyChannelsStore(useZustandShallowSelector(["channels", "loadedAtOnce", "fetchChannels"]));
  useEffect(() => {
    fetchChannels();
  }, [fetchChannels]);

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));
  const { hidePanel } = usePanel();

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
  const {
    form: formInst,
    formPending,
    formProps,
  } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: node?.config ?? initFormModel(),
    onSubmit: async (values) => {
      await formInst.validateFields();
      await updateNode(
        produce(node, (draft) => {
          draft.config = { ...values };
          draft.validated = true;
        })
      );
      hidePanel();
    },
  });

  return (
    <Form {...formProps} form={formInst} disabled={formPending} layout="vertical">
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

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={formPending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

export default memo(NotifyNodeForm);
