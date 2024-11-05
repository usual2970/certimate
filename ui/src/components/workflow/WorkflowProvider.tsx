import { ConfigProvider } from "@/providers/config";
import React from "react";
import { PanelProvider } from "./PanelProvider";

const WorkflowProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <ConfigProvider>
      <PanelProvider>{children}</PanelProvider>
    </ConfigProvider>
  );
};

export default WorkflowProvider;
