import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useDeployEditContext } from "./DeployEdit";

const DeployToTencentCLB = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          region: "ap-guangzhou",
          resourceType: "",
          loadbalancerId: "",
          listenerId: "",
          domain: "",
        },
      });
    }
  }, []);

  useEffect(() => {
    setError({});
  }, []);

  const formSchema = z
    .object({
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

  useEffect(() => {
    const res = formSchema.safeParse(data.config);
    setError({
      ...error,
      region: res.error?.errors?.find((e) => e.path[0] === "region")?.message,
      resourceType: res.error?.errors?.find((e) => e.path[0] === "resourceType")?.message,
      loadbalancerId: res.error?.errors?.find((e) => e.path[0] === "loadbalancerId")?.message,
      listenerId: res.error?.errors?.find((e) => e.path[0] === "listenerId")?.message,
      domain: res.error?.errors?.find((e) => e.path[0] === "domain")?.message,
    });
  }, [data]);

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.tencent_clb_region.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_clb_region.placeholder")}
          className="w-full mt-1"
          value={data?.config?.region}
          onChange={(e) => {
            const newData = produce(data, (draft) => {
              draft.config ??= {};
              draft.config.region = e.target.value?.trim();
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.region}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.tencent_clb_resource_type.label")}</Label>
        <Select
          value={data?.config?.resourceType}
          onValueChange={(value) => {
            const newData = produce(data, (draft) => {
              draft.config ??= {};
              draft.config.resourceType = value;
            });
            setDeploy(newData);
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
        <div className="text-red-600 text-sm mt-1">{error?.resourceType}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.tencent_clb_loadbalancer_id.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_clb_loadbalancer_id.placeholder")}
          className="w-full mt-1"
          value={data?.config?.loadbalancerId}
          onChange={(e) => {
            const newData = produce(data, (draft) => {
              draft.config ??= {};
              draft.config.loadbalancerId = e.target.value?.trim();
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.loadbalancerId}</div>
      </div>

      {data?.config?.resourceType === "ssl-deploy" || data?.config?.resourceType === "listener" || data?.config?.resourceType === "ruledomain" ? (
        <div>
          <Label>{t("domain.deployment.form.tencent_clb_listener_id.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.tencent_clb_listener_id.placeholder")}
            className="w-full mt-1"
            value={data?.config?.listenerId}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.listenerId = e.target.value?.trim();
              });
              setDeploy(newData);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{error?.listenerId}</div>
        </div>
      ) : (
        <></>
      )}

      {data?.config?.resourceType === "ssl-deploy" ? (
        <div>
          <Label>{t("domain.deployment.form.tencent_clb_domain.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.tencent_clb_domain.placeholder")}
            className="w-full mt-1"
            value={data?.config?.domain}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.domain = e.target.value?.trim();
              });
              setDeploy(newData);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{error?.domain}</div>
        </div>
      ) : (
        <></>
      )}

      {data?.config?.resourceType === "ruledomain" ? (
        <div>
          <Label>{t("domain.deployment.form.tencent_clb_ruledomain.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.tencent_clb_ruledomain.placeholder")}
            className="w-full mt-1"
            value={data?.config?.domain}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.domain = e.target.value?.trim();
              });
              setDeploy(newData);
            }}
          />
          <div className="text-red-600 text-sm mt-1">{error?.domain}</div>
        </div>
      ) : (
        <></>
      )}
    </div>
  );
};

export default DeployToTencentCLB;
