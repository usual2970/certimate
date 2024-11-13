import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useDeployEditContext } from "./DeployEdit";

type DeployToTencentCLBParams = {
  region?: string;
  resourceType?: string;
  loadbalancerId?: string;
  listenerId?: string;
  domain?: string;
};

const DeployToTencentCLB = () => {
  const { t } = useTranslation();

  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToTencentCLBParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          region: "ap-guangzhou",
        },
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  const formSchema = z
    .object({
      region: z.string().min(1, t("domain.deployment.form.tencent_clb_region.placeholder")),
      resourceType: z.union([z.literal("ssl-deploy"), z.literal("loadbalancer"), z.literal("listener"), z.literal("ruledomain")], {
        message: t("domain.deployment.form.tencent_clb_resource_type.placeholder"),
      }),
      loadbalancerId: z.string().min(1, t("domain.deployment.form.tencent_clb_loadbalancer_id.placeholder")),
      listenerId: z.string().optional(),
      domain: z.string().optional(),
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
    .refine(
      (data) => {
        switch (data.resourceType) {
          case "ssl-deploy":
          case "ruledomain":
            return !!data.domain?.trim() && /^$|^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/.test(data.domain);
        }
        return true;
      },
      {
        message: t("domain.deployment.form.tencent_clb_ruledomain.placeholder"),
        path: ["domain"],
      }
    );

  useEffect(() => {
    const res = formSchema.safeParse(config.config);
    setErrors({
      ...errors,
      region: res.error?.errors?.find((e) => e.path[0] === "region")?.message,
      resourceType: res.error?.errors?.find((e) => e.path[0] === "resourceType")?.message,
      loadbalancerId: res.error?.errors?.find((e) => e.path[0] === "loadbalancerId")?.message,
      listenerId: res.error?.errors?.find((e) => e.path[0] === "listenerId")?.message,
      domain: res.error?.errors?.find((e) => e.path[0] === "domain")?.message,
    });
  }, [config]);

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.tencent_clb_region.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_clb_region.placeholder")}
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
        <Label>{t("domain.deployment.form.tencent_clb_resource_type.label")}</Label>
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
        <div className="text-red-600 text-sm mt-1">{errors?.resourceType}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.tencent_clb_loadbalancer_id.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_clb_loadbalancer_id.placeholder")}
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

      {config?.config?.resourceType === "ssl-deploy" || config?.config?.resourceType === "listener" || config?.config?.resourceType === "ruledomain" ? (
        <div>
          <Label>{t("domain.deployment.form.tencent_clb_listener_id.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.tencent_clb_listener_id.placeholder")}
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

      {config?.config?.resourceType === "ssl-deploy" ? (
        <div>
          <Label>{t("domain.deployment.form.tencent_clb_domain.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.tencent_clb_domain.placeholder")}
            className="w-full mt-1"
            value={config?.config?.domain}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.domain = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{errors?.domain}</div>
        </div>
      ) : (
        <></>
      )}

      {config?.config?.resourceType === "ruledomain" ? (
        <div>
          <Label>{t("domain.deployment.form.tencent_clb_ruledomain.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.tencent_clb_ruledomain.placeholder")}
            className="w-full mt-1"
            value={config?.config?.domain}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.domain = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{errors?.domain}</div>
        </div>
      ) : (
        <></>
      )}
    </div>
  );
};

export default DeployToTencentCLB;
