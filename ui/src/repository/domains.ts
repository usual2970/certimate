import { Domain } from "@/domain/domain";
import { getPb } from "./api";

export const list = async () => {
  const response = getPb().collection("domains").getFullList<Domain>({
    sort: "-created",
    expand: "lastDeployment",
  });

  return response;
};

export const get = async (id: string) => {
  const response = await getPb().collection("domains").getOne<Domain>(id);
  return response;
};

export const save = async (data: Domain) => {
  if (data.id) {
    return await getPb().collection("domains").update(data.id, data);
  }
  return await getPb().collection("domains").create(data);
};

export const remove = async (id: string) => {
  return await getPb().collection("domains").delete(id);
};
