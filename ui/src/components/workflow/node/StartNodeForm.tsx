import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Alert, Button, Form, Input, Radio } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import dayjs from "dayjs";
import { produce } from "immer";
import { z } from "zod";

import Show from "@/components/Show";
import { type WorkflowNode, type WorkflowNodeConfig } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";
import { validCronExpression, getNextCronExecutions } from "@/utils/cron";
import { usePanel } from "../PanelProvider";

export type StartNodeFormProps = {
  data: WorkflowNode;
};

const initFormModel = (): WorkflowNodeConfig => {
  return {
    executionMethod: "auto",
    crontab: "0 0 * * *",
  };
};

const StartNodeForm = ({ data }: StartNodeFormProps) => {
  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));
  const { hidePanel } = usePanel();

  const formSchema = z
    .object({
      executionMethod: z.string({ message: t("workflow_node.start.form.trigger.placeholder") }).min(1, t("workflow_node.start.form.trigger.placeholder")),
      crontab: z.string().nullish(),
    })
    .superRefine((data, ctx) => {
      if (data.executionMethod !== "auto") {
        return;
      }

      if (!validCronExpression(data.crontab!)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: t("workflow_node.start.form.trigger_cron.errmsg.invalid"),
          path: ["crontab"],
        });
      }
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
      await updateNode(
        produce(data, (draft) => {
          draft.config = { ...values };
          draft.validated = true;
        })
      );
      hidePanel();
    },
  });

  const [triggerType, setTriggerType] = useState(data?.config?.executionMethod);
  const [triggerCronLastExecutions, setTriggerCronExecutions] = useState<Date[]>([]);
  useEffect(() => {
    setTriggerType(data?.config?.executionMethod);
    setTriggerCronExecutions(getNextCronExecutions(data?.config?.crontab as string, 5));
  }, [data?.config?.executionMethod, data?.config?.crontab]);

  const handleTriggerTypeChange = (value: string) => {
    setTriggerType(value);

    if (value === "auto") {
      formInst.setFieldValue("crontab", formInst.getFieldValue("crontab") || initFormModel().crontab);
    }
  };

  const handleTriggerCronChange = (value: string) => {
    setTriggerCronExecutions(getNextCronExecutions(value, 5));
  };

  return (
    <Form {...formProps} form={formInst} disabled={formPending} layout="vertical">
      <Form.Item
        name="executionMethod"
        label={t("workflow_node.start.form.trigger.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.start.form.trigger.tooltip") }}></span>}
      >
        <Radio.Group value={triggerType} onChange={(e) => handleTriggerTypeChange(e.target.value)}>
          <Radio value="auto">{t("workflow_node.start.form.trigger.option.auto.label")}</Radio>
          <Radio value="manual">{t("workflow_node.start.form.trigger.option.manual.label")}</Radio>
        </Radio.Group>
      </Form.Item>

      <Form.Item
        name="crontab"
        label={t("workflow_node.start.form.trigger_cron.label")}
        hidden={triggerType !== "auto"}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.start.form.trigger_cron.tooltip") }}></span>}
        extra={
          <Show when={triggerCronLastExecutions.length > 0}>
            <div>
              {t("workflow_node.start.form.trigger_cron.extra")}
              <br />
              {triggerCronLastExecutions.map((date, index) => (
                <span key={index}>
                  {dayjs(date).format("YYYY-MM-DD HH:mm:ss")}
                  <br />
                </span>
              ))}
            </div>
          </Show>
        }
      >
        <Input placeholder={t("workflow_node.start.form.trigger_cron.placeholder")} onChange={(e) => handleTriggerCronChange(e.target.value)} />
      </Form.Item>

      <Form.Item hidden={triggerType !== "auto"}>
        <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.start.form.trigger_cron_alert.content") }}></span>} />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={formPending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

export default memo(StartNodeForm);
