import { useTranslation } from "react-i18next";
import { z } from "zod";

import { Input } from "@/components/ui/input";
import { DeployFormProps } from "./DeployForm";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { WorkflowNode, WorkflowNodeConfig } from "@/domain/workflow";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useShallow } from "zustand/shallow";
import { usePanel } from "./PanelProvider";
import { Button } from "../ui/button";

import { useEffect, useState } from "react";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { SelectLabel } from "@radix-ui/react-select";

type TencentResourceType = "ssl-deploy" | "loadbalancer" | "listener" | "ruledomain";

const selectState = (state: WorkflowState) => ({
  updateNode: state.updateNode,
  getWorkflowOuptutBeforeId: state.getWorkflowOuptutBeforeId,
});
const DeployToTencentCLB = ({ data }: DeployFormProps) => {
  const { updateNode, getWorkflowOuptutBeforeId } = useWorkflowStore(useShallow(selectState));
  const { hidePanel } = usePanel();
  const { t } = useTranslation();

  const [resourceType, setResourceType] = useState<TencentResourceType>();

  useEffect(() => {
    setResourceType(data.config?.resourceType as TencentResourceType);
  }, [data]);

  const [beforeOutput, setBeforeOutput] = useState<WorkflowNode[]>([]);

  useEffect(() => {
    const rs = getWorkflowOuptutBeforeId(data.id, "certificate");
    setBeforeOutput(rs);
  }, [data]);

  const formSchema = z
    .object({
      providerType: z.string(),
      certificate: z.string().min(1),
      region: z.string().min(1, t("domain.deployment.form.tencent_clb_region.placeholder")),
      resourceType: z.union([z.literal("ssl-deploy"), z.literal("loadbalancer"), z.literal("listener"), z.literal("ruledomain")], {
        message: t("domain.deployment.form.tencent_clb_resource_type.placeholder"),
      }),
      loadbalancerId: z.string().min(1, t("domain.deployment.form.tencent_clb_loadbalancer_id.placeholder")),
      listenerId: z.string().optional(),
      domain: z.string().regex(/^$|^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
        message: t("common.errmsg.domain_invalid"),
      }),
    })
    .refine(
      (data) => {
        switch (data.resourceType) {
          case "ssl-deploy":
          case "listener":
          case "ruledomain":
            return !!data.listenerId?.trim();
        }
        return true;
      },
      {
        message: t("domain.deployment.form.tencent_clb_listener_id.placeholder"),
        path: ["listenerId"],
      }
    )
    .refine((data) => (data.resourceType === "ruledomain" ? !!data.domain?.trim() : true), {
      message: t("domain.deployment.form.tencent_clb_ruledomain.placeholder"),
      path: ["domain"],
    });

  let config: WorkflowNodeConfig = {
    certificate: "",
    providerType: "tencent-clb",
    resouceType: "",
    domain: "",
    loadbalancerId: "",
    listenerId: "",
  };
  if (data) config = data.config ?? config;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      providerType: config.providerType as string,
      certificate: config.certificate as string,
      region: config.region as string,
      resourceType: config.resourceType as TencentResourceType,
      domain: config.certificateId as string,
      loadbalancerId: config.loadbalancerId as string,
      listenerId: config.listenerId as string,
    },
  });

  const onSubmit = async (config: z.infer<typeof formSchema>) => {
    updateNode({ ...data, config: { ...config } });
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
            name="certificate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>证书</FormLabel>
                <FormControl>
                  <Select
                    {...field}
                    value={field.value}
                    onValueChange={(value) => {
                      form.setValue("certificate", value);
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="选择证书来源" />
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
                <FormLabel>{t("domain.deployment.form.huaweicloud_cdn_region.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("domain.deployment.form.huaweicloud_cdn_region.placeholder")} {...field} />
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
                <FormLabel>{t("domain.deployment.form.tencent_clb_resource_type.label")}</FormLabel>
                <FormControl>
                  <Select
                    {...field}
                    value={field.value}
                    onValueChange={(value) => {
                      setResourceType(value as TencentResourceType);
                      form.setValue("resourceType", value as TencentResourceType);
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t("domain.deployment.form.tencent_clb_resource_type.placeholder")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="ssl-deploy">{t("domain.deployment.form.tencent_clb_resource_type.option.ssl_deploy.label")}</SelectItem>
                        <SelectItem value="loadbalancer">{t("domain.deployment.form.tencent_clb_resource_type.option.loadbalancer.label")}</SelectItem>
                        <SelectItem value="listener">{t("domain.deployment.form.tencent_clb_resource_type.option.listener.label")}</SelectItem>
                        <SelectItem value="ruledomain">{t("domain.deployment.form.tencent_clb_resource_type.option.ruledomain.label")}</SelectItem>
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
                <FormLabel>{t("domain.deployment.form.tencent_clb_loadbalancer_id.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("domain.deployment.form.tencent_clb_loadbalancer_id.placeholder")} {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          {(resourceType === "ssl-deploy" || resourceType === "listener" || resourceType === "ruledomain") && (
            <FormField
              control={form.control}
              name="listenerId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("domain.deployment.form.tencent_clb_listener_id.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("domain.deployment.form.tencent_clb_listener_id.placeholder")} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          {resourceType === "ssl-deploy" && (
            <FormField
              control={form.control}
              name="domain"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("domain.deployment.form.tencent_clb_domain.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("domain.deployment.form.tencent_clb_domain.placeholder")} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          {resourceType === "ruledomain" && (
            <FormField
              control={form.control}
              name="domain"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("domain.deployment.form.tencent_clb_ruledomain.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("domain.deployment.form.tencent_clb_ruledomain.placeholder")} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          <div className="flex justify-end">
            <Button type="submit">{t("common.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default DeployToTencentCLB;
