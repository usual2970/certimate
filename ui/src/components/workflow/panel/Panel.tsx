import React, { useEffect } from "react";
import { useControllableValue } from "ahooks";
import { Drawer } from "antd";

export type PanelProps = {
  children: React.ReactNode;
  defaultOpen?: boolean;
  extra?: React.ReactNode;
  footer?: React.ReactNode;
  open?: boolean;
  title?: React.ReactNode;
  onClose?: () => void | Promise<unknown>;
  onOpenChange?: (open: boolean) => void;
};

const Panel = ({ children, extra, footer, title, onClose, ...props }: PanelProps) => {
  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const handleClose = async () => {
    try {
      const ret = await onClose?.();
      if (ret != null && !ret) return;

      setOpen(false);
    } catch {
      return;
    }
  };

  return (
    <Drawer destroyOnClose extra={extra} footer={footer} open={open} title={title} width={640} onClose={handleClose}>
      {children}
    </Drawer>
  );
};

export default Panel;
