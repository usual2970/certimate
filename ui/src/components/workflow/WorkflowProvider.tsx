import React from "react";

import { NotifyProvider } from "@/providers/notify";
import { PanelProvider } from "./PanelProvider";

const WorkflowProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <NotifyProvider>
      <PanelProvider>{children}</PanelProvider>
    </NotifyProvider>
  );
};

export default WorkflowProvider;
