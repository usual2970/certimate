import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";

import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import { DeployFormProps } from "./DeployForm";
import { WorkflowNode, WorkflowNodeConfig } from "@/domain/workflow";
import { useWorkflowStore, WorkflowState } from "@/stores/workflow";
import { useShallow } from "zustand/shallow";
import { usePanel } from "./PanelProvider";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { Button } from "../ui/button";

import AccessSelect from "./AccessSelect";
import AccessEditDialog from "../certimate/AccessEditDialog";
import { Plus } from "lucide-react";

const selectState = (state: WorkflowState) => ({
  updateNode: state.updateNode,
  getWorkflowOuptutBeforeId: state.getWorkflowOuptutBeforeId,
});

const DeployToAliyunCLB = ({ data }: DeployFormProps) => {
  const { updateNode, getWorkflowOuptutBeforeId } = useWorkflowStore(useShallow(selectState));
  const { hidePanel } = usePanel();
  const { t } = useTranslation();

  const [beforeOutput, setBeforeOutput] = useState<WorkflowNode[]>([]);

  useEffect(() => {
    const rs = getWorkflowOuptutBeforeId(data.id, "certificate");
    console.log(rs);
    setBeforeOutput(rs);
  }, [data]);

  const formSchema = z
    .object({
      providerType: z.string(),
      access: z.string().min(1, t("domain.deployment.form.access.placeholder")),
      certificate: z.string().min(1),
      region: z.string().min(1, t("domain.deployment.form.aliyun_clb_region.placeholder")),
      resourceType: z.union([z.literal("certificate"), z.literal("loadbalancer"), z.literal("listener")], {
        message: t("domain.deployment.form.aliyun_clb_resource_type.placeholder"),
      }),
      loadbalancerId: z.string().optional(),
      listenerPort: z.string().optional(),
    })
    .refine((data) => (data.resourceType === "loadbalancer" || data.resourceType === "listener" ? !!data.loadbalancerId?.trim() : true), {
      message: t("domain.deployment.form.aliyun_clb_loadbalancer_id.placeholder"),
      path: ["loadbalancerId"],
    })
    .refine((data) => (data.resourceType === "listener" ? +data.listenerPort! > 0 && +data.listenerPort! < 65535 : true), {
      message: t("domain.deployment.form.aliyun_clb_listener_port.placeholder"),
      path: ["listenerPort"],
    });

  let config: WorkflowNodeConfig = {
    certificate: "",
    providerType: "",
    access: "",
    region: "",
    resourceType: "",
    loadbalancerId: "",
    listenerPort: "",
  };
  if (data) config = data.config ?? config;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      providerType: "aliyun-clb",
      access: config.access as string,
      certificate: config.certificate as string,
      region: config.region as string,
      resourceType: config.resourceType as "loadbalancer" | "listener",
      loadbalancerId: config.loadbalancerId as string,
      listenerPort: config.listenerPort as string,
    },
  });

  const resouceType = form.watch("resourceType");

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

                  <AccessEditDialog
                    trigger={
                      <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                        <Plus size={14} />
                        {t("common.button.add")}
                      </div>
                    }
                    op="add"
                    outConfigType="aliyun"
                  />
                </FormLabel>
                <FormControl>
                  <AccessSelect
                    {...field}
                    value={field.value}
                    onValueChange={(value) => {
                      form.setValue("access", value);
                    }}
                    providerType="aliyun-clb"
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
            name="region"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("domain.deployment.form.aliyun_clb_region.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("domain.deployment.form.aliyun_clb_region.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="resourceType"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("domain.deployment.form.aliyun_clb_resource_type.label")}</FormLabel>
                <FormControl>
                  <Select
                    {...field}
                    value={field.value}
                    onValueChange={(value) => {
                      form.setValue("resourceType", value as "loadbalancer" | "listener");
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t("domain.deployment.form.aliyun_clb_resource_type.placeholder")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="loadbalancer">{t("domain.deployment.form.aliyun_clb_resource_type.option.loadbalancer.label")}</SelectItem>
                        <SelectItem value="listener">{t("domain.deployment.form.aliyun_clb_resource_type.option.listener.label")}</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="loadbalancerId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("domain.deployment.form.aliyun_clb_loadbalancer_id.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("domain.deployment.form.aliyun_clb_loadbalancer_id.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          {resouceType === "listener" && (
            <FormField
              control={form.control}
              name="listenerPort"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("domain.deployment.form.aliyun_clb_listener_port.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("domain.deployment.form.aliyun_clb_listener_port.placeholder")} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          <div className="flex justify-end">
            <Button type="submit">{t("common.button.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default DeployToAliyunCLB;
