import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Alert, Button, Form, Input, Radio } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import dayjs from "dayjs";
import { z } from "zod";

import { usePanel } from "../PanelProvider";
import { useZustandShallowSelector } from "@/hooks";
import { type WorkflowNode, type WorkflowNodeConfig } from "@/domain/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { validCronExpression, getNextCronExecutions } from "@/utils/cron";

export type StartNodeFormProps = {
  data: WorkflowNode;
};

const initFormModel = () => {
  return {
    executionMethod: "auto",
    crontab: "0 0 * * *",
  } as WorkflowNodeConfig;
};

const StartNodeForm = ({ data }: StartNodeFormProps) => {
  const { t } = useTranslation();

  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));
  const { hidePanel } = usePanel();

  const formSchema = z
    .object({
      executionMethod: z.string({ message: t("workflow.nodes.start.form.trigger.placeholder") }).min(1, t("workflow.nodes.start.form.trigger.placeholder")),
      crontab: z.string().nullish(),
    })
    .superRefine((data, ctx) => {
      if (data.executionMethod !== "auto") {
        return;
      }

      if (!validCronExpression(data.crontab!)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: t("workflow.nodes.start.form.trigger_cron.errmsg.invalid"),
          path: ["crontab"],
        });
      }
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
      <Form.Item
        name="executionMethod"
        label={t("workflow.nodes.start.form.trigger.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.start.form.trigger.tooltip") }}></span>}
      >
        <Radio.Group value={triggerType} onChange={(e) => handleTriggerTypeChange(e.target.value)}>
          <Radio value="auto">{t("workflow.nodes.start.form.trigger.option.auto.label")}</Radio>
          <Radio value="manual">{t("workflow.nodes.start.form.trigger.option.manual.label")}</Radio>
        </Radio.Group>
      </Form.Item>

      <Form.Item
        name="crontab"
        label={t("workflow.nodes.start.form.trigger_cron.label")}
        hidden={triggerType !== "auto"}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.start.form.trigger_cron.tooltip") }}></span>}
        extra={
          <span>
            {t("workflow.nodes.start.form.trigger_cron.extra")}
            <br />
            {triggerCronLastExecutions.map((d) => (
              <>
                {dayjs(d).format("YYYY-MM-DD HH:mm:ss")}
                <br />
              </>
            ))}
          </span>
        }
      >
        <Input placeholder={t("workflow.nodes.start.form.trigger_cron.placeholder")} onChange={(e) => handleTriggerCronChange(e.target.value)} />
      </Form.Item>

      <Form.Item hidden={triggerType !== "auto"}>
        <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("workflow.nodes.start.form.trigger_cron_alert.content") }}></span>} />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={formPending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

export default StartNodeForm;
