import { type RecordListOptions } from "pocketbase";

import { type Workflow, type WorkflowNode, type WorkflowRunLog } from "@/domain/workflow";
import { getPocketBase } from "./pocketbase";

export type WorkflowListReq = {
  page?: number;
  perPage?: number;
  enabled?: boolean;
};

export const list = async (req: WorkflowListReq) => {
  const pb = getPocketBase();

  const page = req.page || 1;
  const perPage = req.perPage || 10;

  const options: RecordListOptions = { sort: "-created" };
  if (req.enabled != null) {
    options.filter = pb.filter("enabled={:enabled}", { enabled: req.enabled });
  }

  return await pb.collection("workflow").getList<Workflow>(page, perPage, options);
};

export const get = async (id: string) => {
  return await getPocketBase().collection("workflow").getOne<Workflow>(id);
};

export const save = async (data: Record<string, string | boolean | WorkflowNode>) => {
  if (data.id) {
    return await getPocketBase()
      .collection("workflow")
      .update<Workflow>(data.id as string, data);
  }

  return await getPocketBase().collection("workflow").create<Workflow>(data);
};

export const remove = async (id: string) => {
  return await getPocketBase().collection("workflow").delete(id);
};

type WorkflowLogsReq = {
  id: string;
  page?: number;
  perPage?: number;
};

export const logs = async (req: WorkflowLogsReq) => {
  const page = req.page || 1;
  const perPage = req.perPage || 10;

  return await getPocketBase()
    .collection("workflow_run_log")
    .getList<WorkflowRunLog>(page, perPage, {
      sort: "-created",
      filter: getPocketBase().filter("workflow={:workflowId}", { workflowId: req.id }),
    });
};
