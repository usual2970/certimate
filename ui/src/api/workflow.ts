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

  if (resp.code != 0) {
    throw new Error(resp.msg);
  }

  return resp;
};
