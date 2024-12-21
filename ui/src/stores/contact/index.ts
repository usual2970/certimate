import { create } from "zustand";
import { produce } from "immer";

import { SETTINGS_NAMES, type EmailsSettingsContent, type SettingsModel } from "@/domain/settings";
import { get as getSettings, save as saveSettings } from "@/repository/settings";

export interface ContactState {
  initialized: boolean;
  emails: string[];
  setEmails: (emails: string[]) => void;
  fetchEmails: () => Promise<void>;
}

export const useContactStore = create<ContactState>((set) => {
  let fetcher: Promise<SettingsModel<EmailsSettingsContent>> | null = null; // 防止多次重复请求
  let settings: SettingsModel<EmailsSettingsContent>; // 记录当前设置的其他字段，保存回数据库时用

  return {
    initialized: false,
    emails: [],

    setEmails: async (emails) => {
      settings ??= await getSettings<EmailsSettingsContent>(SETTINGS_NAMES.EMAILS);
      settings = await saveSettings<EmailsSettingsContent>({
        ...settings,
        content: {
          ...settings.content,
          emails: emails,
        },
      });

      set(
        produce((state: ContactState) => {
          state.emails = settings.content.emails;
          state.initialized = true;
        })
      );
    },

    fetchEmails: async () => {
      fetcher ??= getSettings<EmailsSettingsContent>(SETTINGS_NAMES.EMAILS);

      try {
        settings = await fetcher;
        set({ emails: settings.content.emails?.sort() ?? [], initialized: true });
      } finally {
        fetcher = null;
      }
    },
  };
});
