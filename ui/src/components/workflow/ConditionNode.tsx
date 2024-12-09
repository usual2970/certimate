import { useWorkflowStore, WorkflowState } from "@/stores/workflow";
import AddNode from "./AddNode";
import { NodeProps } from "./types";
import { useShallow } from "zustand/shallow";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";
import { Ellipsis, Trash2 } from "lucide-react";

const selectState = (state: WorkflowState) => ({
  updateNode: state.updateNode,
  removeBranch: state.removeBranch,
});
const ConditionNode = ({ data, branchId, branchIndex }: NodeProps) => {
  const { updateNode, removeBranch } = useWorkflowStore(useShallow(selectState));
  const handleNameBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    updateNode({ ...data, name: e.target.innerText });
  };
  return (
    <>
      <div className="rounded-md shadow-md w-[261px] mt-10 relative z-10">
        <DropdownMenu>
          <DropdownMenuTrigger className="absolute right-2 top-1">
            <Ellipsis size={17} className="text-stone-600" />
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem
              className="flex space-x-2 text-red-600"
              onClick={() => {
                removeBranch(branchId ?? "", branchIndex ?? 0);
              }}
            >
              <Trash2 size={16} /> <div>删除分支</div>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

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
