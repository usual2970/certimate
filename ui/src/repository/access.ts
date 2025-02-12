import dayjs from "dayjs";

import { type AccessModel } from "@/domain/access";
import { COLLECTION_NAME_ACCESS, getPocketBase } from "./_pocketbase";

export const list = async () => {
  return await getPocketBase().collection(COLLECTION_NAME_ACCESS).getFullList<AccessModel>({
    filter: "deleted=null",
    sort: "-created",
    requestKey: null,
  });
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
