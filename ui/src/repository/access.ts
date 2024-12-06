import moment from "moment";

import { type Access } from "@/domain/access";
import { getPocketBase } from "./pocketbase";

export const list = async () => {
  return await getPocketBase().collection("access").getFullList<Access>({
    sort: "-created",
    filter: "deleted = null",
  });
};

export const save = async (data: Access) => {
  if (data.id) {
    return await getPocketBase().collection("access").update(data.id, data);
  }
  return await getPocketBase().collection("access").create(data);
};

export const remove = async (data: Access) => {
  data.deleted = moment.utc().format("YYYY-MM-DD HH:mm:ss");
  return await getPocketBase().collection("access").update(data.id, data);
};
