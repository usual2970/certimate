import { useTranslation } from "react-i18next";
import { CloseCircleOutlined as CloseCircleOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { Button, Card, Dropdown, Popover } from "antd";
import { produce } from "immer";

import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import AddNode from "./AddNode";
import { type NodeProps } from "../types";

const ConditionNode = ({ node, branchId, branchIndex }: NodeProps) => {
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
                  label: t("workflow_node.action.delete_branch"),
                  icon: <CloseCircleOutlinedIcon />,
                  danger: true,
                  onClick: () => {
                    removeBranch(branchId ?? "", branchIndex ?? 0);
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
        <Card className="relative w-[256px] shadow-md mt-10 z-[1]" styles={{ body: { padding: 0 } }} hoverable>
          <div className="h-[48px] px-4 py-2 flex flex-col justify-center items-center truncate">
            <div
              className="w-full text-center outline-slate-200 overflow-hidden focus:bg-background focus:text-foreground focus:rounded-sm"
              contentEditable
              suppressContentEditableWarning
              onBlur={handleNodeNameBlur}
            >
              {node.name}
            </div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} />
    </>
  );
};

export default ConditionNode;
