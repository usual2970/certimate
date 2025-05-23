import { memo, useMemo } from "react";
import { useTranslation } from "react-i18next";
import {
  CloudUploadOutlined as CloudUploadOutlinedIcon,
  CopyOutlined as CopyOutlinedIcon,
  DeploymentUnitOutlined as DeploymentUnitOutlinedIcon,
  PlusOutlined as PlusOutlinedIcon,
  SendOutlined as SendOutlinedIcon,
  SisternodeOutlined as SisternodeOutlinedIcon,
  SolutionOutlined as SolutionOutlinedIcon,
} from "@ant-design/icons";
import { Dropdown } from "antd";

import { WorkflowNodeType, hasCloneNode, newNode } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import { type SharedNodeProps } from "./_SharedNode";

export type AddNodeProps = SharedNodeProps;

const AddNode = ({ node, disabled }: AddNodeProps) => {
  const { t } = useTranslation();

  const { addNode, workflow } = useWorkflowStore(useZustandShallowSelector(["addNode", "workflow"]));

  const cloning = hasCloneNode(workflow.draft!);

  const dropdownMenus = useMemo(() => {
    return [
      [WorkflowNodeType.Apply, "workflow_node.apply.label", <SolutionOutlinedIcon />],
      [WorkflowNodeType.Upload, "workflow_node.upload.label", <CloudUploadOutlinedIcon />],
      [WorkflowNodeType.Deploy, "workflow_node.deploy.label", <DeploymentUnitOutlinedIcon />],
      [WorkflowNodeType.Notify, "workflow_node.notify.label", <SendOutlinedIcon />],
      [WorkflowNodeType.Branch, "workflow_node.branch.label", <SisternodeOutlinedIcon />],
      [WorkflowNodeType.ExecuteResultBranch, "workflow_node.execute_result_branch.label", <SisternodeOutlinedIcon />],
      [WorkflowNodeType.Clone, "workflow_node.clone.label", <CopyOutlinedIcon />],
    ]
      .filter(([type]) => {
        if (node.type !== WorkflowNodeType.Apply && node.type !== WorkflowNodeType.Deploy && node.type !== WorkflowNodeType.Notify) {
          return type !== WorkflowNodeType.ExecuteResultBranch;
        }

        return true;
      })
      .map(([type, label, icon]) => {
        return {
          key: type as string,
          disabled: disabled,
          label: t(label as string),
          icon: icon,
          onClick: () => {
            const nextNode = newNode(type as WorkflowNodeType);
            addNode(nextNode, node.id);
          },
        };
      });
  }, [node.id, disabled, node.type]);

  const renderButton = () => {
    const buttonClassName =
      "relative z-[1] flex size-5 items-center justify-center rounded-full " +
      (cloning ? "bg-stone-300 cursor-not-allowed" : "bg-stone-400 cursor-pointer hover:bg-stone-500");

    return (
      <div className={buttonClassName}>
        <PlusOutlinedIcon className="text-white" />
      </div>
    );
  };

  return (
    <div className="relative py-6 before:absolute before:left-1/2 before:top-0 before:h-full before:w-[2px] before:-translate-x-1/2 before:bg-stone-200 before:content-['']">
      {cloning ? (
        <>{renderButton()}</>
      ) : (
        <Dropdown menu={{ items: dropdownMenus }} trigger={["click"]} disabled={disabled}>
          {renderButton()}
        </Dropdown>
      )}
    </div>
  );
};

export default memo(AddNode);
