import { type RecordListOptions, type RecordSubscription } from "pocketbase";

import { type WorkflowModel } from "@/domain/workflow";
import { getPocketBase } from "./_pocketbase";

const COLLECTION_NAME = "workflow";

export type ListWorkflowRequest = {
  page?: number;
  perPage?: number;
  enabled?: boolean;
};

export const list = async (request: ListWorkflowRequest) => {
  const pb = getPocketBase();

  const page = request.page || 1;
  const perPage = request.perPage || 10;

  const options: RecordListOptions = {
    sort: "-created",
    requestKey: null,
  };

  if (request.enabled != null) {
    options.filter = pb.filter("enabled={:enabled}", { enabled: request.enabled });
  }

  return await pb.collection(COLLECTION_NAME).getList<WorkflowModel>(page, perPage, options);
};

export const get = async (id: string) => {
  return await getPocketBase().collection(COLLECTION_NAME).getOne<WorkflowModel>(id, {
    requestKey: null,
  });
};

export const save = async (record: MaybeModelRecord<WorkflowModel>) => {
  if (record.id) {
    return await getPocketBase()
      .collection(COLLECTION_NAME)
      .update<WorkflowModel>(record.id as string, record);
  }

  return await getPocketBase().collection(COLLECTION_NAME).create<WorkflowModel>(record);
};

export const remove = async (record: MaybeModelRecordWithId<WorkflowModel>) => {
  return await getPocketBase().collection(COLLECTION_NAME).delete(record.id);
};

export const subscribe = async (id: string, cb: (e: RecordSubscription<WorkflowModel>) => void) => {
  return getPocketBase().collection(COLLECTION_NAME).subscribe(id, cb);
};

export const unsubscribe = async (id: string) => {
  return getPocketBase().collection(COLLECTION_NAME).unsubscribe(id);
};
