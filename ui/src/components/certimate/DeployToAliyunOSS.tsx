import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

const DeployToAliyunOSS = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          endpoint: "oss.aliyuncs.com",
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
    endpoint: z.string().min(1, {
      message: t("domain.deployment.form.aliyun_oss_endpoint.placeholder"),
    }),
    bucket: z.string().min(1, {
      message: t("domain.deployment.form.aliyun_oss_bucket.placeholder"),
    }),
    domain: z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: t("common.errmsg.domain_invalid"),
    }),
  });

  useEffect(() => {
    const res = formSchema.safeParse(data.config);
    setError({
      ...error,
      endpoint: res.error?.errors?.find((e) => e.path[0] === "endpoint")?.message,
      bucket: res.error?.errors?.find((e) => e.path[0] === "bucket")?.message,
      domain: res.error?.errors?.find((e) => e.path[0] === "domain")?.message,
    });
  }, [data]);

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.aliyun_oss_endpoint.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.aliyun_oss_endpoint.placeholder")}
          className="w-full mt-1"
          value={data?.config?.endpoint}
          onChange={(e) => {
            const newData = produce(data, (draft) => {
              draft.config ??= {};
              draft.config.endpoint = e.target.value?.trim();
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.endpoint}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.aliyun_oss_bucket.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.aliyun_oss_bucket.placeholder")}
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
          placeholder={t("domain.deployment.form.domain.label")}
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

export default DeployToAliyunOSS;
