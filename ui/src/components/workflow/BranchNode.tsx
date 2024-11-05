import { Button } from "@/components/ui/button";
import AddNode from "./AddNode";
import { WorkflowBranchNode, WorkflowNode } from "@/domain/workflow";
import NodeRender from "./NodeRender";
import { memo } from "react";
import { BrandNodeProps } from "./types";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useShallow } from "zustand/shallow";

const selectState = (state: WorkflowState) => ({
  addBranch: state.addBranch,
});
const BranchNode = memo(({ data }: BrandNodeProps) => {
  const { addBranch } = useWorkflowStore(useShallow(selectState));

  const renderNodes = (node: WorkflowBranchNode | WorkflowNode | undefined, branchNodeId?: string, branchIndex?: number) => {
    const elements: JSX.Element[] = [];
    let current = node;
    while (current) {
      elements.push(<NodeRender data={current} branchId={branchNodeId} branchIndex={branchIndex} key={current.id} />);
      current = current.next;
    }
    return elements;
  };

  return (
    <>
      <div className="border-t-[2px] border-b-[2px] relative flex gap-x-16 border-stone-300 bg-slate-50">
        <Button
          onClick={() => {
            addBranch(data.id);
          }}
          size={"sm"}
          variant={"outline"}
          className="text-xs px-2 h-6 rounded-full absolute  -top-3 left-[50%] -translate-x-1/2 z-10"
        >
          添加分支
        </Button>

        {data.branches.map((branch, index) => (
          <div
            key={branch.id}
            className="relative flex flex-col items-center before:content-['']  before:w-[2px] before:bg-stone-300 before:absolute before:h-full before:left-[50%] before:-translate-x-[50%] before:top-0"
          >
            {index == 0 && (
              <>
                <div className="w-[50%]  h-2 absolute -top-1 bg-stone-50 -left-[1px]"></div>
                <div className="w-[50%]  h-2 absolute -bottom-1 bg-stone-50 -left-[1px]"></div>
              </>
            )}
            {index == data.branches.length - 1 && (
              <>
                <div className="w-[50%]  h-2 absolute -top-1 bg-stone-50 -right-[1px]"></div>
                <div className="w-[50%]  h-2 absolute -bottom-1 bg-stone-50 -right-[1px]"></div>
              </>
            )}
            {/* 条件 1 */}
            <div className="relative flex flex-col items-center">{renderNodes(branch, data.id, index)}</div>
          </div>
        ))}
      </div>
      <AddNode data={data} />
    </>
  );
});

export default BranchNode;
