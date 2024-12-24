import { cloneElement, useMemo } from "react";
import { useControllableValue } from "ahooks";
import { Drawer } from "antd";

import Show from "@/components/Show";
import CertificateDetail from "./CertificateDetail";
import { type CertificateModel } from "@/domain/certificate";

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

      <Drawer closable destroyOnClose open={open} loading={loading} placement="right" title={data?.id} width={640} onClose={() => setOpen(false)}>
        <Show when={!!data}>
          <CertificateDetail data={data!} />
        </Show>
      </Drawer>
    </>
  );
};

export default CertificateDetailDrawer;
