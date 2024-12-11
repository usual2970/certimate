import { useEffect, useState } from "react";
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
import { type PbErrorData } from "@/domain/base";
import { useContactStore } from "@/stores/contact";

type EmailsEditProps = {
  className?: string;
  trigger: React.ReactNode;
};

const EmailsEdit = ({ className, trigger }: EmailsEditProps) => {
  const { emails, setEmails, fetchEmails } = useContactStore();

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

  useEffect(() => {
    fetchEmails();
  }, []);

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    if (emails.includes(data.email)) {
      form.setError("email", {
        message: "common.errmsg.email_duplicate",
      });
      return;
    }

    try {
      await setEmails([...emails, data.email]);

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
                <Button type="submit">{t("common.button.save")}</Button>
              </div>
            </form>
          </Form>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default EmailsEdit;
