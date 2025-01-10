import { ClientResponseError, type RecordSubscription } from "pocketbase";

import { WORKFLOW_TRIGGERS, type WorkflowModel } from "@/domain/workflow";
import { getPocketBase } from "@/repository/pocketbase";

export const run = async (id: string) => {
  const pb = getPocketBase();

  const resp = await pb.send("/api/workflow/run", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      workflowId: id,
      trigger: WORKFLOW_TRIGGERS.MANUAL,
    },
  });

  if (resp.code != 0) {
    throw new ClientResponseError({ status: resp.code, response: resp, data: {} });
  }

  return resp;
};

export const subscribe = async (id: string, cb: (e: RecordSubscription<WorkflowModel>) => void) => {
  const pb = getPocketBase();

  return pb.collection("workflow").subscribe(id, cb);
};

export const unsubscribe = async (id: string) => {
  const pb = getPocketBase();

  return pb.collection("workflow").unsubscribe(id);
};
