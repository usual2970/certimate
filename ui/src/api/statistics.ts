import { getPb } from "@/repository/api";

export const get = async () => {
  const pb = getPb();

  const resp = await pb.send("/api/statistics/get", {
    method: "GET",
  });

  if (resp.code != 0) {
    throw new Error(resp.msg);
  }

  return resp.data;
};
