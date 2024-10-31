import { useEffect } from "react";
import { produce } from "immer";

import { useDeployEditContext } from "./DeployEdit";
import KVList from "./KVList";
import { type KVType } from "@/domain/domain";

const DeployToWebhook = () => {
  const { config: data, setConfig: setDeploy, setErrors: setError } = useDeployEditContext();

  useEffect(() => {
    if (!data.id) {
      setDeploy({
        ...data,
        config: {},
      });
    }
  }, []);

  useEffect(() => {
    setError({});
  }, []);

  return (
    <>
      <KVList
        variables={data?.config?.variables}
        onValueChange={(variables: KVType[]) => {
          const nv = produce(data, (draft) => {
            draft.config ??= {};
            draft.config.variables = variables;
          });
          setDeploy(nv);
        }}
      />
    </>
  );
};

export default DeployToWebhook;
