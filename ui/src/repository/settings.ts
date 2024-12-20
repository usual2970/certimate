import { ClientResponseError } from "pocketbase";

import { SETTINGS_NAMES, type SettingsModel } from "@/domain/settings";
import { getPocketBase } from "./pocketbase";

export const get = async <T>(name: (typeof SETTINGS_NAMES)[keyof typeof SETTINGS_NAMES]) => {
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

export const save = async <T>(record: MaybeModelRecordWithId<SettingsModel<T>>) => {
  if (record.id) {
    return await getPocketBase().collection("settings").update<SettingsModel<T>>(record.id, record);
  }

  return await getPocketBase().collection("settings").create<SettingsModel<T>>(record);
};
