import { type WorkflowRunModel } from "@/domain/workflowRun";

import { getPocketBase } from "./_pocketbase";

const COLLECTION_NAME = "workflow_run";

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
    .collection(COLLECTION_NAME)
    .getList<WorkflowRunModel>(page, perPage, {
      filter: getPocketBase().filter(filter, params),
      sort: "-created",
      requestKey: null,
      expand: request.expand ? "workflowId" : undefined,
    });
};

export const remove = async (record: MaybeModelRecordWithId<WorkflowRunModel>) => {
  return await getPocketBase().collection(COLLECTION_NAME).delete(record.id);
};
