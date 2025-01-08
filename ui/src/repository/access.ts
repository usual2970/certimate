import dayjs from "dayjs";

import { type AccessModel } from "@/domain/access";
import { getPocketBase } from "./pocketbase";

const COLLECTION_NAME = "access";

export const list = async () => {
  return await getPocketBase().collection(COLLECTION_NAME).getFullList<AccessModel>({
    filter: "deleted=null",
    sort: "-created",
    requestKey: null,
  });
};

export const save = async (record: MaybeModelRecord<AccessModel>) => {
  if (record.id) {
    return await getPocketBase().collection(COLLECTION_NAME).update<AccessModel>(record.id, record);
  }

  return await getPocketBase().collection(COLLECTION_NAME).create<AccessModel>(record);
};

export const remove = async (record: MaybeModelRecordWithId<AccessModel>) => {
  record = { ...record, deleted: dayjs.utc().format("YYYY-MM-DD HH:mm:ss") };

  // TODO: 仅为兼容旧版本，后续迭代时删除
  if ("provider" in record && record.provider === "httpreq") record.provider = "acmehttpreq";
  if ("provider" in record && record.provider === "tencent") record.provider = "tencentcloud";
  if ("provider" in record && record.provider === "pdns") record.provider = "powerdns";

  await getPocketBase().collection(COLLECTION_NAME).update<AccessModel>(record.id!, record);
};
