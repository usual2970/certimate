import { type WorkflowLogModel } from "@/domain/workflowLog";

import { COLLECTION_NAME_WORKFLOW_LOG, getPocketBase } from "./_pocketbase";

export const listByWorkflowRunId = async (workflowRunId: string) => {
  const pb = getPocketBase();

  const list = await pb.collection(COLLECTION_NAME_WORKFLOW_LOG).getFullList<WorkflowLogModel>({
    batch: 65535,
    filter: pb.filter("runId={:runId}", { runId: workflowRunId }),
    sort: "timestamp",
    requestKey: null,
  });

  return {
    totalItems: list.length,
    items: list,
  };
};
