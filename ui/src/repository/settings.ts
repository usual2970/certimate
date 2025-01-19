import { ClientResponseError } from "pocketbase";

import { type SettingsModel, type SettingsNames } from "@/domain/settings";
import { getPocketBase } from "./_pocketbase";

const COLLECTION_NAME = "settings";

export const get = async <T extends NonNullable<unknown>>(name: SettingsNames) => {
  try {
    const resp = await getPocketBase().collection(COLLECTION_NAME).getFirstListItem<SettingsModel<T>>(`name='${name}'`, {
      requestKey: null,
    });
    return resp;
  } catch (err) {
    if (err instanceof ClientResponseError && err.status === 404) {
      return {
        name: name,
        content: {} as T,
      } as SettingsModel<T>;
    }

    throw err;
  }
};

export const save = async <T extends NonNullable<unknown>>(record: MaybeModelRecordWithId<SettingsModel<T>>) => {
  if (record.id) {
    return await getPocketBase().collection(COLLECTION_NAME).update<SettingsModel<T>>(record.id, record);
  }

  return await getPocketBase().collection(COLLECTION_NAME).create<SettingsModel<T>>(record);
};
