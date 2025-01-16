import { ClientResponseError } from "pocketbase";

import { WORKFLOW_TRIGGERS } from "@/domain/workflow";
import { getPocketBase } from "@/repository/pocketbase";

export const run = async (id: string) => {
  const pb = getPocketBase();

  const resp = await pb.send<BaseResponse>("/api/workflow/run", {
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
