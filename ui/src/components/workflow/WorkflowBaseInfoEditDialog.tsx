import { z } from "zod";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { useWorkflowStore } from "@/stores/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { useTranslation } from "react-i18next";
import { memo, useEffect, useMemo, useState } from "react";
import { Textarea } from "../ui/textarea";

type WorkflowNameEditDialogProps = {
  trigger: React.ReactNode;
};

const formSchema = z.object({
  name: z.string(),
  description: z.string(),
});

const WorkflowNameBaseInfoDialog = ({ trigger }: WorkflowNameEditDialogProps) => {
  const { setBaseInfo, workflow } = useWorkflowStore(useZustandShallowSelector(["setBaseInfo", "workflow"]));

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  const memoWorkflow = useMemo(() => workflow, [workflow]);

  useEffect(() => {
    form.reset({ name: workflow.name, description: workflow.description });
  }, [memoWorkflow]);

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
            <DialogTitle className="dark:text-stone-200">{t("workflow.baseinfo.title")}</DialogTitle>
          </DialogHeader>
          <div>
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
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t("workflow.props.name")}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={t("workflow.props.name.placeholder")}
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
                      <FormLabel>{t("workflow.props.description")}</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t("workflow.props.description.placeholder")}
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
                  <Button type="submit">{t("common.button.save")}</Button>
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
