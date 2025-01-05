import { createContext, useContext } from "react";
import { type ButtonProps } from "antd";

import { type PanelProps } from "./Panel";

export type ShowPanelOptions = Omit<PanelProps, "defaultOpen" | "open" | "onOpenChange">;
export type ShowPanelWithConfirmOptions = Omit<ShowPanelOptions, "footer" | "onClose"> & {
  cancelButtonProps?: ButtonProps;
  cancelText?: React.ReactNode;
  okButtonProps?: ButtonProps;
  okText?: React.ReactNode;
  onCancel?: () => void;
  onOk?: () => void | Promise<unknown>;
};

export type PanelContextProps = {
  open: boolean;
  show: (options: ShowPanelOptions) => void;
  confirm: (options: ShowPanelWithConfirmOptions) => void;
  hide: () => void;
};

const PanelContext = createContext<PanelContextProps | undefined>(undefined);

export const usePanelContext = () => {
  const context = useContext(PanelContext);
  if (!context) {
    throw new Error("`usePanelContext` must be used within `PanelProvider`");
  }

  return context;
};

export default PanelContext;
