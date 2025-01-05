import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Alert, Form, type FormInstance, Input, Radio } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import dayjs from "dayjs";
import { z } from "zod";

import Show from "@/components/Show";
import { WORKFLOW_TRIGGERS, type WorkflowNode, type WorkflowNodeConfigForStart, type WorkflowTriggerType } from "@/domain/workflow";
import { getNextCronExecutions, validCronExpression } from "@/utils/cron";

type StartNodeFormFieldValues = Partial<WorkflowNodeConfigForStart>;

export type StartNodeFormProps = {
  form: FormInstance;
  formName?: string;
  disabled?: boolean;
  workflowNode: WorkflowNode;
  onValuesChange?: (values: StartNodeFormFieldValues) => void;
};

const initFormModel = (): StartNodeFormFieldValues => {
  return {
    trigger: WORKFLOW_TRIGGERS.AUTO,
    triggerCron: "0 0 * * *",
  };
};

const StartNodeForm = ({ form, formName, disabled, workflowNode, onValuesChange }: StartNodeFormProps) => {
  const { t } = useTranslation();

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

  const initialValues: StartNodeFormFieldValues = (workflowNode.config as WorkflowNodeConfigForStart) ?? initFormModel();

  const fieldTrigger = Form.useWatch<WorkflowTriggerType>("trigger", form);
  const fieldTriggerCron = Form.useWatch<string>("triggerCron", form);
  const [fieldTriggerCronExpectedExecutions, setFieldTriggerCronExpectedExecutions] = useState<Date[]>([]);
  useEffect(() => {
    setFieldTriggerCronExpectedExecutions(getNextCronExecutions(fieldTriggerCron, 5));
  }, [fieldTriggerCron]);

  const handleTriggerChange = (value: string) => {
    if (value === WORKFLOW_TRIGGERS.AUTO) {
      form.setFieldValue("triggerCron", initialValues.triggerCron || initFormModel().triggerCron);
    } else {
      form.setFieldValue("triggerCron", undefined);
    }
  };

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as StartNodeFormFieldValues);
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
    </Form>
  );
};

export default memo(StartNodeForm);
