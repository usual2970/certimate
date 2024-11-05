import { WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import AddNode from "./AddNode";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useShallow } from "zustand/shallow";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";
import { Ellipsis, Trash2 } from "lucide-react";
import { usePanel } from "./PanelProvider";
import PanelBody from "./PanelBody";

type NodeProps = {
  data: WorkflowNode;
};

const selectState = (state: WorkflowState) => ({
  updateNode: state.updateNode,
  removeNode: state.removeNode,
});
const Node = ({ data }: NodeProps) => {
  const { updateNode, removeNode } = useWorkflowStore(useShallow(selectState));
  const handleNameBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    updateNode({ ...data, name: e.target.innerText });
  };

  const { showPanel } = usePanel();

  const handleNodeSettingClick = () => {
    showPanel({
      name: data.name,
      children: <PanelBody data={data} />,
    });
  };
  return (
    <>
      <div className="rounded-md shadow-md w-[260px] relative">
        {data.type != WorkflowNodeType.Start && (
          <>
            <DropdownMenu>
              <DropdownMenuTrigger className="absolute right-2 top-1">
                <Ellipsis className="text-white" size={17} />
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem
                  className="flex space-x-2 text-red-600"
                  onClick={() => {
                    removeNode(data.id);
                  }}
                >
                  <Trash2 size={16} /> <div>删除节点</div>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </>
        )}

        <div className="w-[260px] h-[60px] flex flex-col justify-center items-center bg-primary text-white rounded-t-md px-5">
          <div
            contentEditable
            suppressContentEditableWarning
            onBlur={handleNameBlur}
            className="w-full text-center outline-none focus:bg-white focus:text-stone-600 focus:rounded-sm"
          >
            {data.name}
          </div>
        </div>
        <div className="p-2 text-sm text-primary flex flex-col justify-center bg-white">
          <div className="leading-7 text-primary cursor-pointer" onClick={handleNodeSettingClick}>
            设置节点
          </div>
        </div>
      </div>
      <AddNode data={data} />
    </>
  );
};

export default Node;
