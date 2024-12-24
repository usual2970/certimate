import { useTranslation } from "react-i18next";
import { z } from "zod";

import { Input } from "@/components/ui/input";
import { DeployFormProps } from "./DeployForm";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { WorkflowNode, WorkflowNodeConfig } from "@/domain/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { usePanel } from "./PanelProvider";
import { Button } from "../ui/button";

import { useEffect, useState } from "react";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { SelectLabel } from "@radix-ui/react-select";
import AccessEditModal from "../access/AccessEditModal";
import { Plus } from "lucide-react";
import AccessSelect from "./AccessSelect";

const DeployToKubernetesSecret = ({ data }: DeployFormProps) => {
  const { updateNode, getWorkflowOuptutBeforeId } = useWorkflowStore(useZustandShallowSelector(["updateNode", "getWorkflowOuptutBeforeId"]));
  const { hidePanel } = usePanel();
  const { t } = useTranslation();

  const [beforeOutput, setBeforeOutput] = useState<WorkflowNode[]>([]);

  useEffect(() => {
    const rs = getWorkflowOuptutBeforeId(data.id, "certificate");
    console.log(rs);
    setBeforeOutput(rs);
  }, [data]);

  const formSchema = z.object({
    providerType: z.string(),
    access: z.string().min(1, t("domain.deployment.form.access.placeholder")),
    certificate: z.string().min(1),
    namespace: z.string().min(1, {
      message: t("domain.deployment.form.k8s_namespace.placeholder"),
    }),
    secretName: z.string().min(1, {
      message: t("domain.deployment.form.k8s_secret_name.placeholder"),
    }),
    secretDataKeyForCrt: z.string().min(1, {
      message: t("domain.deployment.form.k8s_secret_data_key_for_crt.placeholder"),
    }),
    secretDataKeyForKey: z.string().min(1, {
      message: t("domain.deployment.form.k8s_secret_data_key_for_key.placeholder"),
    }),
  });

  let config: WorkflowNodeConfig = {
    certificate: "",
    providerType: "",
    access: "",
    namespace: "",
    secretName: "",
    secretDataKeyForCrt: "",
    secretDataKeyForKey: "",
  };
  if (data) config = data.config ?? config;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      providerType: "k8s-secret",
      access: config.access as string,
      certificate: config.certificate as string,
      namespace: config.namespace as string,
      secretName: config.secretName as string,
      secretDataKeyForCrt: config.secretDataKeyForCrt as string,
      secretDataKeyForKey: config.secretDataKeyForKey as string,
    },
  });

  const onSubmit = async (config: z.infer<typeof formSchema>) => {
    updateNode({ ...data, config: { ...config }, validated: true });
    hidePanel();
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
            name="access"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="flex justify-between">
                  <div>{t("domain.deployment.form.access.label")}</div>

                  <AccessEditModal
                    data={{ configType: "k8s" }}
                    mode="add"
                    trigger={
                      <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                        <Plus size={14} />
                        {t("common.button.add")}
                      </div>
                    }
                  />
                </FormLabel>
                <FormControl>
                  <AccessSelect
                    {...field}
                    value={field.value}
                    onValueChange={(value) => {
                      form.setValue("access", value);
                    }}
                    providerType="k8s-secret"
                  />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="certificate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("workflow.common.certificate.label")}</FormLabel>
                <FormControl>
                  <Select
                    {...field}
                    value={field.value}
                    onValueChange={(value) => {
                      form.setValue("certificate", value);
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t("workflow.common.certificate.placeholder")} />
                    </SelectTrigger>
                    <SelectContent>
                      {beforeOutput.map((item) => (
                        <>
                          <SelectGroup key={item.id}>
                            <SelectLabel>{item.name}</SelectLabel>
                            {item.output?.map((output) => (
                              <SelectItem key={output.name} value={`${item.id}#${output.name}`}>
                                <div>
                                  {item.name}-{output.label}
                                </div>
                              </SelectItem>
                            ))}
                          </SelectGroup>
                        </>
                      ))}
                    </SelectContent>
                  </Select>
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="namespace"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("domain.deployment.form.k8s_namespace.label")}</FormLabel>
                <FormControl>
                  <Input {...field} placeholder={t("domain.deployment.form.k8s_namespace.label")} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="secretName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("domain.deployment.form.k8s_secret_name.label")}</FormLabel>
                <FormControl>
                  <Input {...field} placeholder={t("domain.deployment.form.k8s_secret_name.label")} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="secretDataKeyForCrt"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("domain.deployment.form.k8s_secret_data_key_for_crt.label")}</FormLabel>
                <FormControl>
                  <Input {...field} placeholder={t("domain.deployment.form.k8s_secret_data_key_for_crt.label")} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="secretDataKeyForKey"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("domain.deployment.form.k8s_secret_data_key_for_key.label")}</FormLabel>
                <FormControl>
                  <Input {...field} placeholder={t("domain.deployment.form.k8s_secret_data_key_for_key.label")} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex justify-end">
            <Button type="submit">{t("common.button.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default DeployToKubernetesSecret;
