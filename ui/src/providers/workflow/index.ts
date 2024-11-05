import { addBranch, addNode, removeBranch, removeNode, updateNode, WorkflowBranchNode, WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import { create } from "zustand";

export type WorkflowState = {
  root: WorkflowNode;
  updateNode: (node: WorkflowNode) => void;
  addNode: (node: WorkflowNode, preId: string) => void;
  addBranch: (branchId: string) => void;
  removeNode: (nodeId: string) => void;
  removeBranch: (branchId: string, index: number) => void;
};

export const useWorkflowStore = create<WorkflowState>((set) => ({
  root: {
    id: "1",
    name: "开始",
    type: WorkflowNodeType.Start,
    next: {
      id: "2",
      name: "结束",
      type: WorkflowNodeType.Branch,
      branches: [
        {
          id: "3",
          name: "条件1",
          type: WorkflowNodeType.Condition,
          next: {
            id: "4",
            name: "条件2",
            type: WorkflowNodeType.Apply,
          },
        },
        {
          id: "5",
          name: "条件2",
          type: WorkflowNodeType.Condition,
        },
      ],
    },
  },
  updateNode: (node: WorkflowNode | WorkflowBranchNode) => {
    set((state: WorkflowState) => {
      const newRoot = updateNode(state.root, node);
      console.log(newRoot);
      return {
        root: newRoot,
      };
    });
  },
  addNode: (node: WorkflowNode | WorkflowBranchNode, preId: string) =>
    set((state: WorkflowState) => {
      const newRoot = addNode(state.root, preId, node);

      return {
        root: newRoot,
      };
    }),
  addBranch: (branchId: string) =>
    set((state: WorkflowState) => {
      const newRoot = addBranch(state.root, branchId);

      return {
        root: newRoot,
      };
    }),

  removeBranch: (branchId: string, index: number) =>
    set((state: WorkflowState) => {
      const newRoot = removeBranch(state.root, branchId, index);

      return {
        root: newRoot,
      };
    }),

  removeNode: (nodeId: string) =>
    set((state: WorkflowState) => {
      const newRoot = removeNode(state.root, nodeId);

      return {
        root: newRoot,
      };
    }),
}));
