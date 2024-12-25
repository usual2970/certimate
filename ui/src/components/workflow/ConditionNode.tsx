import { useWorkflowStore } from "@/stores/workflow";
import AddNode from "./AddNode";
import { NodeProps } from "./types";
import { useZustandShallowSelector } from "@/hooks";
import { Dropdown } from "antd";
import { Ellipsis, Trash2 } from "lucide-react";

const ConditionNode = ({ data, branchId, branchIndex }: NodeProps) => {
  const { updateNode, removeBranch } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeBranch"]));
  const handleNameBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    updateNode({ ...data, name: e.target.innerText });
  };
  return (
    <>
      <div className="rounded-md shadow-md w-[261px] mt-10 relative z-10">
        <Dropdown
          menu={{
            items: [
              {
                key: "delete",
                label: "删除分支",
                icon: <Trash2 size={16} />,
                danger: true,
                onClick: () => {
                  removeBranch(branchId ?? "", branchIndex ?? 0);
                },
              },
            ],
          }}
          trigger={["click"]}
        >
          <div className="absolute right-2 top-1 cursor-pointer">
            <Ellipsis size={17} className="text-stone-600" />
          </div>
        </Dropdown>

        <div className="w-[261px]  flex flex-col justify-center text-foreground rounded-md bg-white px-5 py-5">
          <div contentEditable suppressContentEditableWarning onBlur={handleNameBlur} className="text-center outline-slate-200 dark:text-stone-600">
            {data.name}
          </div>
        </div>
      </div>
      <AddNode data={data} />
    </>
  );
};

export default ConditionNode;
