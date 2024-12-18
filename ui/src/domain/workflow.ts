import { produce } from "immer";
import { nanoid } from "nanoid";

import i18n from "@/i18n";
import { deployTargets, KVType } from "./domain";

export type WorkflowRunLog = {
  id: string;
  workflow: string;
  log: WorkflowRunLogItem[];
  error: string;
  succeed: boolean;
  created: string;
  updated: string;
};

export type WorkflowRunLogItem = {
  nodeName: string;
  error: string;
  outputs: WorkflowOutput[];
};

export type WorkflowOutput = {
  time: string;
  title: string;
  content: string;
  error: string;
};

export interface WorkflowModel extends BaseModel {
  name: string;
  description?: string;
  type: string;
  crontab?: string;
  content?: WorkflowNode;
  draft?: WorkflowNode;
  enabled?: boolean;
  hasDraft?: boolean;
}

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

const i18nPrefix = "workflow.node";

export const workflowNodeTypeDefaultName: Map<WorkflowNodeType, string> = new Map([
  [WorkflowNodeType.Start, i18n.t(`${i18nPrefix}.start.title`)],
  [WorkflowNodeType.End, i18n.t(`${i18nPrefix}.end.title`)],
  [WorkflowNodeType.Branch, i18n.t(`${i18nPrefix}.branch.title`)],
  [WorkflowNodeType.Condition, i18n.t(`${i18nPrefix}.condition.title`)],
  [WorkflowNodeType.Apply, i18n.t(`${i18nPrefix}.apply.title`)],
  [WorkflowNodeType.Deploy, i18n.t(`${i18nPrefix}.deploy.title`)],
  [WorkflowNodeType.Notify, i18n.t(`${i18nPrefix}.notify.title`)],
  [WorkflowNodeType.Custom, i18n.t(`${i18nPrefix}.custom.title`)],
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
        label: i18n.t("workflow.common.certificate.label"),
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
        label: i18n.t("workflow.common.certificate.label"),
      },
    ],
  ],
  [WorkflowNodeType.Deploy, []],
  [WorkflowNodeType.Notify, []],
]);

export type WorkflowNodeConfig = Record<string, string | boolean | number | KVType[] | string[] | undefined>;

export type WorkflowNode = {
  id: string;
  name: string;
  type: WorkflowNodeType;
  validated?: boolean;

  input?: WorkflowNodeIo[];
  config?: WorkflowNodeConfig;
  output?: WorkflowNodeIo[];

  next?: WorkflowNode | WorkflowBranchNode;
};

type NewWorkflowNodeOptions = {
  branchIndex?: number;
  providerType?: string;
};

export const initWorkflow = (): WorkflowModel => {
  // 开始节点
  const rs = newWorkflowNode(WorkflowNodeType.Start, {});
  let root = rs;

  // 申请节点
  root.next = newWorkflowNode(WorkflowNodeType.Apply, {});
  root = root.next;

  // 部署节点
  root.next = newWorkflowNode(WorkflowNodeType.Deploy, {});
  root = root.next;

  // 通知节点
  root.next = newWorkflowNode(WorkflowNodeType.Notify, {});

  return {
    id: "",
    name: i18n.t("workflow.props.name.default"),
    type: "auto",
    crontab: "0 0 * * *",
    enabled: false,
    draft: rs,
    created: new Date().toUTCString(),
    updated: new Date().toUTCString(),
  };
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

  if (type == WorkflowNodeType.Condition) {
    rs.validated = true;
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

export const allNodesValidated = (node: WorkflowNode | WorkflowBranchNode): boolean => {
  let current = node;
  while (current) {
    if (!isWorkflowBranchNode(current) && !current.validated) {
      return false;
    }
    if (isWorkflowBranchNode(current)) {
      for (const branch of current.branches) {
        if (!allNodesValidated(branch)) {
          return false;
        }
      }
    }
    current = current.next as WorkflowNode;
  }
  return true;
};

export const getExecuteMethod = (node: WorkflowNode): { type: string; crontab: string } => {
  if (node.type === WorkflowNodeType.Start) {
    return {
      type: (node.config?.executionMethod as string) ?? "",
      crontab: (node.config?.crontab as string) ?? "",
    };
  } else {
    return {
      type: "",
      crontab: "",
    };
  }
};

export type WorkflowBranchNode = {
  id: string;
  name: string;
  type: WorkflowNodeType;
  branches: WorkflowNode[];
  next?: WorkflowNode | WorkflowBranchNode;
};

type WorkflowNodeDropdwonItem = {
  type: WorkflowNodeType;
  providerType?: string;
  name: string;
  icon: WorkflowNodeDropdwonItemIcon;
  leaf?: boolean;
  children?: WorkflowNodeDropdwonItem[];
};

export enum WorkflowNodeDropdwonItemIconType {
  Icon,
  Provider,
}

export type WorkflowNodeDropdwonItemIcon = {
  type: WorkflowNodeDropdwonItemIconType;
  name: string;
};

const workflowNodeDropdownDeployList: WorkflowNodeDropdwonItem[] = deployTargets.map((item) => {
  return {
    type: WorkflowNodeType.Apply,
    providerType: item.type,
    name: i18n.t(item.name),
    leaf: true,
    icon: {
      type: WorkflowNodeDropdwonItemIconType.Provider,
      name: item.icon,
    },
  };
});

export const workflowNodeDropdownList: WorkflowNodeDropdwonItem[] = [
  {
    type: WorkflowNodeType.Apply,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Apply) ?? "",
    icon: {
      type: WorkflowNodeDropdwonItemIconType.Icon,
      name: "NotebookPen",
    },
    leaf: true,
  },
  {
    type: WorkflowNodeType.Deploy,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Deploy) ?? "",
    icon: {
      type: WorkflowNodeDropdwonItemIconType.Icon,
      name: "CloudUpload",
    },
    children: workflowNodeDropdownDeployList,
  },
  {
    type: WorkflowNodeType.Branch,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Branch) ?? "",
    leaf: true,
    icon: {
      type: WorkflowNodeDropdwonItemIconType.Icon,
      name: "GitFork",
    },
  },
  {
    type: WorkflowNodeType.Notify,
    name: workflowNodeTypeDefaultName.get(WorkflowNodeType.Notify) ?? "",
    leaf: true,
    icon: {
      type: WorkflowNodeDropdwonItemIconType.Icon,
      name: "Megaphone",
    },
  },
];
