import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { ClientResponseError } from "pocketbase";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { PbErrorData } from "@/domain/base";
import { accessProvidersMap, accessTypeFormSchema, type Access, type HttpreqConfig } from "@/domain/access";
import { save } from "@/repository/access";
import { useConfigContext } from "@/providers/config";

type AccessHttpreqFormProps = {
  op: "add" | "edit" | "copy";
  data?: Access;
  onAfterReq: () => void;
};

const AccessHttpreqForm = ({ data, op, onAfterReq }: AccessHttpreqFormProps) => {
  const { addAccess, updateAccess } = useConfigContext();
  const { t } = useTranslation();
  const formSchema = z.object({
    id: z.string().optional(),
    name: z
      .string()
      .min(1, "access.authorization.form.name.placeholder")
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    configType: accessTypeFormSchema,
    endpoint: z.string().url("common.errmsg.url_invalid"),
    mode: z.enum(["RAW", ""]),
    username: z
      .string()
      .min(1, "access.authorization.form.access_key_secret.placeholder")
      .max(128, t("common.errmsg.string_max", { max: 128 })),
    password: z
      .string()
      .min(1, "access.authorization.form.access_key_secret.placeholder")
      .max(128, t("common.errmsg.string_max", { max: 128 })),
  });

  let config: HttpreqConfig = {
    endpoint: "",
    mode: "",
    username: "",
    password: "",
  };
  if (data) config = data.config as HttpreqConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name || "",
      configType: "httpreq",
      endpoint: config.endpoint,
      mode: config.mode === "RAW" ? "RAW" : "",
      username: config.username,
      password: config.password,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: accessProvidersMap.get(data.configType)!.usage,
      config: {
        endpoint: data.endpoint,
        mode: data.mode,
        username: data.username,
        password: data.password,
      },
    };

    try {
      req.id = op == "copy" ? "" : req.id;
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

      Object.entries(err.response.data as PbErrorData).forEach(([key, value]) => {
        form.setError(key as keyof z.infer<typeof formSchema>, {
          type: "manual",
          message: value.message,
        });
      });

      return;
    }
  };

  return (
    <>
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
                <FormLabel>{t("access.authorization.form.name.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.name.placeholder")} {...field} />
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
                <FormLabel>{t("access.authorization.form.config.label")}</FormLabel>
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
                <FormLabel>{t("access.authorization.form.config.label")}</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="endpoint"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.httpreq_endpoint.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.httpreq_endpoint.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="mode"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.httpreq_mode.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.httpreq_mode.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.username.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.username.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.password.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.password.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormMessage />

          <div className="flex justify-end">
            <Button type="submit">{t("common.button.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default AccessHttpreqForm;
