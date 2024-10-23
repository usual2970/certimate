import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

const DeployToTencentCOS = () => {
  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  const { t } = useTranslation();

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
    const bucketResp = bucketSchema.safeParse(data.config?.bucket);
    if (!bucketResp.success) {
      setError({
        ...error,
        bucket: JSON.parse(bucketResp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        bucket: "",
      });
    }
  }, []);

  useEffect(() => {
    const regionResp = regionSchema.safeParse(data.config?.bucket);
    if (!regionResp.success) {
      setError({
        ...error,
        region: JSON.parse(regionResp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        region: "",
      });
    }
  }, []);

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          region: "",
          bucket: "",
          domain: "",
        },
      });
    }
  }, []);

  const domainSchema = z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
    message: t("common.errmsg.domain_invalid"),
  });

  const regionSchema = z.string().regex(/^ap-[a-z]+$/, {
    message: t("domain.deployment.form.cos_region.placeholder"),
  });

  const bucketSchema = z.string().regex(/^.+-\d+$/, {
    message: t("domain.deployment.form.cos_bucket.placeholder"),
  });

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.cos_region.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.cos_region.placeholder")}
          className="w-full mt-1"
          value={data?.config?.region}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = bucketSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                region: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                region: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.region = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.region}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.cos_bucket.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.cos_bucket.placeholder")}
          className="w-full mt-1"
          value={data?.config?.bucket}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = bucketSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                bucket: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                bucket: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.bucket = temp;
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

export default DeployToTencentCOS;
