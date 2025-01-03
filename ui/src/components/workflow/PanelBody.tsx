import { type WorkflowNode, WorkflowNodeType } from "@/domain/workflow";

import ApplyNodeForm from "./node/ApplyNodeForm";
import DeployNodeForm from "./node/DeployNodeForm";
import NotifyNodeForm from "./node/NotifyNodeForm";
import StartNodeForm from "./node/StartNodeForm";

type PanelBodyProps = {
  data: WorkflowNode;
};

const PanelBody = ({ data }: PanelBodyProps) => {
  const getBody = () => {
    switch (data.type) {
      case WorkflowNodeType.Start:
        return <StartNodeForm node={data} />;
      case WorkflowNodeType.Apply:
        return <ApplyNodeForm node={data} />;
      case WorkflowNodeType.Deploy:
        return <DeployNodeForm node={data} />;
      case WorkflowNodeType.Notify:
        return <NotifyNodeForm node={data} />;
      default:
        return <> </>;
    }
  };

  return <>{getBody()}</>;
};

export default PanelBody;
