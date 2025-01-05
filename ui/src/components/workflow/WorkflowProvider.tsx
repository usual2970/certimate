import { PanelProvider } from "./panel/PanelProvider";

const WorkflowProvider = ({ children }: { children: React.ReactNode }) => {
  return <PanelProvider>{children}</PanelProvider>;
};

export default WorkflowProvider;
