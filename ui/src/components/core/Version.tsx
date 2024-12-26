import { useTranslation } from "react-i18next";
import { Divider, Space, Typography } from "antd";
import { BookOutlined as BookOutlinedIcon } from "@ant-design/icons";

import { version } from "@/domain/version";

export type VersionProps = {
  className?: string;
  style?: React.CSSProperties;
};

const Version = ({ className, style }: VersionProps) => {
  const { t } = useTranslation();

  return (
    <Space className={className} style={style} size={4}>
      <Typography.Link type="secondary" href="https://docs.certimate.me" target="_blank">
        <div className="flex items-center justify-center space-x-1">
          <BookOutlinedIcon />
          <span>{t("common.menu.document")}</span>
        </div>
      </Typography.Link>
      <Divider type="vertical" />
      <Typography.Link type="secondary" href="https://github.com/usual2970/certimate/releases" target="_blank">
        {version}
      </Typography.Link>
    </Space>
  );
};

export default Version;
