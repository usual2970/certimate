import { cn } from "@/lib/utils";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";

import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { useConfig } from "@/providers/config";
import { update } from "@/repository/settings";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";
import { useState } from "react";
import { EmailsSetting } from "@/domain/settings";

type EmailsEditProps = {
  className?: string;
  trigger: React.ReactNode;
};

const EmailsEdit = ({ className, trigger }: EmailsEditProps) => {
  const {
    config: { emails },
    setEmails,
  } = useConfig();

  const [open, setOpen] = useState(false);

  const formSchema = z.object({
    email: z.string().email(),
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
        message: "邮箱已存在",
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

      Object.entries(err.response.data as PbErrorData).forEach(
        ([key, value]) => {
          form.setError(key as keyof z.infer<typeof formSchema>, {
            type: "manual",
            message: value.message,
          });
        }
      );
    }
  };

  return (
    <Dialog onOpenChange={setOpen} open={open}>
      <DialogTrigger asChild className={cn(className)}>
        {trigger}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[600px] w-full dark:text-stone-200">
        <DialogHeader>
          <DialogTitle>添加邮箱</DialogTitle>
        </DialogHeader>

        <div className="container py-3">
          <Form {...form}>
            <form
              onSubmit={(e) => {
                console.log(e);
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
                    <FormLabel>邮箱</FormLabel>
                    <FormControl>
                      <Input placeholder="请输入邮箱" {...field} type="email" />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="flex justify-end">
                <Button type="submit">保存</Button>
              </div>
            </form>
          </Form>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default EmailsEdit;
