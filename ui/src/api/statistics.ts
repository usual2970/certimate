import { ClientResponseError } from "pocketbase";

import { type Statistics } from "@/domain/statistics";
import { getPocketBase } from "@/repository/pocketbase";

export const get = async () => {
  const pb = getPocketBase();

  const resp = await pb.send("/api/statistics/get", {
    method: "GET",
  });

  if (resp.status != 0 && resp.status !== 200) {
    throw new ClientResponseError({ status: resp.status, response: resp, data: {} });
  }

  return resp.data as Statistics;
};
