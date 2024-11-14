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

type WorkflowListReq = {
  page: number;
  perPage?: number;
};
export const list = async (req: WorkflowListReq) => {
  let page = 1;
  if (req.page) {
    page = req.page;
  }
  let perPage = 10;
  if (req.perPage) {
    perPage = req.perPage;
  }

  const response = await getPb().collection("workflow").getList<Workflow>(page, perPage, {
    sort: "-created",
  });

  return response;
};

export const remove = async (id: string) => {
  return await getPb().collection("workflow").delete(id);
};
