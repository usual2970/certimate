import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Radio } from "antd";
import { parseExpression } from "cron-parser";
import { z } from "zod";
import { useShallow } from "zustand/shallow";

import { Button } from "../ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { Input } from "../ui/input";
import { useWorkflowStore, WorkflowState } from "@/stores/workflow";
import { WorkflowNode, WorkflowNodeConfig } from "@/domain/workflow";
import { usePanel } from "./PanelProvider";

const formSchema = z
  .object({
    executionMethod: z.string().min(1, "executionMethod is required"),
    crontab: z.string(),
  })
  .superRefine((data, ctx) => {
    if (data.executionMethod != "auto") {
      return;
    }
    try {
      parseExpression(data.crontab);
    } catch (e) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "crontab is invalid",
        path: ["crontab"],
      });
    }
  });

type StartFormProps = {
  data: WorkflowNode;
};

const i18nPrefix = "workflow.node.start.form";

const selectState = (state: WorkflowState) => ({
  updateNode: state.updateNode,
});
const StartForm = ({ data }: StartFormProps) => {
  const { updateNode } = useWorkflowStore(useShallow(selectState));
  const { hidePanel } = usePanel();

  const { t } = useTranslation();

  const [method, setMethod] = useState("auto");

  useEffect(() => {
    if (data.config && data.config.executionMethod) {
      setMethod(data.config.executionMethod as string);
    } else {
      setMethod("auto");
    }
  }, [data]);

  let config: WorkflowNodeConfig = {
    executionMethod: "auto",
    crontab: "0 0 * * *",
  };
  if (data) config = data.config ?? config;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      executionMethod: config.executionMethod as string,
      crontab: config.crontab as string,
    },
  });

  const onSubmit = async (config: z.infer<typeof formSchema>) => {
    updateNode({ ...data, config: { ...config }, validated: true });
    hidePanel();
  };

  return (
    <>
      <Form {...form}>
        <form
          onSubmit={(e) => {
            e.stopPropagation();
            form.handleSubmit(onSubmit)(e);
          }}
          className="space-y-8 dark:text-stone-200"
        >
          <FormField
            control={form.control}
            name="executionMethod"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t(`${i18nPrefix}.executionMethod.label`)}</FormLabel>
                <FormControl>
                  <Radio.Group
                    {...field}
                    value={method}
                    onChange={(e) => {
                      setMethod(e.target.value);
                    }}
                    className="flex space-x-3"
                  >
                    <Radio value="auto">{t(`${i18nPrefix}.executionMethod.options.auto`)}</Radio>
                    <Radio value="manual">{t(`${i18nPrefix}.executionMethod.options.manual`)}</Radio>
                  </Radio.Group>
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="crontab"
            render={({ field }) => (
              <FormItem hidden={method == "manual"}>
                <FormLabel>{t(`${i18nPrefix}.crontab.label`)}</FormLabel>
                <FormControl>
                  <Input {...field} placeholder={t(`${i18nPrefix}.crontab.placeholder`)} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex justify-end">
            <Button type="submit">{t("common.button.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default StartForm;
