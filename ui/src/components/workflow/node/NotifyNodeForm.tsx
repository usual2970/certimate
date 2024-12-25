import { useEffect, useState } from "react";
import { Link } from "react-router";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Button, Form, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";
import { ChevronRight as ChevronRightIcon } from "lucide-react";

import { usePanel } from "../PanelProvider";
import { useZustandShallowSelector } from "@/hooks";
import { notifyChannelsMap } from "@/domain/settings";
import { type WorkflowNode, type WorkflowNodeConfig } from "@/domain/workflow";
import { useNotifyChannelStore } from "@/stores/notify";
import { useWorkflowStore } from "@/stores/workflow";

export type NotifyNodeFormProps = {
  data: WorkflowNode;
};

const initFormModel = () => {
  return {
    subject: "",
    message: "",
  } as WorkflowNodeConfig;
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
  const [formInst] = Form.useForm<z.infer<typeof formSchema>>();
  const [formPending, setFormPending] = useState(false);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(
    (data?.config as Partial<z.infer<typeof formSchema>>) ?? initFormModel()
  );
  useDeepCompareEffect(() => {
    setInitialValues((data?.config as Partial<z.infer<typeof formSchema>>) ?? initFormModel());
  }, [data?.config]);

  const handleFormFinish = async (values: z.infer<typeof formSchema>) => {
    setFormPending(true);

    try {
      await updateNode({ ...data, config: { ...values }, validated: true });

      hidePanel();
    } finally {
      setFormPending(false);
    }
  };

  return (
    <Form form={formInst} disabled={formPending} initialValues={initialValues} layout="vertical" onFinish={handleFormFinish}>
      <Form.Item name="subject" label={t("workflow.nodes.notify.form.subject.label")} rules={[formRule]}>
        <Input placeholder={t("workflow.nodes.notify.form.subject.placeholder")} />
      </Form.Item>

      <Form.Item name="message" label={t("workflow.nodes.notify.form.message.label")} rules={[formRule]}>
        <Input.TextArea autoSize={{ minRows: 3, maxRows: 5 }} placeholder={t("workflow.nodes.notify.form.message.placeholder")} />
      </Form.Item>

      <Form.Item name="channel" rules={[formRule]}>
        <label className="block mb-1">
          <div className="flex items-center justify-between gap-4 w-full overflow-hidden">
            <div className="flex-grow max-w-full truncate">{t("workflow.nodes.notify.form.channel.label")}</div>
            <div className="text-right">
              <Link className="ant-typography" to="/settings/notification" target="_blank">
                <Button className="p-0" type="link">
                  {t("workflow.nodes.notify.form.channel.button")}
                  <ChevronRightIcon size={14} />
                </Button>
              </Link>
            </div>
          </div>
        </label>
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

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={formPending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

export default NotifyNodeForm;
