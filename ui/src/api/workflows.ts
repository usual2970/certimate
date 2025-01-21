import { ClientResponseError } from "pocketbase";

import { WORKFLOW_TRIGGERS } from "@/domain/workflow";
import { getPocketBase } from "@/repository/_pocketbase";

export const startRun = async (workflowId: string) => {
  const pb = getPocketBase();

  const resp = await pb.send<BaseResponse>(`/api/workflows/${encodeURIComponent(workflowId)}/runs`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      trigger: WORKFLOW_TRIGGERS.MANUAL,
    },
  });

  if (resp.code != 0) {
    throw new ClientResponseError({ status: resp.code, response: resp, data: {} });
  }

  return resp;
};

export const cancelRun = async (workflowId: string, runId: string) => {
  const pb = getPocketBase();

  const resp = await pb.send<BaseResponse>(`/api/workflows/${encodeURIComponent(workflowId)}/runs/${encodeURIComponent(runId)}/cancel`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (resp.code != 0) {
    throw new ClientResponseError({ status: resp.code, response: resp, data: {} });
  }

  return resp;
};
