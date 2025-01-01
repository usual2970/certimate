import { WorkflowNode, WorkflowNodeType } from "@/domain/workflow";

import DeployPanelBody from "./DeployPanelBody";
import ApplyNodeForm from "./node/ApplyNodeForm";
import NotifyNodeForm from "./node/NotifyNodeForm";
import StartNodeForm from "./node/StartNodeForm";

type PanelBodyProps = {
  data: WorkflowNode;
};
const PanelBody = ({ data }: PanelBodyProps) => {
  const getBody = () => {
    switch (data.type) {
      case WorkflowNodeType.Start:
        return <StartNodeForm data={data} />;
      case WorkflowNodeType.Apply:
        return <ApplyNodeForm data={data} />;
      case WorkflowNodeType.Deploy:
        return <DeployPanelBody data={data} />;
      case WorkflowNodeType.Notify:
        return <NotifyNodeForm data={data} />;
      case WorkflowNodeType.Branch:
        return <div>分支节点</div>;
      case WorkflowNodeType.Condition:
        return <div>条件节点</div>;
      default:
        return <> </>;
    }
  };

  return <>{getBody()}</>;
};

export default PanelBody;
