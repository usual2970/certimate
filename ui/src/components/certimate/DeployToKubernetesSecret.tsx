import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeployEditContext } from "./DeployEdit";

type DeployToKubernetesSecretConfigParams = {
  namespace?: string;
  secretName?: string;
  secretDataKeyForCrt?: string;
  secretDataKeyForKey?: string;
};

const DeployToKubernetesSecret = () => {
  const { t } = useTranslation();

  const { config, setConfig, errors, setErrors } = useDeployEditContext<DeployToKubernetesSecretConfigParams>();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {
          namespace: "default",
          secretDataKeyForCrt: "tls.crt",
          secretDataKeyForKey: "tls.key",
        },
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  const formSchema = z.object({
    namespace: z.string().min(1, {
      message: t("domain.deployment.form.k8s_namespace.placeholder"),
    }),
    secretName: z.string().min(1, {
      message: t("domain.deployment.form.k8s_secret_name.placeholder"),
    }),
    secretDataKeyForCrt: z.string().min(1, {
      message: t("domain.deployment.form.k8s_secret_data_key_for_crt.placeholder"),
    }),
    secretDataKeyForKey: z.string().min(1, {
      message: t("domain.deployment.form.k8s_secret_data_key_for_key.placeholder"),
    }),
  });

  useEffect(() => {
    const res = formSchema.safeParse(config.config);
    setErrors({
      ...errors,
      namespace: res.error?.errors?.find((e) => e.path[0] === "namespace")?.message,
      secretName: res.error?.errors?.find((e) => e.path[0] === "secretName")?.message,
      secretDataKeyForCrt: res.error?.errors?.find((e) => e.path[0] === "secretDataKeyForCrt")?.message,
      secretDataKeyForKey: res.error?.errors?.find((e) => e.path[0] === "secretDataKeyForKey")?.message,
    });
  }, [config]);

  return (
    <>
      <div className="flex flex-col space-y-8">
        <div>
          <Label>{t("domain.deployment.form.k8s_namespace.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_namespace.label")}
            className="w-full mt-1"
            value={config?.config?.namespace}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.namespace = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.k8s_secret_name.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_secret_name.label")}
            className="w-full mt-1"
            value={config?.config?.secretName}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.secretName = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.k8s_secret_data_key_for_crt.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_secret_data_key_for_crt.label")}
            className="w-full mt-1"
            value={config?.config?.secretDataKeyForCrt}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.secretDataKeyForCrt = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.k8s_secret_data_key_for_key.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.k8s_secret_data_key_for_key.label")}
            className="w-full mt-1"
            value={config?.config?.secretDataKeyForKey}
            onChange={(e) => {
              const nv = produce(config, (draft) => {
                draft.config ??= {};
                draft.config.secretDataKeyForKey = e.target.value?.trim();
              });
              setConfig(nv);
            }}
          />
        </div>
      </div>
    </>
  );
};

export default DeployToKubernetesSecret;
