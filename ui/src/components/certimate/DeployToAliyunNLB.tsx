import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useDeployEditContext } from "./DeployEdit";

type DeployToAliyunNLBConfigParams = {
  region?: string;
  resourceType?: string;
  loadbalancerId?: string;
  listenerId?: string;
};

const DeployToAliyunNLB = () => {
  const { t } = useTranslation();

  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToAliyunNLBConfigParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          region: "cn-hangzhou",
        },
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  const formSchema = z
    .object({
      region: z.string().min(1, t("domain.deployment.form.aliyun_nlb_region.placeholder")),
      resourceType: z.union([z.literal("loadbalancer"), z.literal("listener")], {
        message: t("domain.deployment.form.aliyun_nlb_resource_type.placeholder"),
      }),
      loadbalancerId: z.string().optional(),
      listenerId: z.string().optional(),
    })
    .refine((data) => (data.resourceType === "loadbalancer" ? !!data.loadbalancerId?.trim() : true), {
      message: t("domain.deployment.form.aliyun_nlb_loadbalancer_id.placeholder"),
      path: ["loadbalancerId"],
    })
    .refine((data) => (data.resourceType === "listener" ? !!data.listenerId?.trim() : true), {
      message: t("domain.deployment.form.aliyun_nlb_listener_id.placeholder"),
      path: ["listenerId"],
    });

  useEffect(() => {
    const res = formSchema.safeParse(config.config);
    setErrors({
      ...errors,
      region: res.error?.errors?.find((e) => e.path[0] === "region")?.message,
      resourceType: res.error?.errors?.find((e) => e.path[0] === "resourceType")?.message,
      loadbalancerId: res.error?.errors?.find((e) => e.path[0] === "loadbalancerId")?.message,
      listenerId: res.error?.errors?.find((e) => e.path[0] === "listenerId")?.message,
    });
  }, [config]);

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.aliyun_nlb_region.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.aliyun_nlb_region.placeholder")}
          className="w-full mt-1"
          value={config?.config?.region}
          onChange={(e) => {
            const nv = produce(config, (draft) => {
              draft.config ??= {};
              draft.config.region = e.target.value?.trim();
            });
            setConfig(nv);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{errors?.region}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.aliyun_nlb_resource_type.label")}</Label>
        <Select
          value={config?.config?.resourceType}
          onValueChange={(value) => {
            const nv = produce(config, (draft) => {
              draft.config ??= {};
              draft.config.resourceType = value;
            });
            setConfig(nv);
          }}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("domain.deployment.form.aliyun_nlb_resource_type.placeholder")} />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem value="loadbalancer">{t("domain.deployment.form.aliyun_nlb_resource_type.option.loadbalancer.label")}</SelectItem>
              <SelectItem value="listener">{t("domain.deployment.form.aliyun_nlb_resource_type.option.listener.label")}</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
        <div className="text-red-600 text-sm mt-1">{errors?.resourceType}</div>
      </div>

      {config?.config?.resourceType === "loadbalancer" ? (
        <div>
          <Label>{t("domain.deployment.form.aliyun_nlb_loadbalancer_id.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.aliyun_nlb_loadbalancer_id.placeholder")}
            className="w-full mt-1"
            value={config?.config?.loadbalancerId}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.loadbalancerId = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{errors?.loadbalancerId}</div>
        </div>
      ) : (
        <></>
      )}

      {config?.config?.resourceType === "listener" ? (
        <div>
          <Label>{t("domain.deployment.form.aliyun_nlb_listener_id.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.aliyun_nlb_listener_id.placeholder")}
            className="w-full mt-1"
            value={config?.config?.listenerId}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.listenerId = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{errors?.listenerId}</div>
        </div>
      ) : (
        <></>
      )}
    </div>
  );
};

export default DeployToAliyunNLB;
