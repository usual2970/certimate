import { useTranslation } from "react-i18next";
import { Divider, Space, Typography } from "antd";
import { BookOpen as BookOpenIcon } from "lucide-react";

import { version } from "@/domain/version";

type VersionProps = {
  className?: string;
};

const Version = ({ className }: VersionProps) => {
  const { t } = useTranslation();

  return (
    <Space className={className} size={4}>
      <Typography.Link type="secondary" href="https://docs.certimate.me" target="_blank">
        <div className="flex items-center justify-center space-x-1">
          <BookOpenIcon size={16} />
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
