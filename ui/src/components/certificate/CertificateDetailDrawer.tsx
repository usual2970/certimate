import { useControllableValue } from "ahooks";
import { Drawer } from "antd";

import Show from "@/components/Show";
import { type CertificateModel } from "@/domain/certificate";
import { useTriggerElement } from "@/hooks";
import CertificateDetail from "./CertificateDetail";

export type CertificateDetailDrawerProps = {
  data?: CertificateModel;
  loading?: boolean;
  open?: boolean;
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
};

const CertificateDetailDrawer = ({ data, loading, trigger, ...props }: CertificateDetailDrawerProps) => {
  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  return (
    <>
      {triggerEl}

      <Drawer
        afterOpenChange={setOpen}
        destroyOnHidden
        open={open}
        loading={loading}
        placement="right"
        title={`Certificate #${data?.id}`}
        width={720}
        onClose={() => setOpen(false)}
      >
        <Show when={!!data}>
          <CertificateDetail data={data!} />
        </Show>
      </Drawer>
    </>
  );
};

export default CertificateDetailDrawer;
