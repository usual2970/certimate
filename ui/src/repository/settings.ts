import { type SettingsModel } from "@/domain/settings";
import { getPocketBase } from "./pocketbase";

export const get = async <T>(name: string) => {
  try {
    const resp = await getPocketBase().collection("settings").getFirstListItem<SettingsModel<T>>(`name='${name}'`);
    return resp;
  } catch {
    return {
      name: name,
      content: {} as T,
    } as SettingsModel<T>;
  }
};

export const save = async <T>(record: MaybeModelRecordWithId<SettingsModel<T>>) => {
  if (record.id) {
    return await getPocketBase().collection("settings").update<SettingsModel<T>>(record.id, record);
  }

  return await getPocketBase().collection("settings").create<SettingsModel<T>>(record);
};
