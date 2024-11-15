import {
  addBranch,
  addNode,
  getWorkflowOutputBeforeId,
  initWorkflow,
  removeBranch,
  removeNode,
  updateNode,
  Workflow,
  WorkflowBranchNode,
  WorkflowNode,
  WorkflowNodeType,
} from "@/domain/workflow";
import { save, get as getWrokflow } from "@/repository/workflow";
import { create } from "zustand";

export type WorkflowState = {
  workflow: Workflow;
  initialized: boolean;
  updateNode: (node: WorkflowNode) => void;
  addNode: (node: WorkflowNode, preId: string) => void;
  addBranch: (branchId: string) => void;
  removeNode: (nodeId: string) => void;
  removeBranch: (branchId: string, index: number) => void;
  getWorkflowOuptutBeforeId: (id: string, type: string) => WorkflowNode[];
  switchEnable(): void;
  save(): void;
  init(id?: string): void;
  setBaseInfo: (name: string, description: string) => void;
};

export const useWorkflowStore = create<WorkflowState>((set, get) => ({
  workflow: {
    id: "",
    name: "placeholder",
    type: WorkflowNodeType.Start,
  },
  initialized: false,
  init: async (id?: string) => {
    let data = {
      id: "",
      name: "placeholder",
      type: "auto",
    };

    if (!id) {
      data = initWorkflow();
    } else {
      data = await getWrokflow(id);
    }

    set({
      workflow: data,
      initialized: true,
    });
  },
  setBaseInfo: async (name: string, description: string) => {
    const data: Record<string, string | boolean | WorkflowNode> = {
      id: (get().workflow.id as string) ?? "",
      name: name,
      description: description,
    };
    if (!data.id) {
      data.draft = get().workflow.draft as WorkflowNode;
    }
    const resp = await save(data);
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          name,
          description,
          id: resp.id,
        },
      };
    });
  },
  switchEnable: async () => {
    const resp = await save({
      id: (get().workflow.id as string) ?? "",
      content: get().workflow.draft as WorkflowNode,
      enabled: !get().workflow.enabled,
      hasDraft: false,
    });
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          id: resp.id,
          content: resp.content,
          enabled: resp.enabled,
          hasDraft: false,
        },
      };
    });
  },
  save: async () => {
    const resp = await save({
      id: (get().workflow.id as string) ?? "",
      content: get().workflow.draft as WorkflowNode,
      hasDraft: false,
    });
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          id: resp.id,
          content: resp.content,
          hasDraft: false,
        },
      };
    });
  },
  updateNode: async (node: WorkflowNode | WorkflowBranchNode) => {
    const newRoot = updateNode(get().workflow.draft as WorkflowNode, node);
    const resp = await save({
      id: (get().workflow.id as string) ?? "",
      draft: newRoot,
      hasDraft: true,
    });
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          draft: newRoot,
          id: resp.id,
          hasDraft: true,
        },
      };
    });
  },
  addNode: async (node: WorkflowNode | WorkflowBranchNode, preId: string) => {
    const newRoot = addNode(get().workflow.draft as WorkflowNode, preId, node);
    const resp = await save({
      id: (get().workflow.id as string) ?? "",
      draft: newRoot,
      hasDraft: true,
    });
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          draft: newRoot,
          id: resp.id,
          hasDraft: true,
        },
      };
    });
  },
  addBranch: async (branchId: string) => {
    const newRoot = addBranch(get().workflow.draft as WorkflowNode, branchId);
    const resp = await save({
      id: (get().workflow.id as string) ?? "",
      draft: newRoot,
      hasDraft: true,
    });
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          draft: newRoot,
          id: resp.id,
          hasDraft: true,
        },
      };
    });
  },
  removeBranch: async (branchId: string, index: number) => {
    const newRoot = removeBranch(get().workflow.draft as WorkflowNode, branchId, index);
    const resp = await save({
      id: (get().workflow.id as string) ?? "",
      draft: newRoot,
      hasDraft: true,
    });
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          draft: newRoot,
          id: resp.id,
          hasDraft: true,
        },
      };
    });
  },
  removeNode: async (nodeId: string) => {
    const newRoot = removeNode(get().workflow.draft as WorkflowNode, nodeId);
    const resp = await save({
      id: (get().workflow.id as string) ?? "",
      draft: newRoot,
      hasDraft: true,
    });
    set((state: WorkflowState) => {
      return {
        workflow: {
          ...state.workflow,
          draft: newRoot,
          id: resp.id,
          hasDraft: true,
        },
      };
    });
  },
  getWorkflowOuptutBeforeId: (id: string, type: string) => {
    return getWorkflowOutputBeforeId(get().workflow.draft as WorkflowNode, id, type);
  },
}));

