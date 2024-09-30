import {
  Access,
  accessFormType,
  getUsageByConfigType,
  LocalConfig,
  SSHConfig,
} from "@/domain/access";
import { useConfig } from "@/providers/config";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useTranslation } from "react-i18next";
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
import { Textarea } from "../ui/textarea";
import { save } from "@/repository/access";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";

const AccessLocalForm = ({
  data,
  op,
  onAfterReq,
}: {
  data?: Access;
  op: "add" | "edit" | "copy";
  onAfterReq: () => void;
}) => {
  const { addAccess, updateAccess, reloadAccessGroups } = useConfig();
  const { t } = useTranslation();

  const formSchema = z.object({
    id: z.string().optional(),
    name: z.string().min(1, 'access.form.name.not.empty').max(64, t('zod.rule.string.max', { max: 64 })),
    configType: accessFormType,

    command: z.string().min(1, 'access.form.ssh.command.not.empty').max(2048, t('zod.rule.string.max', { max: 2048 })),
    certPath: z.string().min(0, 'access.form.ssh.cert.path.not.empty').max(2048, t('zod.rule.string.max', { max: 2048 })),
    keyPath: z.string().min(0, 'access.form.ssh.key.path.not.empty').max(2048, t('zod.rule.string.max', { max: 2048 })),
  });

  let config: LocalConfig = {
    command: "sudo service nginx restart",
    certPath: "/etc/nginx/ssl/certificate.crt",
    keyPath: "/etc/nginx/ssl/private.key",
  };
  if (data) config = data.config as SSHConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name || '',
      configType: "local",
      certPath: config.certPath,
      keyPath: config.keyPath,
      command: config.command,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: getUsageByConfigType(data.configType),

      config: {
        command: data.command,
        certPath: data.certPath,
        keyPath: data.keyPath,
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
      } else {
        addAccess(req);
      }

      reloadAccessGroups();
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

      return;
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
            className="space-y-3"
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
              name="certPath"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('access.form.ssh.cert.path')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('access.form.ssh.cert.path.not.empty')} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="keyPath"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('access.form.ssh.key.path')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('access.form.ssh.key.path.not.empty')} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="command"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('access.form.ssh.command')}</FormLabel>
                  <FormControl>
                    <Textarea placeholder={t('access.form.ssh.command.not.empty')} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormMessage />

            <div className="flex justify-end">
              <Button type="submit">{t('save')}</Button>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default AccessLocalForm;
