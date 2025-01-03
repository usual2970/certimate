import {
  CloudUploadOutlined as CloudUploadOutlinedIcon,
  PlusOutlined as PlusOutlinedIcon,
  SendOutlined as SendOutlinedIcon,
  SisternodeOutlined as SisternodeOutlinedIcon,
  SolutionOutlined as SolutionOutlinedIcon,
} from "@ant-design/icons";
import { Dropdown } from "antd";

import { WorkflowNodeType, newNode, workflowNodeTypeDefaultNames } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import { type BrandNodeProps, type NodeProps } from "../types";

const dropdownMenus = [
  {
    type: WorkflowNodeType.Apply,
    label: workflowNodeTypeDefaultNames.get(WorkflowNodeType.Apply),
    icon: <SolutionOutlinedIcon />,
  },
  {
    type: WorkflowNodeType.Deploy,
    label: workflowNodeTypeDefaultNames.get(WorkflowNodeType.Deploy),
    icon: <CloudUploadOutlinedIcon />,
  },
  {
    type: WorkflowNodeType.Branch,
    label: workflowNodeTypeDefaultNames.get(WorkflowNodeType.Branch),
    icon: <SisternodeOutlinedIcon />,
  },
  {
    type: WorkflowNodeType.Notify,
    label: workflowNodeTypeDefaultNames.get(WorkflowNodeType.Notify),
    icon: <SendOutlinedIcon />,
  },
];

const AddNode = ({ node: supnode }: NodeProps | BrandNodeProps) => {
  const { addNode } = useWorkflowStore(useZustandShallowSelector(["addNode"]));

  const handleNodeTypeSelect = (type: WorkflowNodeType) => {
    const node = newNode(type);
    addNode(node, supnode.id);
  };

  return (
    <div className="relative py-6 before:content-[''] before:absolute before:w-[2px] before:h-full before:left-[50%] before:-translate-x-[50%] before:top-0 before:bg-stone-200">
      <Dropdown
        menu={{
          items: dropdownMenus.map((item) => {
            return {
              key: item.type,
              label: item.label,
              icon: item.icon,
              onClick: () => {
                handleNodeTypeSelect(item.type);
              },
            };
          }),
        }}
        trigger={["click"]}
      >
        <div className="bg-stone-400 hover:bg-stone-500 rounded-full size-5 z-[1] relative flex items-center justify-center cursor-pointer">
          <PlusOutlinedIcon className="text-white" />
        </div>
      </Dropdown>
    </div>
  );
};

export default AddNode;
