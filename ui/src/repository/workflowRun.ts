import { type RecordSubscription } from "pocketbase";

import { type WorkflowRunModel } from "@/domain/workflowRun";

import { COLLECTION_NAME_WORKFLOW_RUN, getPocketBase } from "./_pocketbase";

export type ListWorkflowRunsRequest = {
  workflowId?: string;
  page?: number;
  perPage?: number;
  expand?: boolean;
};

export const list = async (request: ListWorkflowRunsRequest) => {
  const page = request.page || 1;
  const perPage = request.perPage || 10;

  let filter = "";
  const params: Record<string, string> = {};
  if (request.workflowId) {
    filter = `workflowId={:workflowId}`;
    params.workflowId = request.workflowId;
  }

  return await getPocketBase()
    .collection(COLLECTION_NAME_WORKFLOW_RUN)
    .getList<WorkflowRunModel>(page, perPage, {
      filter: getPocketBase().filter(filter, params),
      sort: "-created",
      requestKey: null,
      expand: request.expand ? "workflowId" : undefined,
    });
};

export const remove = async (record: MaybeModelRecordWithId<WorkflowRunModel>) => {
  return await getPocketBase().collection(COLLECTION_NAME_WORKFLOW_RUN).delete(record.id);
};

export const subscribe = async (id: string, cb: (e: RecordSubscription<WorkflowRunModel>) => void) => {
  return getPocketBase().collection(COLLECTION_NAME_WORKFLOW_RUN).subscribe(id, cb);
};

export const unsubscribe = async (id: string) => {
  return getPocketBase().collection(COLLECTION_NAME_WORKFLOW_RUN).unsubscribe(id);
};
