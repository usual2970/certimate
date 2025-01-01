import { create } from "zustand";
import { produce } from "immer";

import { SETTINGS_NAMES, type EmailsSettingsContent, type SettingsModel } from "@/domain/settings";
import { get as getSettings, save as saveSettings } from "@/repository/settings";

export interface ContactState {
  emails: string[];
  loading: boolean;
  loadedAtOnce: boolean;

  fetchEmails: () => Promise<void>;
  setEmails: (emails: string[]) => Promise<void>;
  addEmail: (email: string) => Promise<void>;
  removeEmail: (email: string) => Promise<void>;
}

export const useContactStore = create<ContactState>((set, get) => {
  let fetcher: Promise<SettingsModel<EmailsSettingsContent>> | null = null; // 防止多次重复请求
  let settings: SettingsModel<EmailsSettingsContent>; // 记录当前设置的其他字段，保存回数据库时用

  return {
    emails: [],
    loading: false,
    loadedAtOnce: false,

    fetchEmails: async () => {
      fetcher ??= getSettings<EmailsSettingsContent>(SETTINGS_NAMES.EMAILS);

      try {
        set({ loading: true });
        settings = await fetcher;
        set({ emails: settings.content.emails?.sort() ?? [], loadedAtOnce: true });
      } finally {
        fetcher = null;
        set({ loading: false });
      }
    },

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
          state.emails = settings.content.emails?.sort() ?? [];
          state.loadedAtOnce = true;
        })
      );
    },

    addEmail: async (email) => {
      const emails = produce(get().emails, (draft) => {
        if (draft.includes(email)) return;
        draft.push(email);
        draft.sort();
      });
      get().setEmails(emails);
    },

    removeEmail: async (email) => {
      const emails = produce(get().emails, (draft) => {
        draft = draft.filter((e) => e !== email);
        draft.sort();
      });
      get().setEmails(emails);
    },
  };
});
