import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { zodResolver } from "@hookform/resolvers/zod";

import z from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";

import { Access, accessFormType, CloudflareConfig, getUsageByConfigType } from "@/domain/access";
import { save } from "@/repository/access";
import { useConfig } from "@/providers/config";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";

const AccessCloudflareForm = ({
  data,
  onAfterReq,
}: {
  data?: Access;
  onAfterReq: () => void;
}) => {
  const { addAccess, updateAccess } = useConfig();
  const { t } = useTranslation();
  const formSchema = z.object({
    id: z.string().optional(),
    name: z.string().min(1, 'access.form.name.not.empty').max(64, t('zod.rule.string.max', { max: 64 })),
    configType: accessFormType,
    dnsApiToken: z.string().min(1, 'access.form.cloud.dns.api.token.not.empty').max(64, t('zod.rule.string.max', { max: 64 })),
  });

  let config: CloudflareConfig = {
    dnsApiToken: "",
  };
  if (data) config = data.config as CloudflareConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name || '',
      configType: "cloudflare",
      dnsApiToken: config.dnsApiToken,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    console.log(data);
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: getUsageByConfigType(data.configType),
      config: {
        dnsApiToken: data.dnsApiToken,
      },
    };

    try {
      const rs = await save(req);

      onAfterReq();

      req.id = rs.id;
      req.created = rs.created;
      req.updated = rs.updated;
      if (data.id) {
        updateAccess(req);
        return;
      }
      addAccess(req);
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
    <>
      <div className="max-w-[35em] mx-auto mt-10">
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
                  <FormLabel>{t('name')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('access.form.name.not.empty')} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="id"
              render={({ field }) => (
                <FormItem className="hidden">
                  <FormLabel>{t('access.form.config.field')}</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="configType"
              render={({ field }) => (
                <FormItem className="hidden">
                  <FormLabel>{t('access.form.config.field')}</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="dnsApiToken"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('access.form.cloud.dns.api.token')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('access.form.cloud.dns.api.token.not.empty')} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end">
              <Button type="submit">{t('save')}</Button>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default AccessCloudflareForm;
