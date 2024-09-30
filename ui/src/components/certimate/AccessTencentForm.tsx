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

import { Access, accessFormType, getUsageByConfigType, TencentConfig } from "@/domain/access";
import { save } from "@/repository/access";
import { useConfig } from "@/providers/config";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";

const AccessTencentForm = ({
  data,
  op,
  onAfterReq,
}: {
  data?: Access;
  op: "add" | "edit" | "copy";
  onAfterReq: () => void;
}) => {
  const { addAccess, updateAccess } = useConfig();
  const { t } = useTranslation();
  const formSchema = z.object({
    id: z.string().optional(),
    name: z.string().min(1, 'access.form.name.not.empty').max(64, t('zod.rule.string.max', { max: 64 })),
    configType: accessFormType,
    secretId: z.string().min(1, 'access.form.secret.id.not.empty').max(64, t('zod.rule.string.max', { max: 64 })),
    secretKey: z.string().min(1, 'access.form.secret.key.not.empty').max(64, t('zod.rule.string.max', { max: 64 })),
  });

  let config: TencentConfig = {
    secretId: "",
    secretKey: "",
  };
  if (data) config = data.config as TencentConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name || '',
      configType: "tencent",
      secretId: config.secretId,
      secretKey: config.secretKey,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: getUsageByConfigType(data.configType),
      config: {
        secretId: data.secretId,
        secretKey: data.secretKey,
      },
    };

    try {
      req.id =  op == "copy" ? "" : req.id;
      const rs = await save(req);

      onAfterReq();

      req.id = rs.id;
      req.created = rs.created;
      req.updated = rs.updated;
      if (data.id && op == "edit") {
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
              name="secretId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('access.form.secret.id')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('access.form.secret.id.not.empty')} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="secretKey"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('access.form.secret.key')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('access.form.secret.key.not.empty')} {...field} />
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

export default AccessTencentForm;
