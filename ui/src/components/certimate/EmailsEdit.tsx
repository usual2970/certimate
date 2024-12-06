import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { ClientResponseError } from "pocketbase";

import { cn } from "@/components/ui/utils";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { PbErrorData } from "@/domain/base";
import { EmailsSetting } from "@/domain/settings";
import { update } from "@/repository/settings";
import { useConfigContext } from "@/providers/config";

type EmailsEditProps = {
  className?: string;
  trigger: React.ReactNode;
};

const EmailsEdit = ({ className, trigger }: EmailsEditProps) => {
  const {
    config: { emails },
    setEmails,
  } = useConfigContext();

  const [open, setOpen] = useState(false);
  const { t } = useTranslation();

  const formSchema = z.object({
    email: z.string().email("common.errmsg.email_invalid"),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    if ((emails.content as EmailsSetting).emails.includes(data.email)) {
      form.setError("email", {
        message: "common.errmsg.email_duplicate",
      });
      return;
    }

    // 保存到 config
    const newEmails = [...(emails.content as EmailsSetting).emails, data.email];

    try {
      const resp = await update({
        ...emails,
        name: "emails",
        content: {
          emails: newEmails,
        },
      });

      // 更新本地状态
      setEmails(resp);

      // 关闭弹窗
      form.reset();
      form.clearErrors();

      setOpen(false);
    } catch (e) {
      const err = e as ClientResponseError;

      Object.entries(err.response.data as PbErrorData).forEach(([key, value]) => {
        form.setError(key as keyof z.infer<typeof formSchema>, {
          type: "manual",
          message: value.message,
        });
      });
    }
  };

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild className={cn(className)}>
        {trigger}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[600px] w-full dark:text-stone-200">
        <DialogHeader>
          <DialogTitle>{t("domain.application.form.email.add")}</DialogTitle>
        </DialogHeader>

        <div className="container py-3">
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
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("domain.application.form.email.label")}</FormLabel>
                    <FormControl>
                      <Input placeholder={t("common.errmsg.email_empty")} {...field} type="email" />
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
  );
};

export default EmailsEdit;
