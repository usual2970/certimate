import React from "react";

import { PanelProvider } from "./PanelProvider";

const WorkflowProvider = ({ children }: { children: React.ReactNode }) => {
  return <PanelProvider>{children}</PanelProvider>;
};

export default WorkflowProvider;
