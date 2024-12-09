import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { getErrMsg } from "@/utils/error";
import { getPocketBase } from "@/repository/pocketbase";

const formSchema = z
  .object({
    oldPassword: z.string().min(10, {
      message: "settings.password.password.errmsg.length",
    }),
    newPassword: z.string().min(10, {
      message: "settings.password.password.errmsg.length",
    }),
    confirmPassword: z.string().min(10, {
      message: "settings.password.password.errmsg.length",
    }),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "settings.password.password.errmsg.not_matched",
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
      await getPocketBase().admins.authWithPassword(getPocketBase().authStore.model?.email, values.oldPassword);
    } catch (e) {
      const message = getErrMsg(e);
      form.setError("oldPassword", { message });
    }

    try {
      await getPocketBase().admins.update(getPocketBase().authStore.model?.id, {
        password: values.newPassword,
        passwordConfirm: values.confirmPassword,
      });

      getPocketBase().authStore.clear();
      toast({
        title: t("settings.password.changed.message"),
        description: t("settings.account.relogin.message"),
      });
      setTimeout(() => {
        navigate("/login");
      }, 500);
    } catch (e) {
      const message = getErrMsg(e);
      toast({
        title: t("settings.password.failed.message"),
        description: message,
        variant: "destructive",
      });
    }
  };

  return (
    <>
      <div className="w-full md:max-w-[35em]">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 dark:text-stone-200">
            <FormField
              control={form.control}
              name="oldPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("settings.password.current_password.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("settings.password.current_password.placeholder")} {...field} type="password" />
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
                  <FormLabel>{t("settings.password.new_password.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("settings.password.new_password.placeholder")} {...field} type="password" />
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
                  <FormLabel>{t("settings.password.confirm_password.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("settings.password.confirm_password.placeholder")} {...field} type="password" />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex justify-end">
              <Button type="submit">{t("common.update")}</Button>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default Password;
