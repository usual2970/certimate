import { Settings } from "@/domain/settings";
import { getPocketBase } from "./pocketbase";

export const get = async <T>(name: string) => {
  try {
    const resp = await getPocketBase().collection("settings").getFirstListItem<Settings<T>>(`name='${name}'`);
    return resp;
  } catch {
    return {
      name: name,
      content: {} as T,
    } as Settings<T>;
  }
};

export const save = async <T>(record: Settings<T>) => {
  if (record.id) {
    return await getPocketBase().collection("settings").update<Settings<T>>(record.id, record);
  }

  return await getPocketBase().collection("settings").create<Settings<T>>(record);
};
