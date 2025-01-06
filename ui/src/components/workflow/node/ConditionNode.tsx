import { memo } from "react";
import { useTranslation } from "react-i18next";
import { CloseCircleOutlined as CloseCircleOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { Button, Card, Dropdown, Popover } from "antd";
import { produce } from "immer";

import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import AddNode from "./AddNode";
import { type SharedNodeProps } from "./_SharedNode";

export type ConditionNodeProps = SharedNodeProps & {
  branchId: string;
  branchIndex: number;
};

const ConditionNode = ({ node, disabled, branchId, branchIndex }: ConditionNodeProps) => {
  const { t } = useTranslation();

  const { updateNode, removeBranch } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeBranch"]));

  const handleNodeNameBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    const oldName = node.name;
    const newName = e.target.innerText.trim();
    if (oldName === newName) {
      return;
    }

    updateNode(
      produce(node, (draft) => {
        draft.name = newName;
      })
    );
  };

  return (
    <>
      <Popover
        arrow={false}
        content={
          <Dropdown
            menu={{
              items: [
                {
                  key: "delete",
                  disabled: disabled,
                  label: t("workflow_node.action.delete_branch"),
                  icon: <CloseCircleOutlinedIcon />,
                  danger: true,
                  onClick: () => {
                    if (disabled) return;

                    removeBranch(branchId!, branchIndex!);
                  },
                },
              ],
            }}
            trigger={["click"]}
          >
            <Button color="primary" icon={<EllipsisOutlinedIcon />} variant="text" />
          </Dropdown>
        }
        overlayClassName="shadow-md"
        overlayInnerStyle={{ padding: 0 }}
        placement="rightTop"
      >
        <Card className="relative z-[1] mt-10 w-[256px] shadow-md" styles={{ body: { padding: 0 } }} hoverable>
          <div className="flex h-[48px] flex-col items-center justify-center truncate px-4 py-2">
            <div
              className="focus:bg-background focus:text-foreground w-full overflow-hidden text-center outline-slate-200 focus:rounded-sm"
              contentEditable
              suppressContentEditableWarning
              onBlur={handleNodeNameBlur}
            >
              {node.name}
            </div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} disabled={disabled} />
    </>
  );
};

export default memo(ConditionNode);
