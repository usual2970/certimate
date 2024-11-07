import { produce } from "immer";
import { nanoid } from "nanoid";
import { accessProviders } from "./access";
import i18n from "@/i18n";
import { deployTargets } from "./domain";

export enum WorkflowNodeType {
  Start = "start",
  End = "end",
  Branch = "branch",
  Condition = "condition",
  Apply = "apply",
  Deploy = "deploy",
  Notify = "notify",
  Custom = "custom",
}

export const workflowNodeTypeDefaultName: Map<WorkflowNodeType, string> = new Map([
  [WorkflowNodeType.Start, "开始"],
  [WorkflowNodeType.End, "结束"],
  [WorkflowNodeType.Branch, "分支"],
  [WorkflowNodeType.Condition, "分支"],
  [WorkflowNodeType.Apply, "申请"],
  [WorkflowNodeType.Deploy, "部署"],
  [WorkflowNodeType.Notify, "通知"],
  [WorkflowNodeType.Custom, "自定义"],
]);

export type WorkflowNodeIo = {
  name: string;
  type: string;
  required: boolean;
  label: string;
  value?: string;
  valueSelector?: WorkflowNodeIoValueSelector;
};

export type WorkflowNodeIoValueSelector = {
  id: string;
  name: string;
};

export const workflowNodeTypeDefaultInput: Map<WorkflowNodeType, WorkflowNodeIo[]> = new Map([
  [WorkflowNodeType.Apply, []],
  [
    WorkflowNodeType.Deploy,
    [
      {
        name: "certificate",
        type: " certificate",
        required: true,
        label: "证书",
      },
    ],
  ],
  [WorkflowNodeType.Notify, []],
]);

export const workflowNodeTypeDefaultOutput: Map<WorkflowNodeType, WorkflowNodeIo[]> = new Map([
  [
    WorkflowNodeType.Apply,
    [
      {
        name: "certificate",
        type: "certificate",
        required: true,
        label: "证书",
      },
    ],
  ],
  [WorkflowNodeType.Deploy, []],
  [WorkflowNodeType.Notify, []],
]);

export type WorkflowNodeConfig = Record<string, string | boolean | number | string[] | undefined>;

export type WorkflowNode = {
  id: string;
  name: string;
  type: WorkflowNodeType;

  input?: WorkflowNodeIo[];
  config?: WorkflowNodeConfig;
  configured?: boolean;
  output?: WorkflowNodeIo[];

  next?: WorkflowNode | WorkflowBranchNode;
};

type NewWorkflowNodeOptions = {
  branchIndex?: number;
  providerType?: string;
};

export const newWorkflowNode = (type: WorkflowNodeType, options: NewWorkflowNodeOptions): WorkflowNode | WorkflowBranchNode => {
  const id = nanoid();
  const typeName = workflowNodeTypeDefaultName.get(type) || "";
  const name = options.branchIndex !== undefined ? `${typeName} ${options.branchIndex + 1}` : typeName;

  let rs: WorkflowNode | WorkflowBranchNode = {
    id,
    name,
    type,
  };

  if (type === WorkflowNodeType.Apply || type === WorkflowNodeType.Deploy) {
    rs = {
      ...rs,
      config: {
        providerType: options.providerType,
      },
      input: workflowNodeTypeDefaultInput.get(type),
      output: workflowNodeTypeDefaultOutput.get(type),
    };
  }

  if (type === WorkflowNodeType.Branch) {
    rs = {
      ...rs,
      branches: [newWorkflowNode(WorkflowNodeType.Condition, { branchIndex: 0 }), newWorkflowNode(WorkflowNodeType.Condition, { branchIndex: 1 })],
    };
  }

  return rs;
};

export const isWorkflowBranchNode = (node: WorkflowNode | WorkflowBranchNode): node is WorkflowBranchNode => {
  return node.type === WorkflowNodeType.Branch;
};

export const updateNode = (node: WorkflowNode | WorkflowBranchNode, targetNode: WorkflowNode | WorkflowBranchNode) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.id === targetNode.id) {
        Object.assign(current, targetNode);
        break;
      }
      if (isWorkflowBranchNode(current)) {
        current.branches = current.branches.map((branch) => updateNode(branch, targetNode));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const addNode = (node: WorkflowNode | WorkflowBranchNode, preId: string, targetNode: WorkflowNode | WorkflowBranchNode) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.id === preId && !isWorkflowBranchNode(targetNode)) {
        targetNode.next = current.next;
        current.next = targetNode;
        break;
      } else if (current.id === preId && isWorkflowBranchNode(targetNode)) {
        targetNode.branches[0].next = current.next;
        current.next = targetNode;
        break;
      }
      if (isWorkflowBranchNode(current)) {
        current.branches = current.branches.map((branch) => addNode(branch, preId, targetNode));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const addBranch = (node: WorkflowNode | WorkflowBranchNode, branchNodeId: string) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.id === branchNodeId) {
        if (!isWorkflowBranchNode(current)) {
          return draft;
        }
        current.branches.push(
          newWorkflowNode(WorkflowNodeType.Condition, {
            branchIndex: current.branches.length,
          })
        );
        break;
      }
      if (isWorkflowBranchNode(current)) {
        current.branches = current.branches.map((branch) => addBranch(branch, branchNodeId));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const removeNode = (node: WorkflowNode | WorkflowBranchNode, targetNodeId: string) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.next?.id === targetNodeId) {
        current.next = current.next.next;
        break;
      }
      if (isWorkflowBranchNode(current)) {
        current.branches = current.branches.map((branch) => removeNode(branch, targetNodeId));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const removeBranch = (node: WorkflowNode | WorkflowBranchNode, branchNodeId: string, branchIndex: number) => {
  return produce(node, (draft) => {
    let current = draft;
    let last: WorkflowNode | WorkflowBranchNode | undefined = {
      id: "",
      name: "",
      type: WorkflowNodeType.Start,
      next: draft,
    };
    while (current && last) {
      if (current.id === branchNodeId) {
        if (!isWorkflowBranchNode(current)) {
          return draft;
        }
        current.branches.splice(branchIndex, 1);

        // 如果仅剩一个分支，删除分支节点，将分支节点的下一个节点挂载到当前节点
        if (current.branches.length === 1) {
          const branch = current.branches[0];
          if (branch.next) {
            last.next = branch.next;
            let lastNode: WorkflowNode | WorkflowBranchNode | undefined = branch.next;
            while (lastNode?.next) {
              lastNode = lastNode.next;
            }
            lastNode.next = current.next;
          } else {
            last.next = current.next;
          }
        }

        break;
      }
      if (isWorkflowBranchNode(current)) {
        current.branches = current.branches.map((branch) => removeBranch(branch, branchNodeId, branchIndex));
      }
      current = current.next as WorkflowNode;
      last = last.next;
    }
    return draft;
  });
};

// 1 个分支的节点，不应该能获取到相邻分支上节点的输出
export const getWorkflowOutputBeforeId = (node: WorkflowNode | WorkflowBranchNode, id: string, type: string): WorkflowNode[] => {
  const output: WorkflowNode[] = [];

  const traverse = (current: WorkflowNode | WorkflowBranchNode, output: WorkflowNode[]) => {
    if (!current) {
      return false;
    }
    if (current.id === id) {
      return true;
    }

    if (!isWorkflowBranchNode(current) && current.output && current.output.some((io) => io.type === type)) {
      output.push({
        ...current,
        output: current.output.filter((io) => io.type === type),
      });
    }

    if (isWorkflowBranchNode(current)) {
      const currentLength = output.length;
      console.log(currentLength);
      for (const branch of current.branches) {
        if (traverse(branch, output)) {
          return true;
        }
        // 如果当前分支没有输出，清空之前的输出
        if (output.length > currentLength) {
          output.splice(currentLength);
        }
      }
    }

    return traverse(current.next as WorkflowNode, output);
  };

  traverse(node, output);
  return output;
};

export type WorkflowBranchNode = {
  id: string;
  name: string;
  type: WorkflowNodeType;
  branches: WorkflowNode[];
  next?: WorkflowNode | WorkflowBranchNode;
};

type WorkflowwNodeDropdwonItem = {
  type: WorkflowNodeType;
  providerType?: string;
  name: string;
  icon: WorkflowwNodeDropdwonItemIcon;
  leaf?: boolean;
  children?: WorkflowwNodeDropdwonItem[];
};

export enum WorkflowwNodeDropdwonItemIconType {
  Icon,
  Provider,
}

export type WorkflowwNodeDropdwonItemIcon = {
  type: WorkflowwNodeDropdwonItemIconType;
  name: string;
};

const workflowNodeDropdownDeployList: WorkflowwNodeDropdwonItem[] = deployTargets.map((item) => {
  return {
    type: WorkflowNodeType.Apply,
    providerType: item.type,
    name: i18n.t(item.name),
    leaf: true,
    icon: {
      type: WorkflowwNodeDropdwonItemIconType.Provider,
      name: item.icon,
    },
  };
});

export const workflowNodeDropdownList: WorkflowwNodeDropdwonItem[] = [
  {
    type: WorkflowNodeType.Apply,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Apply) ?? "",
    icon: {
      type: WorkflowwNodeDropdwonItemIconType.Icon,
      name: "NotebookPen",
    },
    leaf: true,
  },
  {
    type: WorkflowNodeType.Deploy,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Deploy) ?? "",
    icon: {
      type: WorkflowwNodeDropdwonItemIconType.Icon,
      name: "CloudUpload",
    },
    children: workflowNodeDropdownDeployList,
  },
  {
    type: WorkflowNodeType.Branch,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Branch) ?? "",
    leaf: true,
    icon: {
      type: WorkflowwNodeDropdwonItemIconType.Icon,
      name: "GitFork",
    },
  },
  {
    type: WorkflowNodeType.Notify,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Notify) ?? "",
    leaf: true,
    icon: {
      type: WorkflowwNodeDropdwonItemIconType.Icon,
      name: "Megaphone",
    },
  },
];

