import { useEffect } from "react";
import { produce } from "immer";

import { useDeployEditContext } from "./DeployEdit";
import KVList from "./KVList";
import { type KVType } from "@/domain/domain";

const DeployToWebhook = () => {
  const { deploy: data, setDeploy } = useDeployEditContext();

  const { setError } = useDeployEditContext();

  useEffect(() => {
    setError({});
  }, []);

  return (
    <>
      <KVList
        variables={data?.config?.variables}
        onValueChange={(variables: KVType[]) => {
          const newData = produce(data, (draft) => {
            if (!draft.config) {
              draft.config = {};
            }
            draft.config.variables = variables;
          });
          setDeploy(newData);
        }}
      />
    </>
  );
};

export default DeployToWebhook;
