import { Plus } from "lucide-react";

import { BrandNodeProps, NodeProps } from "./types";

import { newWorkflowNode, workflowNodeDropdownList, WorkflowNodeType } from "@/domain/workflow";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useShallow } from "zustand/shallow";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import DropdownMenuItemIcon from "./DropdownMenuItemIcon";
import Show from "../Show";

const selectState = (state: WorkflowState) => ({
  addNode: state.addNode,
});

const AddNode = ({ data }: NodeProps | BrandNodeProps) => {
  const { addNode } = useWorkflowStore(useShallow(selectState));

  const handleTypeSelected = (type: WorkflowNodeType, provider?: string) => {
    const node = newWorkflowNode(type, {
      providerType: provider,
    });

    addNode(node, data.id);
  };

  return (
    <div className="before:content-['']  before:w-[2px] before:bg-stone-300 before:absolute before:h-full before:left-[50%] before:-translate-x-[50%] before:top-0 pt-6 pb-9 relative flex flex-col items-center">
      <DropdownMenu>
        <DropdownMenuTrigger className="">
          <div className="bg-stone-400 hover:bg-stone-500 rounded-full z-10 relative outline-none">
            <Plus size={18} className="text-white" />
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuLabel>选择节点类型</DropdownMenuLabel>
          <DropdownMenuSeparator />
          {workflowNodeDropdownList.map((item) => (
            <Show
              key={item.type}
              when={!!item.leaf}
              fallback={
                <DropdownMenuSub>
                  <DropdownMenuSubTrigger className="flex space-x-2">
                    <DropdownMenuItemIcon type={item.icon.type} name={item.icon.name} /> <div>{item.name}</div>
                  </DropdownMenuSubTrigger>
                  <DropdownMenuPortal>
                    <DropdownMenuSubContent>
                      {item.children?.map((subItem) => {
                        return (
                          <DropdownMenuItem
                            key={subItem.providerType}
                            className="flex space-x-2"
                            onClick={() => {
                              handleTypeSelected(item.type, subItem.providerType);
                            }}
                          >
                            <DropdownMenuItemIcon type={subItem.icon.type} name={subItem.icon.name} /> <div>{subItem.name}</div>
                          </DropdownMenuItem>
                        );
                      })}
                    </DropdownMenuSubContent>
                  </DropdownMenuPortal>
                </DropdownMenuSub>
              }
            >
              <DropdownMenuItem
                key={item.type}
                className="flex space-x-2"
                onClick={() => {
                  handleTypeSelected(item.type, item.providerType);
                }}
              >
                <DropdownMenuItemIcon type={item.icon.type} name={item.icon.name} /> <div>{item.name}</div>
              </DropdownMenuItem>
            </Show>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

export default AddNode;
