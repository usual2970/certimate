import { createContext, useContext, useState } from "react";

import Panel from "./Panel";

type PanelContentProps = { name: string; children: React.ReactNode };

type PanelContextProps = {
  open: boolean;
  showPanel: ({ name, children }: PanelContentProps) => void;
  hidePanel: () => void;
};

const PanelContext = createContext<PanelContextProps | undefined>(undefined);

export const PanelProvider = ({ children }: { children: React.ReactNode }) => {
  const [open, setOpen] = useState(false);
  const [panelContent, setPanelContent] = useState<PanelContentProps | null>(null);

  const showPanel = (panelContent: PanelContentProps) => {
    setOpen(true);
    setPanelContent(panelContent);
  };

  const hidePanel = () => {
    setOpen(false);
    setPanelContent(null);
  };

  return (
    <PanelContext.Provider value={{ open, showPanel, hidePanel }}>
      {children}
      <Panel open={open} onOpenChange={setOpen} children={panelContent?.children} name={panelContent?.name ?? ""} />
    </PanelContext.Provider>
  );
};

export const usePanel = () => {
  const context = useContext(PanelContext);
  if (!context) {
    throw new Error("`usePanel` must be used within PanelProvider");
  }
  return context;
};
