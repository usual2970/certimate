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
  Apply = "apply",
  Upload = "upload",
  Inspect = "inspect",
  Deploy = "deploy",
  Notify = "notify",
  Branch = "branch",
  Condition = "condition",
  ExecuteResultBranch = "execute_result_branch",
  ExecuteSuccess = "execute_success",
  ExecuteFailure = "execute_failure",
  Custom = "custom",
}

const workflowNodeTypeDefaultNames: Map<WorkflowNodeType, string> = new Map([
  [WorkflowNodeType.Start, i18n.t("workflow_node.start.label")],
  [WorkflowNodeType.End, i18n.t("workflow_node.end.label")],
  [WorkflowNodeType.Apply, i18n.t("workflow_node.apply.label")],
  [WorkflowNodeType.Upload, i18n.t("workflow_node.upload.label")],
  [WorkflowNodeType.Inspect, i18n.t("workflow_node.inspect.label")],
  [WorkflowNodeType.Deploy, i18n.t("workflow_node.deploy.label")],
  [WorkflowNodeType.Notify, i18n.t("workflow_node.notify.label")],
  [WorkflowNodeType.Branch, i18n.t("workflow_node.branch.label")],
  [WorkflowNodeType.Condition, i18n.t("workflow_node.condition.label")],
  [WorkflowNodeType.ExecuteResultBranch, i18n.t("workflow_node.execute_result_branch.label")],
  [WorkflowNodeType.ExecuteSuccess, i18n.t("workflow_node.execute_success.label")],
  [WorkflowNodeType.ExecuteFailure, i18n.t("workflow_node.execute_failure.label")],
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
        label: i18n.t("workflow.variables.certificate.label"),
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
        label: i18n.t("workflow.variables.certificate.label"),
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
        label: i18n.t("workflow.variables.certificate.label"),
      },
    ],
  ],
  [
    WorkflowNodeType.Inspect,
    [
      {
        name: "certificate",
        type: "certificate",
        required: true,
        label: i18n.t("workflow.variables.certificate.label"),
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
  challengeType: string;
  provider: string;
  providerAccessId: string;
  providerConfig?: Record<string, unknown>;
  caProvider?: string;
  caProviderAccessId?: string;
  caProviderConfig?: Record<string, unknown>;
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

export type WorkflowNodeConfigForInspect = {
  domain: string;
  port: string;
  host: string;
  path: string;
};

export type WorkflowNodeConfigForDeploy = {
  certificate: string;
  provider: string;
  providerAccessId?: string;
  providerConfig?: Record<string, unknown>;
  skipOnLastSucceeded: boolean;
};

export type WorkflowNodeConfigForNotify = {
  subject: string;
  message: string;
  /**
   * @deprecated
   */
  channel?: string;
  provider: string;
  providerAccessId: string;
  providerConfig?: Record<string, unknown>;
};

export type WorkflowNodeConfigForCondition = {
  expression: Expr;
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

export const VALUE_TYPES = Object.freeze({
  STRING: "string",
  NUMBER: "number",
  BOOLEAN: "boolean",
} as const);

export type WorkflowNodeIoValueType = (typeof VALUE_TYPES)[keyof typeof VALUE_TYPES];

export type WorkflowNodeIOValueSelector = {
  id: string;
  name: string;
  type: WorkflowNodeIoValueType;
};

type WorkflowNodeIOOptions = {
  label: string;
  value: string;
};

export const workflowNodeIOOptions = (node: WorkflowNode) => {
  const rs = {
    label: node.name,
    options: Array<WorkflowNodeIOOptions>(),
  };

  if (node.outputs) {
    for (const output of node.outputs) {
      switch (output.type) {
        case "certificate":
          rs.options.push({
            label: `${node.name} - ${output.label} - ${i18n.t("workflow.variables.is_validated.label")}`,
            value: `${node.id}#${output.name}.validated#boolean`,
          });

          rs.options.push({
            label: `${node.name} - ${output.label} - ${i18n.t("workflow.variables.days_left.label")}`,
            value: `${node.id}#${output.name}.daysLeft#number`,
          });
          break;
        default:
          rs.options.push({
            label: `${node.name} - ${output.label}`,
            value: `${node.id}#${output.name}#${output.type}`,
          });
          break;
      }
    }
  }

  return rs;
};

// #endregion

// #region Condition expression

export type Value = string | number | boolean;

export type ComparisonOperator = ">" | "<" | ">=" | "<=" | "==" | "!=" | "is";

export enum LogicalOperator {
  And = "and",
  Or = "or",
  Not = "not",
}

export enum ExprType {
  Const = "const",
  Var = "var",
  Compare = "compare",
  Logical = "logical",
  Not = "not",
}

export type ConstExpr = { type: ExprType.Const; value: string; valueType: WorkflowNodeIoValueType };
export type VarExpr = { type: ExprType.Var; selector: WorkflowNodeIOValueSelector };
export type CompareExpr = { type: ExprType.Compare; op: ComparisonOperator; left: Expr; right: Expr };
export type LogicalExpr = { type: ExprType.Logical; op: LogicalOperator; left: Expr; right: Expr };
export type NotExpr = { type: ExprType.Not; expr: Expr };

export type Expr = ConstExpr | VarExpr | CompareExpr | LogicalExpr | NotExpr;

export const isConstExpr = (expr: Expr): expr is ConstExpr => {
  return expr.type === ExprType.Const;
};

export const isVarExpr = (expr: Expr): expr is VarExpr => {
  return expr.type === ExprType.Var;
};

// #endregion

const isBranchLike = (node: WorkflowNode) => {
  return node.type === WorkflowNodeType.Branch || node.type === WorkflowNodeType.ExecuteResultBranch;
};

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
    current.next = newNode(WorkflowNodeType.ExecuteResultBranch, {});

    current = current.next!.branches![1];
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
    case WorkflowNodeType.Inspect:
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
        // Object.assign(current, targetNode);
        // TODO: 暂时这么处理，避免 #485 #489，后续再优化
        current.type = targetNode.type;
        current.name = targetNode.name;
        current.config = targetNode.config;
        current.inputs = targetNode.inputs;
        current.outputs = targetNode.outputs;
        current.next = targetNode.next;
        current.branches = targetNode.branches;
        current.validated = targetNode.validated;
        break;
      }

      if (isBranchLike(current)) {
        current.branches ??= [];
        current.branches = current.branches.map((branch) => updateNode(branch, targetNode));
      }

      current = current.next as WorkflowNode;
    }

    return draft;
  });
};

export const addNode = (node: WorkflowNode, previousNodeId: string, targetNode: WorkflowNode) => {
  return produce(node, (draft) => {
    let current = draft;
    while (current) {
      if (current.id === previousNodeId && !isBranchLike(targetNode)) {
        targetNode.next = current.next;
        current.next = targetNode;
        break;
      } else if (current.id === previousNodeId && isBranchLike(targetNode)) {
        targetNode.branches![0].next = current.next;
        current.next = targetNode;
        break;
      }

      if (isBranchLike(current)) {
        current.branches ??= [];
        current.branches = current.branches.map((branch) => addNode(branch, previousNodeId, targetNode));
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

        current.branches ??= [];
        current.branches.push(
          newNode(WorkflowNodeType.Condition, {
            branchIndex: current.branches.length,
          })
        );
        break;
      }

      if (isBranchLike(current)) {
        current.branches ??= [];
        current.branches = current.branches.map((branch) => addBranch(branch, branchNodeId));
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

      if (isBranchLike(current)) {
        current.branches ??= [];
        current.branches = current.branches.map((branch) => removeNode(branch, targetNodeId));
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
        if (!isBranchLike(current)) {
          return draft;
        }

        current.branches ??= [];
        current.branches.splice(branchIndex, 1);

        // 如果仅剩一个分支，删除分支节点，将分支节点的下一个节点挂载到当前节点
        if (current.branches.length === 1) {
          const branch = current.branches[0];
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

      if (isBranchLike(current)) {
        current.branches ??= [];
        current.branches = current.branches.map((branch) => removeBranch(branch, branchNodeId, branchIndex));
      }

      current = current.next as WorkflowNode;
      last = last.next;
    }

    return draft;
  });
};

const typeEqual = (a: WorkflowNodeIO, t: string) => {
  if (t === "all") {
    return true;
  }
  if (a.type === t) {
    return true;
  }
  return false;
};

export const getOutputBeforeNodeId = (root: WorkflowNode, nodeId: string, type: string = "all"): WorkflowNode[] => {
  // 某个分支的节点，不应该能获取到相邻分支上节点的输出
  const outputs: WorkflowNode[] = [];

  const traverse = (current: WorkflowNode, output: WorkflowNode[]) => {
    if (!current) {
      return false;
    }
    if (current.id === nodeId) {
      return true;
    }

    if (current.type !== WorkflowNodeType.Branch && current.outputs && current.outputs.some((io) => typeEqual(io, type))) {
      output.push({
        ...current,
        outputs: current.outputs.filter((io) => typeEqual(io, type)),
      });
    }

    if (isBranchLike(current)) {
      let currentLength = output.length;
      const latestOutput = output.length > 0 ? output[output.length - 1] : null;
      for (const branch of current.branches!) {
        if (branch.type === WorkflowNodeType.ExecuteFailure) {
          output.splice(output.length - 1);
          currentLength -= 1;
        }
        if (traverse(branch, output)) {
          return true;
        }
        // 如果当前分支没有输出，清空之前的输出
        if (output.length > currentLength) {
          output.splice(currentLength);
        }
        if (latestOutput && branch.type === WorkflowNodeType.ExecuteFailure) {
          output.push(latestOutput);
          currentLength += 1;
        }
      }
    }

    return traverse(current.next as WorkflowNode, output);
  };

  traverse(root, outputs);
  return outputs;
};

export const isAllNodesValidated = (node: WorkflowNode): boolean => {
  let current = node as typeof node | undefined;
  while (current) {
    if (isBranchLike(current)) {
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
