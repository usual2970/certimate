import { Access } from "@/domain/access";
import { getPb } from "./api";
import moment from "moment";

export const list = async () => {
  return await getPb().collection("access").getFullList<Access>({
    sort: "-created",
    filter: "deleted = null",
  });
};

export const save = async (data: Access) => {
  if (data.id) {
    return await getPb().collection("access").update(data.id, data);
  }
  return await getPb().collection("access").create(data);
};

export const remove = async (data: Access) => {
  data.deleted = moment.utc().format("YYYY-MM-DD HH:mm:ss");
  return await getPb().collection("access").update(data.id, data);
};
