import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";

const DeployToIIS = () => {
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
    const passwordResp = passwordSchema.safeParse(data.config?.password);
    if (!passwordResp.success) {
      setError({
        ...error,
        password: JSON.parse(passwordResp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        password: "",
      });
    }
  }, [data]);

  useEffect(() => {
    const siteNameResp = siteNameSchema.safeParse(data.config?.siteName);
    if (!siteNameResp.success) {
      setError({
        ...error,
        siteName: JSON.parse(siteNameResp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        siteName: "",
      });
    }
  }, [data]);

  useEffect(() => {
    const ipResp = ipSchema.safeParse(data.config?.ip);
    if (!ipResp.success) {
      setError({
        ...error,
        ip: JSON.parse(ipResp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        ip: "",
      });
    }
  }, []);

  useEffect(() => {
    const portResp = portSchema.safeParse(data.config?.bindingPort);
    if (!portResp.success) {
      setError({
        ...error,
        bindingPort: JSON.parse(portResp.error.message)[0].message,
      });
    } else {
      setError({
        ...error,
        bindingPort: "",
      });
    }
  }, []);

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          password: "",
          siteName: "",
          ip: "*",
          bindingPort: "443",
          domain: "",
        },
      });
    }
  }, []);

  const domainSchema = z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
    message: t("common.errmsg.domain_invalid"),
  });

  const passwordSchema = z.string().min(1, {
    message: t("domain.deployment.form.iis_password.placeholder"),
  });

  const siteNameSchema = z.string().min(1, {
    message: t("domain.deployment.form.iis_site_name.placeholder"),
  });

  const ipSchema = z.string().regex(/^\*|(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})))$/, {
    message: t("domain.deployment.form.iis_ip.placeholder"),
  });

  const portSchema = z.string().regex(/[0-9]+/, {
    message: t("domain.deployment.form.iis_port.placeholder"),
  });

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.iis_password.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.iis_password.placeholder")}
          className="w-full mt-1"
          value={data?.config?.password}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = passwordSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                password: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                password: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.password = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.password}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.iis_site_name.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.iis_site_name.placeholder")}
          className="w-full mt-1"
          value={data?.config?.siteName}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = siteNameSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                siteName: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                siteName: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.siteName = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.siteName}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.iis_ip.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.iis_ip.placeholder")}
          className="w-full mt-1"
          value={data?.config?.ip}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = ipSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                ip: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                ip: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.ip = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.ip}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.iis_port.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.iis_port.placeholder")}
          className="w-full mt-1"
          value={data?.config?.bindingPort}
          onChange={(e) => {
            const temp = e.target.value;

            const resp = portSchema.safeParse(temp);
            if (!resp.success) {
              setError({
                ...error,
                bindingPort: JSON.parse(resp.error.message)[0].message,
              });
            } else {
              setError({
                ...error,
                bindingPort: "",
              });
            }

            const newData = produce(data, (draft) => {
              if (!draft.config) {
                draft.config = {};
              }
              draft.config.bindingPort = temp;
            });
            setDeploy(newData);
          }}
        />
        <div className="text-red-600 text-sm mt-1">{error?.bindingPort}</div>
      </div>

      <div>
        <Label>{t("domain.deployment.form.iis_domain.label")}</Label>
        <Textarea
          placeholder={t("domain.deployment.form.iis_domain.placeholder")}
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
      <div>
        <Label className="text-red-600">{t("domain.deployment.form.iis_note")}</Label>
      </div>
    </div>
  );
};

export default DeployToIIS;
