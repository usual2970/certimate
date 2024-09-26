import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { getErrMessage } from "@/lib/error";
import { getPb } from "@/repository/api";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { z } from "zod";

const formSchema = z
  .object({
    oldPassword: z.string().min(10, {
      message: "setting.password.length.message",
    }),
    newPassword: z.string().min(10, {
      message: "setting.password.length.message",
    }),
    confirmPassword: z.string().min(10, {
      message: "setting.password.length.message",
    }),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "setting.password.not.match",
    path: ["confirmPassword"],
  });

const Password = () => {
  const { toast } = useToast();
  const navigate = useNavigate();
  const { t } = useTranslation();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await getPb().admins.authWithPassword(
        getPb().authStore.model?.email,
        values.oldPassword
      );
    } catch (e) {
      const message = getErrMessage(e);
      form.setError("oldPassword", { message });
    }

    try {
      await getPb().admins.update(getPb().authStore.model?.id, {
        password: values.newPassword,
        passwordConfirm: values.confirmPassword,
      });

      getPb().authStore.clear();
      toast({
        title: t('setting.password.change.succeed'),
        description: t("setting.account.log.back.in"),
      });
      setTimeout(() => {
        navigate("/login");
      }, 500);
    } catch (e) {
      const message = getErrMessage(e);
      toast({
        title: t('setting.password.change.failed'),
        description: message,
        variant: "destructive",
      });
    }
  };

  return (
    <>
      <div className="w-full md:max-w-[35em]">
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-8 dark:text-stone-200"
          >
            <FormField
              control={form.control}
              name="oldPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('setting.password.current.password')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('setting.password.current.password')} {...field} type="password" />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="newPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('setting.password.new.password')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="newPassword"
                      {...field}
                      type="password"
                    />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('setting.password.confirm.password')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="confirmPassword"
                      {...field}
                      type="password"
                    />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex justify-end">
              <Button type="submit">{t('setting.submit')}</Button>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default Password;
