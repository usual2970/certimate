import { Setting } from "@/domain/settings";
import { getPb } from "./api";

export const getEmails = async () => {
  try {
    const resp = await getPb()
      .collection("settings")
      .getFirstListItem<Setting>("name='emails'");
    return resp;
  } catch (e) {
    return {
      content: { emails: [] },
    };
  }
};

export const update = async (setting: Setting) => {
  const pb = getPb();
  let resp: Setting;
  if (setting.id) {
    resp = await pb.collection("settings").update(setting.id, setting);
  } else {
    resp = await pb.collection("settings").create(setting);
  }
  return resp;
};
