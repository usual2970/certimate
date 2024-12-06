import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { ClientResponseError } from "pocketbase";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { readFileContent } from "@/lib/file";
import { PbErrorData } from "@/domain/base";
import { accessProvidersMap, accessTypeFormSchema, type Access, type KubernetesConfig } from "@/domain/access";
import { save } from "@/repository/access";
import { useConfigContext } from "@/providers/config";

type AccessKubernetesFormProps = {
  op: "add" | "edit" | "copy";
  data?: Access;
  onAfterReq: () => void;
};

const AccessKubernetesForm = ({ data, op, onAfterReq }: AccessKubernetesFormProps) => {
  const { addAccess, updateAccess } = useConfigContext();

  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [fileName, setFileName] = useState("");

  const { t } = useTranslation();

  const formSchema = z.object({
    id: z.string().optional(),
    name: z
      .string()
      .min(1, "access.authorization.form.name.placeholder")
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    configType: accessTypeFormSchema,
    kubeConfig: z
      .string()
      .min(0, "access.authorization.form.k8s_kubeconfig.placeholder")
      .max(20480, t("common.errmsg.string_max", { max: 20480 })),
    kubeConfigFile: z.any().optional(),
  });

  let config: KubernetesConfig & { kubeConfigFile?: string } = {
    kubeConfig: "",
    kubeConfigFile: "",
  };
  if (data) config = data.config as typeof config;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name || "",
      configType: "k8s",
      kubeConfig: config.kubeConfig,
      kubeConfigFile: config.kubeConfigFile,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: accessProvidersMap.get(data.configType)!.usage,
      config: {
        kubeConfig: data.kubeConfig,
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
      } else {
        addAccess(req);
      }
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

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;
    const savedFile = file;
    setFileName(savedFile.name);
    const content = await readFileContent(savedFile);
    form.setValue("kubeConfig", content);
  };

  const handleSelectFileClick = () => {
    fileInputRef.current?.click();
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
            name="kubeConfig"
            render={({ field }) => (
              <FormItem hidden>
                <FormLabel>{t("access.authorization.form.k8s_kubeconfig.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.k8s_kubeconfig.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="kubeConfigFile"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.k8s_kubeconfig.label")}</FormLabel>
                <FormControl>
                  <div>
                    <Button type={"button"} variant={"secondary"} size={"sm"} className="w-48" onClick={handleSelectFileClick}>
                      {fileName ? fileName : t("access.authorization.form.k8s_kubeconfig_file.placeholder")}
                    </Button>
                    <Input
                      placeholder={t("access.authorization.form.k8s_kubeconfig.placeholder")}
                      {...field}
                      ref={fileInputRef}
                      className="hidden"
                      hidden
                      type="file"
                      onChange={handleFileChange}
                    />
                  </div>
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormMessage />

          <div className="flex justify-end">
            <Button type="submit">{t("common.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default AccessKubernetesForm;
