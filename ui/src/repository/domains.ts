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
    return await getPb().collection("domains").update<Domain>(data.id, data);
  }
  return await getPb().collection("domains").create<Domain>(data);
};

export const remove = async (id: string) => {
  return await getPb().collection("domains").delete(id);
};

type Callback = (data: Domain) => void;
export const subscribeId = (id: string, callback: Callback) => {
  return getPb()
    .collection("domains")
    .subscribe<Domain>(
      id,
      (e) => {
        if (e.action === "update") {
          callback(e.record);
        }
      },
      {
        expand: "lastDeployment",
      }
    );
};

export const unsubscribeId = (id: string) => {
  getPb().collection("domains").unsubscribe(id);
};
