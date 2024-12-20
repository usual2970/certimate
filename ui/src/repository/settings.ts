import { ClientResponseError } from "pocketbase";

import { type SettingsModel, type SettingsNames } from "@/domain/settings";
import { getPocketBase } from "./pocketbase";

export const get = async <T extends NonNullable<unknown>>(name: SettingsNames) => {
  try {
    const resp = await getPocketBase().collection("settings").getFirstListItem<SettingsModel<T>>(`name='${name}'`, {
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
    return await getPocketBase().collection("settings").update<SettingsModel<T>>(record.id, record);
  }

  return await getPocketBase().collection("settings").create<SettingsModel<T>>(record);
};
