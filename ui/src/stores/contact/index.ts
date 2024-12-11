import { create } from "zustand";
import { produce } from "immer";

import { type EmailsSettingsContent, type SettingsModel } from "@/domain/settings";
import { get as getSettings, save as saveSettings } from "@/repository/settings";

export interface ContactState {
  emails: string[];
  setEmails: (emails: string[]) => void;
  fetchEmails: () => Promise<void>;
}

export const useContactStore = create<ContactState>((set) => {
  let settings: SettingsModel<EmailsSettingsContent>;

  return {
    emails: [],

    setEmails: async (emails: string[]) => {
      settings ??= await getSettings<EmailsSettingsContent>("emails");
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
        })
      );
    },

    fetchEmails: async () => {
      settings = await getSettings<EmailsSettingsContent>("emails");

      set({
        emails: settings.content.emails?.sort() ?? [],
      });
    },
  };
});
