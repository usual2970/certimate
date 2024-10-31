import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";

type DeployToTencentTEOParams = {
  zoneId?: string;
  domain?: string;
};

const DeployToTencentTEO = () => {
  const { t } = useTranslation();

  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToTencentTEOParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {},
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  const formSchema = z.object({
    zoneId: z.string().min(1, t("domain.deployment.form.tencent_teo_zone_id.placeholder")),
    domain: z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: t("common.errmsg.domain_invalid"),
    }),
  });

  useEffect(() => {
    const res = formSchema.safeParse(config.config);
    setErrors({
      ...errors,
      zoneId: res.error?.errors?.find((e) => e.path[0] === "zoneId")?.message,
      domain: res.error?.errors?.find((e) => e.path[0] === "domain")?.message,
    });
  }, [config]);

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.tencent_teo_zone_id.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_teo_zone_id.placeholder")}
          className="w-full mt-1"
          value={config?.config?.zoneId}
          onChange={(e) => {
            const nv = produce(config, (draft) => {
              draft.config ??= {};
              draft.config.zoneId = e.target.value?.trim();
            });
            setConfig(nv);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{errors?.zoneId}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.tencent_teo_domain.label")}</Label>
        <Textarea
          placeholder={t("domain.deployment.form.tencent_teo_domain.placeholder")}
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

export default DeployToTencentTEO;
