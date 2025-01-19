import { memo } from "react";
import { useTranslation } from "react-i18next";
import { MoreOutlined as MoreOutlinedIcon } from "@ant-design/icons";
import { Button, Card, Popover } from "antd";

import { CheckCircleIcon, XCircleIcon } from "lucide-react";
import { WorkflowNodeType } from "@/domain/workflow";
import AddNode from "./AddNode";
import SharedNode, { type SharedNodeProps } from "./_SharedNode";

export type ConditionNodeProps = SharedNodeProps & {
  branchId: string;
  branchIndex: number;
};

const ExecuteResultNode = ({ node, disabled, branchId, branchIndex }: ConditionNodeProps) => {
  const { t } = useTranslation();

  return (
    <>
      <Popover
        arrow={false}
        content={
          <SharedNode.Menu
            node={node}
            branchId={branchId}
            branchIndex={branchIndex}
            disabled={disabled}
            trigger={<Button color="primary" icon={<MoreOutlinedIcon />} variant="text" />}
          />
        }
        overlayClassName="shadow-md"
        overlayInnerStyle={{ padding: 0 }}
        placement="rightTop"
      >
        <Card className="relative z-[1] mt-10 w-[256px] shadow-md" styles={{ body: { padding: 0 } }} hoverable>
          <div className="flex h-[48px] flex-col items-center justify-center truncate px-4 py-2">
            <div className="flex items-center space-x-2">
              {node.type === WorkflowNodeType.ExecuteSuccess ? (
                <>
                  <CheckCircleIcon size={18} className="text-green-500" />
                  <div>{t("workflow_node.execute_success.label")}</div>
                </>
              ) : (
                <>
                  <XCircleIcon size={18} className="text-red-500" />
                  <div>{t("workflow_node.execute_failure.label")}</div>
                </>
              )}
            </div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} disabled={disabled} />
    </>
  );
};

export default memo(ExecuteResultNode);
