import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { produce } from "immer";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useDeployEditContext } from "./DeployEdit";

const DeployToSSH = () => {
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
          certPath: "/etc/nginx/ssl/nginx.crt",
          keyPath: "/etc/nginx/ssl/nginx.key",
          preCommand: "",
          command: "sudo service nginx reload",
        },
      });
    }
  }, []);

  return (
    <>
      <div className="flex flex-col space-y-8">
        <div>
          <Label>{t("domain.deployment.form.ssh_cert_path.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.ssh_cert_path.label")}
            className="w-full mt-1"
            value={data?.config?.certPath}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.certPath = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.ssh_key_path.label")}</Label>
          <Input
            placeholder={t("domain.deployment.form.ssh_key_path.placeholder")}
            className="w-full mt-1"
            value={data?.config?.keyPath}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.keyPath = e.target.value;
              });
              setDeploy(newData);
            }}
          />
        </div>

        <div>
          <Label>{t("domain.deployment.form.ssh_pre_command.label")}</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.preCommand}
            placeholder={t("domain.deployment.form.ssh_pre_command.placeholder")}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.preCommand = e.target.value;
              });
              setDeploy(newData);
            }}
          ></Textarea>
        </div>

        <div>
          <Label>{t("domain.deployment.form.ssh_command.label")}</Label>
          <Textarea
            className="mt-1"
            value={data?.config?.command}
            placeholder={t("domain.deployment.form.ssh_command.placeholder")}
            onChange={(e) => {
              const newData = produce(data, (draft) => {
                if (!draft.config) {
                  draft.config = {};
                }
                draft.config.command = e.target.value;
              });
              setDeploy(newData);
            }}
          ></Textarea>
        </div>
      </div>
    </>
  );
};

export default DeployToSSH;
