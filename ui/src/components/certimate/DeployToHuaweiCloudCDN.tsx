import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

type DeployToHuaweiCloudCDNConfigParams = {
  region?: string;
  domain?: string;
};

const DeployToHuaweiCloudCDN = () => {
  const { t } = useTranslation();

  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToHuaweiCloudCDNConfigParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          region: "cn-north-1",
        },
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  const formSchema = z.object({
    region: z.string().min(1, {
      message: t("domain.deployment.form.huaweicloud_cdn_region.placeholder"),
    }),
    domain: z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: t("common.errmsg.domain_invalid"),
    }),
  });

  useEffect(() => {
    const res = formSchema.safeParse(config.config);
    setErrors({
      ...errors,
      region: res.error?.errors?.find((e) => e.path[0] === "region")?.message,
      domain: res.error?.errors?.find((e) => e.path[0] === "domain")?.message,
    });
  }, [config]);

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.huaweicloud_cdn_region.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.huaweicloud_cdn_region.placeholder")}
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
        <Label>{t("domain.deployment.form.domain.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.domain.placeholder")}
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
    </div>
  );
};

export default DeployToHuaweiCloudCDN;
