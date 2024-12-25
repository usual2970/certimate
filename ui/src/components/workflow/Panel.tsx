import { useEffect } from "react";
import { Drawer } from "antd";

type AddNodePanelProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  children: React.ReactNode;
  name: string;
};

const Panel = ({ open, onOpenChange, children, name }: AddNodePanelProps) => {
  useEffect(() => {
    onOpenChange(open);
  }, [open, onOpenChange]);

  return (
    <Drawer open={open} title={name} width={640} onClose={() => onOpenChange(false)}>
      {children}
    </Drawer>
  );
};

export default Panel;
