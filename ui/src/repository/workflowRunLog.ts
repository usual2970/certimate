import { type WorkflowRunLog } from "@/domain/workflow";
import { getPocketBase } from "./pocketbase";

const COLLECTION_NAME = "workflow_run_log";

export type ListWorkflowLogsRequest = {
  id: string;
  page?: number;
  perPage?: number;
};

export const list = async (request: ListWorkflowLogsRequest) => {
  const page = request.page || 1;
  const perPage = request.perPage || 10;

  return await getPocketBase()
    .collection(COLLECTION_NAME)
    .getList<WorkflowRunLog>(page, perPage, {
      filter: getPocketBase().filter("workflow={:workflowId}", { workflowId: request.id }),
      sort: "-created",
      requestKey: null,
    });
};
