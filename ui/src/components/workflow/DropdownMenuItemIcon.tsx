import {
  CloudUploadOutlined as CloudUploadOutlinedIcon,
  SendOutlined as SendOutlinedIcon,
  SisternodeOutlined as SisternodeOutlinedIcon,
  SolutionOutlined as SolutionOutlinedIcon,
} from "@ant-design/icons";
import { Avatar } from "antd";

import { type WorkflowNodeDropdwonItemIcon, WorkflowNodeDropdwonItemIconType } from "@/domain/workflow";

const icons = new Map([
  ["ApplyNodeIcon", <SolutionOutlinedIcon />],
  ["DeployNodeIcon", <CloudUploadOutlinedIcon />],
  ["BranchNodeIcon", <SisternodeOutlinedIcon />],
  ["NotifyNodeIcon", <SendOutlinedIcon />],
]);

const DropdownMenuItemIcon = ({ type, name }: WorkflowNodeDropdwonItemIcon) => {
  const getIcon = () => {
    if (type === WorkflowNodeDropdwonItemIconType.Icon) {
      return icons.get(name);
    } else {
      return <Avatar src={name} size="small" />;
    }
  };

  return getIcon();
};

export default DropdownMenuItemIcon;
