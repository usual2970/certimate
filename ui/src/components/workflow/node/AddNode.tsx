import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import {
  CloudUploadOutlined as CloudUploadOutlinedIcon,
  PlusOutlined as PlusOutlinedIcon,
  SendOutlined as SendOutlinedIcon,
  SisternodeOutlined as SisternodeOutlinedIcon,
  SolutionOutlined as SolutionOutlinedIcon,
} from "@ant-design/icons";
import { Dropdown } from "antd";

import { type WorkflowNode, WorkflowNodeType, newNode } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

export type AddNodeProps = {
  node: WorkflowNode;
  disabled?: boolean;
};

const AddNode = ({ node, disabled }: AddNodeProps) => {
  const { t } = useTranslation();

  const { addNode } = useWorkflowStore(useZustandShallowSelector(["addNode"]));

  const dropdownMenus = useMemo(() => {
    return [
      [WorkflowNodeType.Apply, "workflow_node.apply.label", <SolutionOutlinedIcon />],
      [WorkflowNodeType.Deploy, "workflow_node.deploy.label", <CloudUploadOutlinedIcon />],
      [WorkflowNodeType.Branch, "workflow_node.branch.label", <SisternodeOutlinedIcon />],
      [WorkflowNodeType.Notify, "workflow_node.notify.label", <SendOutlinedIcon />],
    ].map(([type, label, icon]) => {
      return {
        key: type as string,
        disabled: disabled,
        label: t(label as string),
        icon: icon,
        onClick: () => {
          if (disabled) return;

          const nextNode = newNode(type as WorkflowNodeType);
          addNode(nextNode, node.id);
        },
      };
    });
  }, []);

  return (
    <div className="relative py-6 before:absolute before:left-1/2 before:top-0 before:h-full before:w-[2px] before:-translate-x-1/2 before:bg-stone-200 before:content-['']">
      <Dropdown menu={{ items: dropdownMenus }} trigger={["click"]}>
        <div className="relative z-[1] flex size-5 cursor-pointer items-center justify-center rounded-full bg-stone-400 hover:bg-stone-500">
          <PlusOutlinedIcon className="text-white" />
        </div>
      </Dropdown>
    </div>
  );
};

export default AddNode;
