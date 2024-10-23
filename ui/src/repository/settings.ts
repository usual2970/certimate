import { EmailsSetting, Setting } from "@/domain/settings";
import { getPb } from "./api";

export const getEmails = async () => {
  try {
    const resp = await getPb().collection("settings").getFirstListItem<Setting<EmailsSetting>>("name='emails'");
    return resp;
  } catch (e) {
    return {
      content: { emails: [] },
    };
  }
};

export const getSetting = async <T>(name: string) => {
  try {
    const resp = await getPb().collection("settings").getFirstListItem<Setting<T>>(`name='${name}'`);
    return resp;
  } catch (e) {
    const rs: Setting<T> = {
      name: name,
    };
    return rs;
  }
};

export const update = async <T>(setting: Setting<T>) => {
  const pb = getPb();
  let resp: Setting<T>;
  if (setting.id) {
    resp = await pb.collection("settings").update(setting.id, setting);
  } else {
    resp = await pb.collection("settings").create(setting);
  }
  return resp;
};
