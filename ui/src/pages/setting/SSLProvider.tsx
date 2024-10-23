import { useContext, useEffect, useState, createContext } from "react";
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
import { produce } from "immer";

type SSLProviderContext = {
  setting: Setting<SSLProviderSetting>;
  onSubmit: (data: Setting<SSLProviderSetting>) => void;
  setConfig: (config: Setting<SSLProviderSetting>) => void;
};

const Context = createContext({} as SSLProviderContext);

export const useSSLProviderContext = () => {
  return useContext(Context);
};

const SSLProvider = () => {
  const { t } = useTranslation();

  const [config, setConfig] = useState<Setting<SSLProviderSetting>>({
    id: "",
    content: {
      provider: "letsencrypt",
      config: {},
    },
  });

  const { toast } = useToast();

  useEffect(() => {
    const fetchData = async () => {
      const setting = await getSetting<SSLProviderSetting>("ssl-provider");

      if (setting) {
        setConfig(setting);
      }
    };
    fetchData();
  }, []);

  const setProvider = (val: SSLProviderType) => {
    const newData = produce(config, (draft) => {
      if (draft.content) {
        draft.content.provider = val;
      } else {
        draft.content = {
          provider: val,
          config: {},
        };
      }
    });
    setConfig(newData);
  };

  const getOptionCls = (val: string) => {
    if (config.content?.provider === val) {
      return "border-primary";
    }

    return "";
  };

  const onSubmit = async (data: Setting<SSLProviderSetting>) => {
    try {
      console.log(data);
      const resp = await update({ ...data });
      setConfig(resp);
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
      <Context.Provider value={{ onSubmit, setConfig, setting: config }}>
        <div className="w-full md:max-w-[35em]">
          <Label>{t("common.text.ca")}</Label>
          <RadioGroup
            className="flex mt-3"
            onValueChange={(val) => {
              setProvider(val as SSLProviderType);
            }}
            value={config.content?.provider}
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
            <div className="flex items-center space-x-2 ">
              <RadioGroupItem value="zerossl" id="zerossl" />
              <Label htmlFor="zerossl">
                <div className={cn("flex items-center space-x-2 border p-2 rounded cursor-pointer", getOptionCls("zerossl"))}>
                  <img src={"/imgs/providers/zerossl.svg"} className="h-6" />
                  <div>{"ZeroSSL"}</div>
                </div>
              </Label>
            </div>

            <div className="flex items-center space-x-2">
              <RadioGroupItem value="gts" id="gts" />
              <Label htmlFor="gts">
                <div className={cn("flex items-center space-x-2 border p-2 rounded cursor-pointer", getOptionCls("gts"))}>
                  <img src={"/imgs/providers/google.svg"} className="h-6" />
                  <div>{"Google Trust Services"}</div>
                </div>
              </Label>
            </div>
          </RadioGroup>

          <SSLProviderForm kind={config.content?.provider ?? ""} />
        </div>
      </Context.Provider>
    </>
  );
};

const SSLProviderForm = ({ kind }: { kind: string }) => {
  const getForm = () => {
    switch (kind) {
      case "zerossl":
        return <SSLProviderZeroSSLForm />;
      case "gts":
        return <SSLProviderGtsForm />;
      default:
        return <SSLProviderLetsEncryptForm />;
    }
  };

  return (
    <>
      <div className="mt-5">{getForm()}</div>
    </>
  );
};

const getConfigStr = (content: SSLProviderSetting, kind: string, key: string) => {
  if (!content.config) {
    return "";
  }
  if (!content.config[kind]) {
    return "";
  }
  return content.config[kind][key] ?? "";
};

const SSLProviderLetsEncryptForm = () => {
  const { t } = useTranslation();

  const { setting, onSubmit } = useSSLProviderContext();

  const formSchema = z.object({
    kind: z.literal("letsencrypt"),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      kind: "letsencrypt",
    },
  });

  const onLocalSubmit = async (data: z.infer<typeof formSchema>) => {
    const newData = produce(setting, (draft) => {
      if (!draft.content) {
        draft.content = {
          provider: data.kind,
          config: {
            letsencrypt: {},
          },
        };
      }
    });
    onSubmit(newData);
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onLocalSubmit)} className="space-y-8 dark:text-stone-200">
        <FormField
          control={form.control}
          name="kind"
          render={({ field }) => (
            <FormItem hidden>
              <FormLabel>kind</FormLabel>
              <FormControl>
                <Input {...field} type="text" />
              </FormControl>

              <FormMessage />
            </FormItem>
          )}
        />

        <FormMessage />

        <div className="flex justify-end">
          <Button type="submit">{t("common.update")}</Button>
        </div>
      </form>
    </Form>
  );
};
const SSLProviderZeroSSLForm = () => {
  const { t } = useTranslation();

  const { setting, onSubmit } = useSSLProviderContext();

  const formSchema = z.object({
    kind: z.literal("zerossl"),
    eabKid: z.string().min(1, { message: t("settings.ca.eab_kid_hmac_key.errmsg.empty") }),
    eabHmacKey: z.string().min(1, { message: t("settings.ca.eab_kid_hmac_key.errmsg.empty") }),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      kind: "zerossl",
      eabKid: "",
      eabHmacKey: "",
    },
  });

  useEffect(() => {
    if (setting.content) {
      const content = setting.content;

      form.reset({
        eabKid: getConfigStr(content, "zerossl", "eabKid"),
        eabHmacKey: getConfigStr(content, "zerossl", "eabHmacKey"),
      });
    }
  }, [setting]);

  const onLocalSubmit = async (data: z.infer<typeof formSchema>) => {
    const newData = produce(setting, (draft) => {
      if (!draft.content) {
        draft.content = {
          provider: "zerossl",
          config: {
            zerossl: {},
          },
        };
      }

      draft.content.config.zerossl = {
        eabKid: data.eabKid,
        eabHmacKey: data.eabHmacKey,
      };
    });
    onSubmit(newData);
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onLocalSubmit)} className="space-y-8 dark:text-stone-200">
        <FormField
          control={form.control}
          name="kind"
          render={({ field }) => (
            <FormItem hidden>
              <FormLabel>kind</FormLabel>
              <FormControl>
                <Input {...field} type="text" />
              </FormControl>

              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="eabKid"
          render={({ field }) => (
            <FormItem>
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
            <FormItem>
              <FormLabel>EAB_HMAC_KEY</FormLabel>
              <FormControl>
                <Input placeholder={t("settings.ca.eab_hmac_key.errmsg.empty")} {...field} type="text" />
              </FormControl>

              <FormMessage />
            </FormItem>
          )}
        />

        <FormMessage />

        <div className="flex justify-end">
          <Button type="submit">{t("common.update")}</Button>
        </div>
      </form>
    </Form>
  );
};

const SSLProviderGtsForm = () => {
  const { t } = useTranslation();

  const { setting, onSubmit } = useSSLProviderContext();

  const formSchema = z.object({
    kind: z.literal("gts"),
    eabKid: z.string().min(1, { message: t("settings.ca.eab_kid_hmac_key.errmsg.empty") }),
    eabHmacKey: z.string().min(1, { message: t("settings.ca.eab_kid_hmac_key.errmsg.empty") }),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      kind: "gts",
      eabKid: "",
      eabHmacKey: "",
    },
  });

  useEffect(() => {
    if (setting.content) {
      const content = setting.content;

      form.reset({
        eabKid: getConfigStr(content, "gts", "eabKid"),
        eabHmacKey: getConfigStr(content, "gts", "eabHmacKey"),
      });
    }
  }, [setting]);

  const onLocalSubmit = async (data: z.infer<typeof formSchema>) => {
    const newData = produce(setting, (draft) => {
      if (!draft.content) {
        draft.content = {
          provider: "gts",
          config: {
            zerossl: {},
          },
        };
      }

      draft.content.config.gts = {
        eabKid: data.eabKid,
        eabHmacKey: data.eabHmacKey,
      };
    });
    onSubmit(newData);
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onLocalSubmit)} className="space-y-8 dark:text-stone-200">
        <FormField
          control={form.control}
          name="kind"
          render={({ field }) => (
            <FormItem hidden>
              <FormLabel>kind</FormLabel>
              <FormControl>
                <Input {...field} type="text" />
              </FormControl>

              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="eabKid"
          render={({ field }) => (
            <FormItem>
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
            <FormItem>
              <FormLabel>EAB_HMAC_KEY</FormLabel>
              <FormControl>
                <Input placeholder={t("settings.ca.eab_hmac_key.errmsg.empty")} {...field} type="text" />
              </FormControl>

              <FormMessage />
            </FormItem>
          )}
        />

        <FormMessage />

        <div className="flex justify-end">
          <Button type="submit">{t("common.update")}</Button>
        </div>
      </form>
    </Form>
  );
};

export default SSLProvider;
