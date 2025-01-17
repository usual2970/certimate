import { ClientResponseError } from "pocketbase";

import { WORKFLOW_TRIGGERS } from "@/domain/workflow";
import { getPocketBase } from "@/repository/_pocketbase";

export const run = async (id: string) => {
  const pb = getPocketBase();

  const resp = await pb.send<BaseResponse>(`/api/workflows/${encodeURIComponent(id)}/run`, {
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
