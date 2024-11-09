import { ConfigProvider } from "@/providers/config";
import React from "react";
import { PanelProvider } from "./PanelProvider";
import { NotifyProvider } from "@/providers/notify";

const WorkflowProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <ConfigProvider>
      <NotifyProvider>
        <PanelProvider>{children}</PanelProvider>
      </NotifyProvider>
    </ConfigProvider>
  );
};

export default WorkflowProvider;
