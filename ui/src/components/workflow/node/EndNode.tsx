import { useTranslation } from "react-i18next";
import { Typography } from "antd";

const EndNode = () => {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col items-center">
      <div className="size-[20px] rounded-full bg-stone-400"></div>
      <div className="text-sm mt-2">
        <Typography.Text type="secondary">{t("workflow_node.end.label")}</Typography.Text>
      </div>
    </div>
  );
};

export default EndNode;
