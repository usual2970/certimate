import { Workflow, WorkflowNode } from "@/domain/workflow";
import { getPb } from "./api";

export const get = async (id: string) => {
  const response = await getPb().collection("workflow").getOne<Workflow>(id);
  return response;
};

export const save = async (data: Record<string, string | boolean | WorkflowNode>) => {
  if (data.id) {
    return await getPb()
      .collection("workflow")
      .update<Workflow>(data.id as string, data);
  }
  return await getPb().collection("workflow").create<Workflow>(data);
};
