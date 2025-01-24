import { memo } from "react";
import { MoreOutlined as MoreOutlinedIcon } from "@ant-design/icons";
import { Button, Card, Popover } from "antd";

import SharedNode, { type SharedNodeProps } from "./_SharedNode";
import AddNode from "./AddNode";

export type ConditionNodeProps = SharedNodeProps & {
  branchId: string;
  branchIndex: number;
};

const ConditionNode = ({ node, disabled, branchId, branchIndex }: ConditionNodeProps) => {
  // TODO: 条件分支

  return (
    <>
      <Popover
        classNames={{ root: "shadow-md" }}
        styles={{ body: { padding: 0 } }}
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
        placement="rightTop"
      >
        <Card className="relative z-[1] mt-10 w-[256px] shadow-md" styles={{ body: { padding: 0 } }} hoverable>
          <div className="flex h-[48px] flex-col items-center justify-center truncate px-4 py-2">
            <SharedNode.Title
              className="focus:bg-background focus:text-foreground overflow-hidden outline-slate-200 focus:rounded-sm"
              node={node}
              disabled={disabled}
            />
          </div>
        </Card>
      </Popover>

      <AddNode node={node} disabled={disabled} />
    </>
  );
};

export default memo(ConditionNode);
