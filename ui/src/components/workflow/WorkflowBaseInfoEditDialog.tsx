import { z } from "zod";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useShallow } from "zustand/shallow";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { useTranslation } from "react-i18next";
import { memo, useEffect, useState } from "react";
import { Textarea } from "../ui/textarea";

type WorkflowNameEditDialogProps = {
  trigger: React.ReactNode;
};

const formSchema = z.object({
  name: z.string(),
  description: z.string(),
});

const selectState = (state: WorkflowState) => ({
  setBaseInfo: state.setBaseInfo,
  workflow: state.workflow,
});
const WorkflowNameBaseInfoDialog = ({ trigger }: WorkflowNameEditDialogProps) => {
  const { setBaseInfo, workflow } = useWorkflowStore(useShallow(selectState));

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  useEffect(() => {
    form.reset({ name: workflow.name, description: workflow.description });
  }, [workflow]);

  const { t } = useTranslation();

  const [open, setOpen] = useState(false);

  const onSubmit = async (config: z.infer<typeof formSchema>) => {
    await setBaseInfo(config.name, config.description);
    setOpen(false);
  };

  return (
    <>
      <Dialog
        open={open}
        onOpenChange={(val) => {
          setOpen(val);
        }}
      >
        <DialogTrigger>{trigger}</DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>基础信息</DialogTitle>
          </DialogHeader>
          <div>
            <Form {...form}>
              <form
                onSubmit={(e) => {
                  e.stopPropagation();
                  form.handleSubmit(onSubmit)(e);
                }}
                className="space-y-8"
              >
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>名称</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="请输入流程名称"
                          {...field}
                          value={field.value}
                          defaultValue={workflow.name}
                          onChange={(e) => {
                            form.setValue("name", e.target.value);
                          }}
                        />
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>说明</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="请输入流程说明"
                          {...field}
                          value={field.value}
                          defaultValue={workflow.description}
                          onChange={(e) => {
                            form.setValue("description", e.target.value);
                          }}
                        />
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="flex justify-end">
                  <Button type="submit">{t("common.save")}</Button>
                </div>
              </form>
            </Form>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
};

export default memo(WorkflowNameBaseInfoDialog);
