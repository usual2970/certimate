import moment from "moment";

import { type Access } from "@/domain/access";
import { getPocketBase } from "./pocketbase";

export const list = async () => {
  return await getPocketBase().collection("access").getFullList<Access>({
    sort: "-created",
    filter: "deleted = null",
  });
};

export const save = async (record: Access) => {
  if (record.id) {
    return await getPocketBase().collection("access").update(record.id, record);
  }

  return await getPocketBase().collection("access").create(record);
};

export const remove = async (record: Access) => {
  record.deleted = moment.utc().format("YYYY-MM-DD HH:mm:ss");
  return await getPocketBase().collection("access").update(record.id, record);
};
