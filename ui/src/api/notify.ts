import { ClientResponseError } from "pocketbase";

import { getPocketBase } from "@/repository/pocketbase";

export const notifyTest = async (channel: string) => {
  const pb = getPocketBase();

  const resp = await pb.send("/api/notify/test", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      channel,
    },
  });

  if (resp.status != 0 && resp.status !== 200) {
    throw new ClientResponseError({ status: resp.status, response: resp, data: {} });
  }

  return resp;
};
