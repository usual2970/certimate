import dayjs from "dayjs";

import { type AccessModel } from "@/domain/access";
import { getPocketBase } from "./_pocketbase";

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
  await getPocketBase()
    .collection(COLLECTION_NAME)
    .update<AccessModel>(record.id!, { deleted: dayjs.utc().format("YYYY-MM-DD HH:mm:ss") });
  return true;
};
