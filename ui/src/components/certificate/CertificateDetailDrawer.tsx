import { cloneElement, useEffect, useMemo, useState } from "react";
import { useControllableValue } from "ahooks";
import { Drawer } from "antd";

import { type CertificateModel } from "@/domain/certificate";
import CertificateDetail from "./CertificateDetail";

type CertificateDetailDrawerProps = {
  data?: CertificateModel;
  open?: boolean;
  trigger?: React.ReactElement;
  onOpenChange?: (open: boolean) => void;
};

const CertificateDetailDrawer = ({ data, trigger, ...props }: CertificateDetailDrawerProps) => {
  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const [loading, setLoading] = useState(true);
  useEffect(() => {
    setLoading(data == null);
  }, [data]);

  const triggerDom = useMemo(() => {
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
      {triggerDom}

      <Drawer closable destroyOnClose open={open} loading={loading} placement="right" width={480} onClose={() => setOpen(false)}>
        {data ? <CertificateDetail data={data} /> : <></>}
      </Drawer>
    </>
  );
};

export default CertificateDetailDrawer;
