import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

const DeployToQiniuCDN = () => {
  const { t } = useTranslation();

  const { deploy: data, setDeploy, error, setError } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          domain: "",
        },
      });
    }
  }, []);

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

  const domainSchema = z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
    message: t("common.errmsg.domain_invalid"),
  });

  return (
    <div className="flex flex-col space-y-8">
      <div>
        <Label>{t("domain.deployment.form.domain.label")}</Label>
        <Input
          placeholder={t("domain.deployment.form.domain.placeholder")}
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

export default DeployToQiniuCDN;
