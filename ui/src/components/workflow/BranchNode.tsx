import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "antd";

import { type WorkflowBranchNode, type WorkflowNode } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import AddNode from "./AddNode";
import NodeRender from "./NodeRender";
import { type BrandNodeProps } from "./types";

const BranchNode = memo(({ data }: BrandNodeProps) => {
  const { addBranch } = useWorkflowStore(useZustandShallowSelector(["addBranch"]));

  const { t } = useTranslation();

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
      <div className="border-t-[2px] border-b-[2px] relative flex gap-x-16 border-stone-200 bg-background">
        <Button
          size="small"
          variant="outlined"
          onClick={() => {
            addBranch(data.id);
          }}
          className="text-xs px-2 h-6 rounded-full absolute -top-3 left-[50%] -translate-x-1/2 z-[1] dark:text-stone-200"
        >
          {t("workflow.node.addBranch.label")}
        </Button>

        {data.branches.map((branch, index) => (
          <div
            key={branch.id}
            className="relative flex flex-col items-center before:content-['']  before:w-[2px] before:bg-stone-200 before:absolute before:h-full before:left-[50%] before:-translate-x-[50%] before:top-0"
          >
            {index == 0 && (
              <>
                <div className="w-[50%]  h-2 absolute -top-1 bg-background -left-[1px]"></div>
                <div className="w-[50%]  h-2 absolute -bottom-1 bg-background -left-[1px]"></div>
              </>
            )}
            {index == data.branches.length - 1 && (
              <>
                <div className="w-[50%]  h-2 absolute -top-1 bg-background -right-[1px]"></div>
                <div className="w-[50%]  h-2 absolute -bottom-1 bg-background -right-[1px]"></div>
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
