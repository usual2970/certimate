import { WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import StartForm from "./StartForm";
import DeployPanelBody from "./DeployPanelBody";
import ApplyForm from "./ApplyForm";
import NotifyForm from "./NotifyForm";

type PanelBodyProps = {
  data: WorkflowNode;
};
const PanelBody = ({ data }: PanelBodyProps) => {
  const getBody = () => {
    switch (data.type) {
      case WorkflowNodeType.Start:
        return <StartForm data={data} />;
      case WorkflowNodeType.Apply:
        return <ApplyForm data={data} />;
      case WorkflowNodeType.Deploy:
        return <DeployPanelBody data={data} />;
      case WorkflowNodeType.Notify:
        return <NotifyForm data={data} />;
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

