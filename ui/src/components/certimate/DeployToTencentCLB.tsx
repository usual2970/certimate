import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

const DeployToTencentCLB = () => {
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
    const clbIdresp = clbIdSchema.safeParse(data.config?.clbId);
    if (!clbIdresp.success) {
      setError({
        ...error,
        clbId: JSON.parse(clbIdresp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        clbId: "",
      });
    }
  }, [data]);

  useEffect(() => {
    const lsnIdresp = lsnIdSchema.safeParse(data.config?.lsnId);
    if (!lsnIdresp.success) {
      setError({
        ...error,
        lsnId: JSON.parse(lsnIdresp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        lsnId: "",
      });
    }
  }, [data]);

  useEffect(() => {
    const regionResp = regionSchema.safeParse(data.config?.region);
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
          lsnId: "",
          clbId: "",
          domain: "",
          region: "",
        },
      });
    }
  }, []);

  const regionSchema = z.string().regex(/^ap-[a-z]+$/, {
    message: t("domain.deployment.form.clb_region.placeholder"),
  });

  const domainSchema = z.string().regex(/^$|^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
    message: t("common.errmsg.domain_invalid"),
  });

  const clbIdSchema = z.string().regex(/^lb-[a-zA-Z0-9]{8}$/, {
    message: t("domain.deployment.form.clb_id.placeholder"),
  });

  const lsnIdSchema = z.string().regex(/^lbl-.{8}$/, {
    message: t("domain.deployment.form.clb_listener.placeholder"),
  });

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.clb_region.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.clb_region.placeholder")}
          className="w-full mt-1"
          value={data?.config?.region}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = regionSchema.safeParse(temp);
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
        <Label>{t("domain.deployment.form.clb_id.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.clb_id.placeholder")}
          className="w-full mt-1"
          value={data?.config?.clbId}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = clbIdSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                clbId: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                clbId: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.clbId = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.clbId}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.clb_listener.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.clb_listener.placeholder")}
          className="w-full mt-1"
          value={data?.config?.lsnId}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = lsnIdSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                lsnId: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                lsnId: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.lsnId = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.lsnId}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.clb_domain.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.clb_domain.placeholder")}
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

export default DeployToTencentCLB;
