import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

const DeployToTencentCOS = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          region: "ap-guangzhou",
          bucket: "",
          domain: "",
        },
      });
    }
  }, []);

  useEffect(() => {
    setError({});
  }, []);

  const formSchema = z.object({
    region: z.string().min(1, t("domain.deployment.form.tencent_cos_region.placeholder")),
    bucket: z.string().min(1, t("domain.deployment.form.tencent_cos_bucket.placeholder")),
    domain: z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: t("common.errmsg.domain_invalid"),
    }),
  });

  useEffect(() => {
    const res = formSchema.safeParse(data.config);
    setError({
      ...error,
      region: res.error?.errors?.find((e) => e.path[0] === "region")?.message,
      bucket: res.error?.errors?.find((e) => e.path[0] === "bucket")?.message,
      domain: res.error?.errors?.find((e) => e.path[0] === "domain")?.message,
    });
  }, [data]);

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.tencent_cos_region.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_cos_region.placeholder")}
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
        <Label>{t("domain.deployment.form.tencent_cos_bucket.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.tencent_cos_bucket.placeholder")}
          className="w-full mt-1"
          value={data?.config?.bucket}
          onChange={(e) => {
            const newData = produce(data, (draft) => {
              draft.config ??= {};
              draft.config.bucket = e.target.value?.trim();
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.bucket}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.domain.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.domain.placeholder")}
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
    </div>
  );
};

export default DeployToTencentCOS;
