import { getPocketBase } from "@/repository/pocketbase";

export const get = async () => {
  const pb = getPocketBase();

  const resp = await pb.send("/api/statistics/get", {
    method: "GET",
  });

  if (resp.code != 0) {
    throw new Error(resp.msg);
  }

  return resp.data;
};
