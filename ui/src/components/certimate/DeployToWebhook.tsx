import { useEffect } from "react";
import { produce } from "immer";

import { useDeployEditContext } from "./DeployEdit";
import KVList from "./KVList";
import { type KVType } from "@/domain/domain";

const DeployToWebhook = () => {
  const { config, setConfig, setErrors } = useDeployEditContext();

  useEffect(() => {
    if (!config.id) {
      setConfig({
        ...config,
        config: {},
      });
    }
  }, []);

  useEffect(() => {
    setErrors({});
  }, []);

  return (
    <>
      <KVList
        variables={config?.config?.variables}
        onValueChange={(variables: KVType[]) => {
          const nv = produce(config, (draft) => {
            draft.config ??= {};
            draft.config.variables = variables;
          });
          setConfig(nv);
        }}
      />
    </>
  );
};

export default DeployToWebhook;
