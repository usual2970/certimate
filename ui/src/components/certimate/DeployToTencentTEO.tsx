import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";

const DeployToTencentTEO = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  useEffect(() => {
    setError({});
  }, []);

  useEffect(() => {
    const resp = domainSchema.safeParse(data.config?.domain);
    if (!resp.success) {
      setError({
        ...error,
        domain: JSON.parse(resp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        domain: "",
      });
    }
  }, [data]);

  useEffect(() => {
    const resp = zoneIdSchema.safeParse(data.config?.zoneId);
    if (!resp.success) {
      setError({
        ...error,
        zoneId: JSON.parse(resp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        zoneId: "",
      });
    }
  }, [data]);

  const domainSchema = z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
    message: t("common.errmsg.domain_invalid"),
  });

  const zoneIdSchema = z.string().regex(/^zone-[0-9a-zA-Z]{9}$/, {
    message: t("common.errmsg.zoneid_invalid"),
  });

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.tencent_teo_zone_id.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_teo_zone_id.placeholder")}
          className="w-full mt-1"
          value={data?.config?.zoneId}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = zoneIdSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                zoneId: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                zoneId: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.zoneId = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.zoneId}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.tencent_teo_domain.label")}</Label>
        <Textarea
          placeholder={t("domain.deployment.form.tencent_teo_domain.placeholder")}
          className="w-full mt-1"
          value={data?.config?.domain}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = domainSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                domain: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                domain: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.domain = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.domain}</div>
      </div>
    </div>
  );
};

export default DeployToTencentTEO;
