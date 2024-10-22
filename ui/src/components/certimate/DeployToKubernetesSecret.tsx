import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

const DeployToKubernetesSecret = () => {
  const { t } = useTranslation();
  const { setError } = useDeployEditContext();

  useEffect(() => {
    setError({});
  }, []);

  const { deploy: data, setDeploy } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {
          namespace: "default",
          secretName: "",
          secretDataKeyForCrt: "tls.crt",
          secretDataKeyForKey: "tls.key",
        },
      });
    }
  }, []);

  return (
    <>
      <div className="flex flex-col space-y-8">
        <div>
          <Label>{t("domain.deployment.form.k8s_namespace.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_namespace.label")}
            className="w-full mt-1"
            value={data?.config?.namespace}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.namespace = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.k8s_secret_name.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_secret_name.label")}
            className="w-full mt-1"
            value={data?.config?.secretName}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.secretName = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.k8s_secret_data_key_for_crt.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_secret_data_key_for_crt.label")}
            className="w-full mt-1"
            value={data?.config?.secretDataKeyForCrt}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.secretDataKeyForCrt = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.k8s_secret_data_key_for_key.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_secret_data_key_for_key.label")}
            className="w-full mt-1"
            value={data?.config?.secretDataKeyForKey}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                draft.config ??= {};
                draft.config.secretDataKeyForKey = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>
      </div>
    </>
  );
};

export default DeployToKubernetesSecret;
