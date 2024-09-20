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
import { update } from "@/repository/access_group";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";
import { useState } from "react";

type AccessGroupEditProps = {
  className?: string;
  trigger: React.ReactNode;
};

const AccessGroupEdit = ({ className, trigger }: AccessGroupEditProps) => {
  const { reloadAccessGroups } = useConfig();

  const [open, setOpen] = useState(false);

  const formSchema = z.object({
    name: z.string().min(1).max(64),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    try {
      await update({
        name: data.name,
      });

      // 更新本地状态
      reloadAccessGroups();

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
          <DialogTitle>添加分组</DialogTitle>
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
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>组名</FormLabel>
                    <FormControl>
                      <Input placeholder="请输入组名" {...field} type="text" />
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

export default AccessGroupEdit;
