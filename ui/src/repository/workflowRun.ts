import { type WorkflowRunModel } from "@/domain/workflowRun";

import { getPocketBase } from "./pocketbase";

const COLLECTION_NAME = "workflow_run_log";

export type ListWorkflowRunsRequest = {
  workflowId: string;
  page?: number;
  perPage?: number;
};

export const list = async (request: ListWorkflowRunsRequest) => {
  const page = request.page || 1;
  const perPage = request.perPage || 10;

  return await getPocketBase()
    .collection(COLLECTION_NAME)
    .getList<WorkflowRunModel>(page, perPage, {
      filter: getPocketBase().filter("workflowId={:workflowId}", { workflowId: request.workflowId }),
      sort: "-created",
      requestKey: null,
    });
};
