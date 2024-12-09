import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { getErrMsg } from "@/utils/error";
import { getPocketBase } from "@/repository/pocketbase";

const formSchema = z.object({
  email: z.string().email("settings.account.email.errmsg.invalid"),
});

const Account = () => {
  const { toast } = useToast();
  const navigate = useNavigate();
  const { t } = useTranslation();

  const [changed, setChanged] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: getPocketBase().authStore.model?.email,
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await getPocketBase().admins.update(getPocketBase().authStore.model?.id, {
        email: values.email,
      });

      getPocketBase().authStore.clear();
      toast({
        title: t("settings.account.email.changed.message"),
        description: t("settings.account.relogin.message"),
      });
      setTimeout(() => {
        navigate("/login");
      }, 500);
    } catch (e) {
      const message = getErrMsg(e);
      toast({
        title: t("settings.account.email.failed.message"),
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
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("settings.account.email.label")}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t("settings.account.email.placeholder")}
                      {...field}
                      type="email"
                      onChange={(e) => {
                        setChanged(true);
                        form.setValue("email", e.target.value);
                      }}
                    />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end">
              {changed ? (
                <Button type="submit">{t("common.update")}</Button>
              ) : (
                <Button type="submit" disabled variant={"secondary"}>
                  {t("common.update")}
                </Button>
              )}
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default Account;
