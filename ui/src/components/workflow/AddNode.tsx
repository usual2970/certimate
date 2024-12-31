import { useTranslation } from "react-i18next";
import { Dropdown } from "antd";
import { PlusOutlined as PlusOutlinedIcon } from "@ant-design/icons";

import { useZustandShallowSelector } from "@/hooks";
import { newWorkflowNode, workflowNodeDropdownList, WorkflowNodeType } from "@/domain/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { type BrandNodeProps, type NodeProps } from "./types";
import DropdownMenuItemIcon from "./DropdownMenuItemIcon";

const AddNode = ({ data }: NodeProps | BrandNodeProps) => {
  const { t } = useTranslation();

  const { addNode } = useWorkflowStore(useZustandShallowSelector(["addNode"]));

  const handleTypeSelected = (type: WorkflowNodeType, provider?: string) => {
    const node = newWorkflowNode(type, {
      providerType: provider,
    });

    addNode(node, data.id);
  };

  return (
    <div className="before:content-[''] before:w-[2px] before:bg-stone-200 before:absolute before:h-full before:left-[50%] before:-translate-x-[50%] before:top-0 py-6 relative">
      <Dropdown
        menu={{
          items: workflowNodeDropdownList.map((item) => {
            if (item.leaf) {
              return {
                key: item.type,
                label: <div className="ml-2">{item.name}</div>,
                icon: <DropdownMenuItemIcon type={item.icon.type} name={item.icon.name} />,
                onClick: () => {
                  handleTypeSelected(item.type);
                },
              };
            }

            return {
              key: item.type,
              label: <div className="ml-2">{item.name}</div>,
              icon: <DropdownMenuItemIcon type={item.icon.type} name={item.icon.name} />,
              children: item.children?.map((subItem) => {
                return {
                  key: subItem.providerType,
                  label: <div className="ml-2">{subItem.name}</div>,
                  icon: <DropdownMenuItemIcon type={subItem.icon.type} name={subItem.icon.name} />,
                  onClick: () => {
                    handleTypeSelected(item.type, subItem.providerType);
                  },
                };
              }),
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
