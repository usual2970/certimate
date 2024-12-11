import { type RecordListOptions } from "pocketbase";

import { type WorkflowModel, type WorkflowNode, type WorkflowRunLog } from "@/domain/workflow";
import { getPocketBase } from "./pocketbase";

const COLLECTION_NAME = "workflow";

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

  return await pb.collection(COLLECTION_NAME).getList<WorkflowModel>(page, perPage, options);
};

export const get = async (id: string) => {
  return await getPocketBase().collection(COLLECTION_NAME).getOne<WorkflowModel>(id);
};

export const save = async (record: Record<string, string | boolean | WorkflowNode>) => {
  if (record.id) {
    return await getPocketBase()
      .collection(COLLECTION_NAME)
      .update<WorkflowModel>(record.id as string, record);
  }

  return await getPocketBase().collection(COLLECTION_NAME).create<WorkflowModel>(record);
};

export const remove = async (record: WorkflowModel) => {
  return await getPocketBase().collection(COLLECTION_NAME).delete(record.id);
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
