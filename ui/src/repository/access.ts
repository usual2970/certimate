import dayjs from "dayjs";

import { type AccessModel } from "@/domain/access";
import { COLLECTION_NAME_ACCESS, getPocketBase } from "./_pocketbase";

export const list = async () => {
  const list = await getPocketBase().collection(COLLECTION_NAME_ACCESS).getFullList<AccessModel>({
    batch: 65535,
    filter: "deleted=null",
    sort: "-created",
    requestKey: null,
  });
  return {
    totalItems: list.length,
    items: list,
  };
};

export const save = async (record: MaybeModelRecord<AccessModel>) => {
  if (record.id) {
    return await getPocketBase().collection(COLLECTION_NAME_ACCESS).update<AccessModel>(record.id, record);
  }

  return await getPocketBase().collection(COLLECTION_NAME_ACCESS).create<AccessModel>(record);
};

export const remove = async (record: MaybeModelRecordWithId<AccessModel>) => {
  await getPocketBase()
    .collection(COLLECTION_NAME_ACCESS)
    .update<AccessModel>(record.id!, { deleted: dayjs.utc().format("YYYY-MM-DD HH:mm:ss") });
  return true;
};
