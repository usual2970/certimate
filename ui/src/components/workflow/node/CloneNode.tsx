import { Card } from "antd";
import { type SharedNodeProps } from "./_SharedNode";
import AddNode from "./AddNode";
import { useTranslation } from "react-i18next";

export type UploadNodeProps = SharedNodeProps;
const CloneNode = ({ node, disabled }: SharedNodeProps) => {
  const { t } = useTranslation();
  return (
    <>
      <Card className="relative z-[1] w-[256px] shadow-md" styles={{ body: { padding: 0 } }} hoverable>
        <div className="flex h-[64px] flex-col items-center justify-center truncate px-4 py-2">{t("workflow_node.clone.description")}</div>
      </Card>
      <AddNode node={node} disabled={disabled} />
    </>
  );
};

export default CloneNode;
