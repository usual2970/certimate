import { useEffect, useState } from "react";
import { Drawer } from "antd";

import { type CertificateModel } from "@/domain/certificate";
import CertificateDetail from "./CertificateDetail";

type CertificateDetailDrawerProps = {
  data?: CertificateModel;
  open?: boolean;
  onClose?: () => void;
};

const CertificateDetailDrawer = ({ data, open, onClose }: CertificateDetailDrawerProps) => {
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    setLoading(data == null);
  }, [data]);

  return (
    <Drawer closable destroyOnClose open={open} loading={loading} placement="right" width={480} onClose={onClose}>
      {data ? <CertificateDetail data={data} /> : <></>}
    </Drawer>
  );
};

export default CertificateDetailDrawer;
