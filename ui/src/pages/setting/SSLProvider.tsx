import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { useToast } from "@/components/ui/use-toast";
import { getErrMessage } from "@/lib/error";
import { cn } from "@/lib/utils";
import { SSLProvider as SSLProviderType, SSLProviderSetting, Setting } from "@/domain/settings";
import { getSetting, update } from "@/repository/settings";

const SSLProvider = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    provider: z.enum(["letsencrypt", "zerossl"], {
      message: t("settings.ca.provider.errmsg.empty"),
    }),
    eabKid: z.string().optional(),
    eabHmacKey: z.string().optional(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      provider: "letsencrypt",
    },
  });

  const [provider, setProvider] = useState("letsencrypt");

  const [config, setConfig] = useState<Setting>();

  const { toast } = useToast();

  useEffect(() => {
    const fetchData = async () => {
      const setting = await getSetting("ssl-provider");

      if (setting) {
        setConfig(setting);
        const content = setting.content as SSLProviderSetting;

        form.setValue("provider", content.provider);
        form.setValue("eabKid", content.config[content.provider].eabKid);
        form.setValue("eabHmacKey", content.config[content.provider].eabHmacKey);
        setProvider(content.provider);
      } else {
        form.setValue("provider", "letsencrypt");
        setProvider("letsencrypt");
      }
    };
    fetchData();
  }, []);

  const getOptionCls = (val: string) => {
    if (provider === val) {
      return "border-primary";
    }

    return "";
  };

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    if (values.provider === "zerossl") {
      if (!values.eabKid) {
        form.setError("eabKid", {
          message: t("settings.ca.eab_kid_hmac_key.errmsg.empty"),
        });
      }
      if (!values.eabHmacKey) {
        form.setError("eabHmacKey", {
          message: t("settings.ca.eab_kid_hmac_key.errmsg.empty"),
        });
      }
      if (!values.eabKid || !values.eabHmacKey) {
        return;
      }
    }

    const setting: Setting = {
      id: config?.id,
      name: "ssl-provider",
      content: {
        provider: values.provider,
        config: {
          letsencrypt: {},
          zerossl: {
            eabKid: values.eabKid ?? "",
            eabHmacKey: values.eabHmacKey ?? "",
          },
        },
      },
    };

    try {
      await update(setting);
      toast({
        title: t("common.update.succeeded.message"),
        description: t("common.update.succeeded.message"),
      });
    } catch (e) {
      const message = getErrMessage(e);
      toast({
        title: t("common.update.failed.message"),
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
              name="provider"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("common.text.ca")}</FormLabel>
                  <FormControl>
                    <RadioGroup
                      {...field}
                      className="flex"
                      onValueChange={(val) => {
                        setProvider(val);
                        form.setValue("provider", val as SSLProviderType);
                      }}
                      value={provider}
                    >
                      <div className="flex items-center space-x-2">
                        <RadioGroupItem value="letsencrypt" id="letsencrypt" />
                        <Label htmlFor="letsencrypt">
                          <div className={cn("flex items-center space-x-2 border p-2 rounded cursor-pointer", getOptionCls("letsencrypt"))}>
                            <img src={"/imgs/providers/letsencrypt.svg"} className="h-6" />
                            <div>{"Let's Encrypt"}</div>
                          </div>
                        </Label>
                      </div>
                      <div className="flex items-center space-x-2">
                        <RadioGroupItem value="zerossl" id="zerossl" />
                        <Label htmlFor="zerossl">
                          <div className={cn("flex items-center space-x-2 border p-2 rounded cursor-pointer", getOptionCls("zerossl"))}>
                            <img src={"/imgs/providers/zerossl.svg"} className="h-6" />
                            <div>{"ZeroSSL"}</div>
                          </div>
                        </Label>
                      </div>
                    </RadioGroup>
                  </FormControl>

                  <FormField
                    control={form.control}
                    name="eabKid"
                    render={({ field }) => (
                      <FormItem hidden={provider !== "zerossl"}>
                        <FormLabel>EAB_KID</FormLabel>
                        <FormControl>
                          <Input placeholder={t("settings.ca.eab_kid.errmsg.empty")} {...field} type="text" />
                        </FormControl>

                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="eabHmacKey"
                    render={({ field }) => (
                      <FormItem hidden={provider !== "zerossl"}>
                        <FormLabel>EAB_HMAC_KEY</FormLabel>
                        <FormControl>
                          <Input placeholder={t("settings.ca.eab_hmac_key.errmsg.empty")} {...field} type="text" />
                        </FormControl>

                        <FormMessage />
                      </FormItem>
                    )}
                  />

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

export default SSLProvider;
