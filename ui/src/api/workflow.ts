import { ClientResponseError } from "pocketbase";

import { getPocketBase } from "@/repository/pocketbase";

export const run = async (id: string) => {
  const pb = getPocketBase();

  const resp = await pb.send("/api/workflow/run", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      id,
    },
  });

  if (resp.status != 0 && resp.status !== 200) {
    throw new ClientResponseError({ status: resp.status, response: resp, data: {} });
  }

  return resp;
};
