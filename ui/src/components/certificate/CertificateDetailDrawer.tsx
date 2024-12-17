import { cloneElement, useMemo } from "react";
import { useControllableValue } from "ahooks";
import { Drawer } from "antd";

import { type CertificateModel } from "@/domain/certificate";
import CertificateDetail from "./CertificateDetail";

export type CertificateDetailDrawerProps = {
  data?: CertificateModel;
  loading?: boolean;
  open?: boolean;
  trigger?: React.ReactElement;
  onOpenChange?: (open: boolean) => void;
};

const CertificateDetailDrawer = ({ data, loading, trigger, ...props }: CertificateDetailDrawerProps) => {
  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useMemo(() => {
    if (!trigger) {
      return null;
    }

    return cloneElement(trigger, {
      ...trigger.props,
      onClick: () => {
        setOpen(true);
        trigger.props?.onClick?.();
      },
    });
  }, [trigger, setOpen]);

  return (
    <>
      {triggerEl}

      <Drawer closable destroyOnClose open={open} loading={loading} placement="right" width={480} onClose={() => setOpen(false)}>
        {data ? <CertificateDetail data={data} /> : <></>}
      </Drawer>
    </>
  );
};

export default CertificateDetailDrawer;
