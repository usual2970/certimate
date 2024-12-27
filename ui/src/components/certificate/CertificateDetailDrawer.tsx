import { useControllableValue } from "ahooks";
import { Drawer } from "antd";

import Show from "@/components/Show";
import CertificateDetail from "./CertificateDetail";
import { useTriggerElement } from "@/hooks";
import { type CertificateModel } from "@/domain/certificate";

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

  const triggerDom = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  return (
    <>
      {triggerDom}

      <Drawer
        afterOpenChange={setOpen}
        closable
        destroyOnClose
        open={open}
        loading={loading}
        placement="right"
        title={data?.id}
        width={640}
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
