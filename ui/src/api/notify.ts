import { getPb } from "@/repository/api";

export const notifyTest = async (channel: string) => {
  const pb = getPb();

  const resp = await pb.send("/api/notify/test", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      channel,
    },
  });

  if (resp.code != 0) {
    throw new Error(resp.msg);
  }

  return resp;
};
