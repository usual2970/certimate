import dayjs from "dayjs";
import { produce } from "immer";
import { nanoid } from "nanoid";

import i18n from "@/i18n";

export interface WorkflowModel extends BaseModel {
  name: string;
  description?: string;
  trigger: string;
  triggerCron?: string;
  enabled?: boolean;
  content?: WorkflowNode;
  draft?: WorkflowNode;
  hasDraft?: boolean;
  lastRunId?: string;
  lastRunStatus?: string;
  lastRunTime?: string;
}

export const WORKFLOW_TRIGGERS = Object.freeze({
  AUTO: "auto",
  MANUAL: "manual",
} as const);

export type WorkflowTriggerType = (typeof WORKFLOW_TRIGGERS)[keyof typeof WORKFLOW_TRIGGERS];

// #region Node
export enum WorkflowNodeType {
  Start = "start",
  End = "end",
  Branch = "branch",
  ExecuteResultBranch = "execute_result_branch",
  ExecuteSuccess = "execute_success",
  ExecuteFailure = "execute_failure",
  Condition = "condition",
  Apply = "apply",
  Upload = "upload",
  Deploy = "deploy",
  Notify = "notify",
  Custom = "custom",
}

const workflowNodeTypeDefaultNames: Map<WorkflowNodeType, string> = new Map([
  [WorkflowNodeType.Start, i18n.t("workflow_node.start.label")],
  [WorkflowNodeType.End, i18n.t("workflow_node.end.label")],
  [WorkflowNodeType.Branch, i18n.t("workflow_node.branch.label")],
  [WorkflowNodeType.ExecuteResultBranch, i18n.t("workflow_node.execute_result_branch.label")],
  [WorkflowNodeType.ExecuteSuccess, i18n.t("workflow_node.execute_success.label")],
  [WorkflowNodeType.ExecuteFailure, i18n.t("workflow_node.execute_failure.label")],
  [WorkflowNodeType.Condition, i18n.t("workflow_node.condition.label")],
  [WorkflowNodeType.Apply, i18n.t("workflow_node.apply.label")],
  [WorkflowNodeType.Upload, i18n.t("workflow_node.upload.label")],
  [WorkflowNodeType.Deploy, i18n.t("workflow_node.deploy.label")],
  [WorkflowNodeType.Notify, i18n.t("workflow_node.notify.label")],
  [WorkflowNodeType.Custom, i18n.t("workflow_node.custom.title")],
]);

const workflowNodeTypeDefaultInputs: Map<WorkflowNodeType, WorkflowNodeIO[]> = new Map([
  [WorkflowNodeType.Apply, []],
  [
    WorkflowNodeType.Deploy,
    [
      {
        name: "certificate",
        type: "certificate",
        required: true,
        label: "证书",
      },
    ],
  ],
  [WorkflowNodeType.Notify, []],
]);

const workflowNodeTypeDefaultOutputs: Map<WorkflowNodeType, WorkflowNodeIO[]> = new Map([
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
  [
    WorkflowNodeType.Upload,
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

export type WorkflowNode = {
  id: string;
  name: string;
  type: WorkflowNodeType;

  config?: Record<string, unknown>;
  inputs?: WorkflowNodeIO[];
  outputs?: WorkflowNodeIO[];

  next?: WorkflowNode;
  branches?: WorkflowNode[];

  validated?: boolean;
};

export type WorkflowNodeConfigForStart = {
  trigger: string;
  triggerCron?: string;
};

export type WorkflowNodeConfigForApply = {
  domains: string;
  contactEmail: string;
  provider: string;
  providerAccessId: string;
  providerConfig?: Record<string, unknown>;
  keyAlgorithm: string;
  nameservers?: string;
  dnsPropagationTimeout?: number;
  dnsTTL?: number;
  disableFollowCNAME?: boolean;
  disableARI?: boolean;
  skipBeforeExpiryDays: number;
};

export type WorkflowNodeConfigForUpload = {
  certificateId: string;
  domains: string;
  certificate: string;
  privateKey: string;
};

export type WorkflowNodeConfigForDeploy = {
  certificate: string;
  provider: string;
  providerAccessId: string;
  providerConfig: Record<string, unknown>;
  skipOnLastSucceeded: boolean;
};

export type WorkflowNodeConfigForNotify = {
  channel: string;
  subject: string;
  message: string;
};

export type WorkflowNodeConfigForBranch = never;

export type WorkflowNodeConfigForEnd = never;

export type WorkflowNodeIO = {
  name: string;
  type: string;
  required: boolean;
  label: string;
  value?: string;
  valueSelector?: WorkflowNodeIOValueSelector;
};

export type WorkflowNodeIOValueSelector = {
  id: string;
  name: string;
};
// #endregion

type InitWorkflowOptions = {
  template?: "standard";
};

export const initWorkflow = (options: InitWorkflowOptions = {}): WorkflowModel => {
  const root = newNode(WorkflowNodeType.Start, {}) as WorkflowNode;
  root.config = { trigger: WORKFLOW_TRIGGERS.MANUAL };

  if (options.template === "standard") {
    let current = root;
    current.next = newNode(WorkflowNodeType.Apply, {});

    current = current.next;
    current.next = newNode(WorkflowNodeType.Deploy, {});

    current = current.next;
    current.next = newNode(WorkflowNodeType.Notify, {});
  }

  return {
    id: null!,
    name: `MyWorkflow-${dayjs().format("YYYYMMDDHHmmss")}`,
    trigger: root.config!.trigger as string,
    triggerCron: root.config!.triggerCron as string,
    enabled: false,
    draft: root,
    hasDraft: true,
    created: new Date().toISOString(),
    updated: new Date().toISOString(),
  };
};

type NewNodeOptions = {
  branchIndex?: number;
};

export const newNode = (nodeType: WorkflowNodeType, options: NewNodeOptions = {}): WorkflowNode => {
  const nodeTypeName = workflowNodeTypeDefaultNames.get(nodeType) || "";
  const nodeName = options.branchIndex != null ? `${nodeTypeName} ${options.branchIndex + 1}` : nodeTypeName;

  const node: WorkflowNode = {
    id: nanoid(),
    name: nodeName,
    type: nodeType,
  };

  switch (nodeType) {
    case WorkflowNodeType.Apply:
    case WorkflowNodeType.Upload:
    case WorkflowNodeType.Deploy:
      {
        node.inputs = workflowNodeTypeDefaultInputs.get(nodeType);
        node.outputs = workflowNodeTypeDefaultOutputs.get(nodeType);
      }
      break;

    case WorkflowNodeType.Condition:
      {
        node.validated = true;
      }
      break;

    case WorkflowNodeType.Branch:
      {
        node.branches = [newNode(WorkflowNodeType.Condition, { branchIndex: 0 }), newNode(WorkflowNodeType.Condition, { branchIndex: 1 })];
      }
      break;

    case WorkflowNodeType.ExecuteResultBranch:
      {
        node.branches = [newNode(WorkflowNodeType.ExecuteSuccess), newNode(WorkflowNodeType.ExecuteFailure)];
      }
      break;

    case WorkflowNodeType.ExecuteSuccess:
    case WorkflowNodeType.ExecuteFailure:
      {
        node.validated = true;
      }
      break;
  }

  return node;
};

export const updateNode = (node: WorkflowNode, targetNode: WorkflowNode) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.id === targetNode.id) {
        Object.assign(current, targetNode);
        break;
      }
      if (current.type === WorkflowNodeType.Branch || current.type === WorkflowNodeType.ExecuteResultBranch) {
        current.branches = current.branches!.map((branch) => updateNode(branch, targetNode));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const addNode = (node: WorkflowNode, preId: string, targetNode: WorkflowNode) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.id === preId && targetNode.type !== WorkflowNodeType.Branch && targetNode.type !== WorkflowNodeType.ExecuteResultBranch) {
        targetNode.next = current.next;
        current.next = targetNode;
        break;
      } else if (current.id === preId && (targetNode.type === WorkflowNodeType.Branch || targetNode.type === WorkflowNodeType.ExecuteResultBranch)) {
        targetNode.branches![0].next = current.next;
        current.next = targetNode;
        break;
      }
      if (current.type === WorkflowNodeType.Branch || current.type === WorkflowNodeType.ExecuteResultBranch) {
        current.branches = current.branches!.map((branch) => addNode(branch, preId, targetNode));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const addBranch = (node: WorkflowNode, branchNodeId: string) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.id === branchNodeId) {
        if (current.type !== WorkflowNodeType.Branch) {
          return draft;
        }
        current.branches!.push(
          newNode(WorkflowNodeType.Condition, {
            branchIndex: current.branches!.length,
          })
        );
        break;
      }
      if (current.type === WorkflowNodeType.Branch || current.type === WorkflowNodeType.ExecuteResultBranch) {
        current.branches = current.branches!.map((branch) => addBranch(branch, branchNodeId));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const removeNode = (node: WorkflowNode, targetNodeId: string) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.next?.id === targetNodeId) {
        current.next = current.next.next;
        break;
      }
      if (current.type === WorkflowNodeType.Branch || current.type === WorkflowNodeType.ExecuteResultBranch) {
        current.branches = current.branches!.map((branch) => removeNode(branch, targetNodeId));
      }
      current = current.next as WorkflowNode;
    }
    return draft;
  });
};

export const removeBranch = (node: WorkflowNode, branchNodeId: string, branchIndex: number) => {
  return produce(node, (draft) => {
    let current = draft;
    let last: WorkflowNode | undefined = {
      id: "",
      name: "",
      type: WorkflowNodeType.Start,
      next: draft,
    };
    while (current && last) {
      if (current.id === branchNodeId) {
        if (current.type !== WorkflowNodeType.Branch && current.type !== WorkflowNodeType.ExecuteResultBranch) {
          return draft;
        }
        current.branches!.splice(branchIndex, 1);

        // 如果仅剩一个分支，删除分支节点，将分支节点的下一个节点挂载到当前节点
        if (current.branches!.length === 1) {
          const branch = current.branches![0];
          if (branch.next) {
            last.next = branch.next;
            let lastNode: WorkflowNode | undefined = branch.next;
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
      if (current.type === WorkflowNodeType.Branch || current.type === WorkflowNodeType.ExecuteResultBranch) {
        current.branches = current.branches!.map((branch) => removeBranch(branch, branchNodeId, branchIndex));
      }
      current = current.next as WorkflowNode;
      last = last.next;
    }
    return draft;
  });
};

// 1 个分支的节点，不应该能获取到相邻分支上节点的输出
export const getWorkflowOutputBeforeId = (node: WorkflowNode, id: string, type: string): WorkflowNode[] => {
  const output: WorkflowNode[] = [];

  const traverse = (current: WorkflowNode, output: WorkflowNode[]) => {
    if (!current) {
      return false;
    }
    if (current.id === id) {
      return true;
    }

    // 如果当前节点是execute_failure,清除execute_result_branch节点前一个节点的输出
    if (current.type === WorkflowNodeType.ExecuteFailure) {
      output.splice(output.length - 1);
    }

    if (current.type !== WorkflowNodeType.Branch && current.outputs && current.outputs.some((io) => io.type === type)) {
      output.push({
        ...current,
        outputs: current.outputs.filter((io) => io.type === type),
      });
    }

    if (current.type === WorkflowNodeType.Branch || current.type === WorkflowNodeType.ExecuteResultBranch) {
      const currentLength = output.length;
      for (const branch of current.branches!) {
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

export const isAllNodesValidated = (node: WorkflowNode): boolean => {
  let current = node as typeof node | undefined;
  while (current) {
    if (current.type === WorkflowNodeType.Branch || current.type === WorkflowNodeType.ExecuteResultBranch) {
      for (const branch of current.branches!) {
        if (!isAllNodesValidated(branch)) {
          return false;
        }
      }
    } else {
      if (!current.validated) {
        return false;
      }
    }

    current = current.next;
  }

  return true;
};

/**
 * @deprecated
 */
export const getExecuteMethod = (node: WorkflowNode): { trigger: string; triggerCron: string } => {
  if (node.type === WorkflowNodeType.Start) {
    const config = node.config as WorkflowNodeConfigForStart;
    return {
      trigger: config.trigger ?? "",
      triggerCron: config.triggerCron ?? "",
    };
  } else {
    return {
      trigger: "",
      triggerCron: "",
    };
  }
};
