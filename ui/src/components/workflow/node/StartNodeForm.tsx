import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Alert, Button, Form, Input, Radio } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import dayjs from "dayjs";
import { produce } from "immer";
import { z } from "zod";

import Show from "@/components/Show";
import { WORKFLOW_TRIGGERS, type WorkflowNode, type WorkflowNodeConfigForStart } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";
import { getNextCronExecutions, validCronExpression } from "@/utils/cron";
import { usePanel } from "../PanelProvider";

export type StartNodeFormProps = {
  node: WorkflowNode;
};

const initFormModel = (): WorkflowNodeConfigForStart => {
  return {
    trigger: WORKFLOW_TRIGGERS.AUTO,
    triggerCron: "0 0 * * *",
  };
};

const StartNodeForm = ({ node }: StartNodeFormProps) => {
  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));
  const { hidePanel } = usePanel();

  const formSchema = z
    .object({
      trigger: z.string({ message: t("workflow_node.start.form.trigger.placeholder") }).min(1, t("workflow_node.start.form.trigger.placeholder")),
      triggerCron: z.string().nullish(),
    })
    .superRefine((data, ctx) => {
      if (data.trigger !== WORKFLOW_TRIGGERS.AUTO) {
        return;
      }

      if (!validCronExpression(data.triggerCron!)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: t("workflow_node.start.form.trigger_cron.errmsg.invalid"),
          path: ["triggerCron"],
        });
      }
    });
  const formRule = createSchemaFieldRule(formSchema);
  const {
    form: formInst,
    formPending,
    formProps,
  } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: (node?.config as WorkflowNodeConfigForStart) ?? initFormModel(),
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

  const fieldTrigger = Form.useWatch<string>("trigger", formInst);
  const fieldTriggerCron = Form.useWatch<string>("triggerCron", formInst);
  const [fieldTriggerCronExpectedExecutions, setFieldTriggerCronExpectedExecutions] = useState<Date[]>([]);
  useEffect(() => {
    setFieldTriggerCronExpectedExecutions(getNextCronExecutions(fieldTriggerCron, 5));
  }, [fieldTriggerCron]);

  const handleTriggerChange = (value: string) => {
    if (value === WORKFLOW_TRIGGERS.AUTO) {
      formInst.setFieldValue("triggerCron", formInst.getFieldValue("triggerCron") || initFormModel().triggerCron);
    }
  };

  return (
    <Form {...formProps} form={formInst} disabled={formPending} layout="vertical">
      <Form.Item
        name="trigger"
        label={t("workflow_node.start.form.trigger.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.start.form.trigger.tooltip") }}></span>}
      >
        <Radio.Group onChange={(e) => handleTriggerChange(e.target.value)}>
          <Radio value={WORKFLOW_TRIGGERS.AUTO}>{t("workflow_node.start.form.trigger.option.auto.label")}</Radio>
          <Radio value={WORKFLOW_TRIGGERS.MANUAL}>{t("workflow_node.start.form.trigger.option.manual.label")}</Radio>
        </Radio.Group>
      </Form.Item>

      <Form.Item
        name="triggerCron"
        label={t("workflow_node.start.form.trigger_cron.label")}
        hidden={fieldTrigger !== WORKFLOW_TRIGGERS.AUTO}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.start.form.trigger_cron.tooltip") }}></span>}
        extra={
          <Show when={fieldTriggerCronExpectedExecutions.length > 0}>
            <div>
              {t("workflow_node.start.form.trigger_cron.extra")}
              <br />
              {fieldTriggerCronExpectedExecutions.map((date, index) => (
                <span key={index}>
                  {dayjs(date).format("YYYY-MM-DD HH:mm:ss")}
                  <br />
                </span>
              ))}
            </div>
          </Show>
        }
      >
        <Input placeholder={t("workflow_node.start.form.trigger_cron.placeholder")} />
      </Form.Item>

      <Show when={fieldTrigger === WORKFLOW_TRIGGERS.AUTO}>
        <Form.Item>
          <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.start.form.trigger_cron_alert.content") }}></span>} />
        </Form.Item>
      </Show>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={formPending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

export default memo(StartNodeForm);
