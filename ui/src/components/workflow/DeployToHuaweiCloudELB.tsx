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

const DeployToHuaweiCloudELB = ({ data }: DeployFormProps) => {
  const { updateNode, getWorkflowOuptutBeforeId } = useWorkflowStore(useZustandShallowSelector(["updateNode", "getWorkflowOuptutBeforeId"]));
  const { hidePanel } = usePanel();
  const { t } = useTranslation();

  const [resourceType, setResourceType] = useState<"certificate" | "loadbalancer" | "listener">();

  useEffect(() => {
    setResourceType(data.config?.resourceType as "certificate" | "loadbalancer" | "listener");
  }, [data]);

  const [beforeOutput, setBeforeOutput] = useState<WorkflowNode[]>([]);

  useEffect(() => {
    const rs = getWorkflowOuptutBeforeId(data.id, "certificate");
    setBeforeOutput(rs);
  }, [data]);

  const formSchema = z
    .object({
      providerType: z.string(),
      access: z.string().min(1, t("domain.deployment.form.access.placeholder")),
      certificate: z.string().min(1),
      region: z.string().min(1, t("domain.deployment.form.huaweicloud_elb_region.placeholder")),
      resourceType: z.union([z.literal("certificate"), z.literal("loadbalancer"), z.literal("listener")], {
        message: t("domain.deployment.form.huaweicloud_elb_resource_type.placeholder"),
      }),
      certificateId: z.string().optional(),
      loadbalancerId: z.string().optional(),
      listenerId: z.string().optional(),
    })
    .refine((data) => (data.resourceType === "certificate" ? !!data.certificateId?.trim() : true), {
      message: t("domain.deployment.form.huaweicloud_elb_certificate_id.placeholder"),
      path: ["certificateId"],
    })
    .refine((data) => (data.resourceType === "loadbalancer" ? !!data.loadbalancerId?.trim() : true), {
      message: t("domain.deployment.form.huaweicloud_elb_loadbalancer_id.placeholder"),
      path: ["loadbalancerId"],
    })
    .refine((data) => (data.resourceType === "listener" ? !!data.listenerId?.trim() : true), {
      message: t("domain.deployment.form.huaweicloud_elb_listener_id.placeholder"),
      path: ["listenerId"],
    });

  let config: WorkflowNodeConfig = {
    certificate: "",
    providerType: "",
    access: "",
    resouceType: "",
    certificateId: "",
    loadbalancerId: "",
    listenerId: "",
  };
  if (data) config = data.config ?? config;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      providerType: "huaweicloud-elb",
      access: config.access as string,
      certificate: config.certificate as string,
      region: config.region as string,
      resourceType: config.resourceType as "certificate" | "loadbalancer" | "listener",
      certificateId: config.certificateId as string,
      loadbalancerId: config.loadbalancerId as string,
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
                    data={{ configType: "huaweicloud" }}
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
                    providerType="huaweicloud-elb"
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
                <FormLabel>{t("domain.deployment.form.huaweicloud_elb_resource_type.label")}</FormLabel>
                <FormControl>
                  <Select
                    {...field}
                    value={field.value}
                    onValueChange={(value) => {
                      setResourceType(value as "certificate" | "loadbalancer" | "listener");
                      form.setValue("resourceType", value as "loadbalancer" | "certificate" | "listener");
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t("domain.deployment.form.huaweicloud_elb_resource_type.placeholder")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="certificate">{t("domain.deployment.form.huaweicloud_elb_resource_type.option.certificate.label")}</SelectItem>
                        <SelectItem value="loadbalancer">{t("domain.deployment.form.huaweicloud_elb_resource_type.option.loadbalancer.label")}</SelectItem>
                        <SelectItem value="listener">{t("domain.deployment.form.huaweicloud_elb_resource_type.option.listener.label")}</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          {resourceType === "certificate" && (
            <FormField
              control={form.control}
              name="certificateId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("domain.deployment.form.huaweicloud_elb_certificate_id.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("domain.deployment.form.huaweicloud_elb_certificate_id.placeholder")} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          {resourceType === "loadbalancer" && (
            <FormField
              control={form.control}
              name="loadbalancerId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("domain.deployment.form.huaweicloud_elb_loadbalancer_id.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("domain.deployment.form.huaweicloud_elb_loadbalancer_id.placeholder")} {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
          )}

          {resourceType === "listener" && (
            <FormField
              control={form.control}
              name="listenerId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t("domain.deployment.form.huaweicloud_elb_listener_id.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("domain.deployment.form.huaweicloud_elb_listener_id.placeholder")} {...field} />
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

export default DeployToHuaweiCloudELB;
