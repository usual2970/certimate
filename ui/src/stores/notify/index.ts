import { produce } from "immer";
import { create } from "zustand";

import { type NotifyChannelsSettingsContent, SETTINGS_NAMES, type SettingsModel } from "@/domain/settings";
import { get as getSettings, save as saveSettings } from "@/repository/settings";

/**
 * @deprecated
 */
export interface NotifyChannelsState {
  channels: NotifyChannelsSettingsContent;
  loading: boolean;
  loadedAtOnce: boolean;

  fetchChannels: () => Promise<void>;
  setChannel: (channel: keyof NotifyChannelsSettingsContent, config: NotifyChannelsSettingsContent[keyof NotifyChannelsSettingsContent]) => Promise<void>;
  setChannels: (channels: NotifyChannelsSettingsContent) => Promise<void>;
}

/**
 * @deprecated
 */
export const useNotifyChannelsStore = create<NotifyChannelsState>((set, get) => {
  let fetcher: Promise<SettingsModel<NotifyChannelsSettingsContent>> | null = null; // 防止多次重复请求
  let settings: SettingsModel<NotifyChannelsSettingsContent>; // 记录当前设置的其他字段，保存回数据库时用

  return {
    channels: {},
    loading: false,
    loadedAtOnce: false,

    fetchChannels: async () => {
      fetcher ??= getSettings<NotifyChannelsSettingsContent>(SETTINGS_NAMES.NOTIFY_CHANNELS);

      try {
        set({ loading: true });
        settings = await fetcher;
        set({ channels: settings.content ?? {}, loadedAtOnce: true });
      } finally {
        fetcher = null;
        set({ loading: false });
      }
    },

    setChannel: async (channel, config) => {
      settings ??= await getSettings<NotifyChannelsSettingsContent>(SETTINGS_NAMES.NOTIFY_CHANNELS);
      return get().setChannels(
        produce(settings, (draft) => {
          draft.content ??= {};
          draft.content[channel] = { ...draft.content[channel], ...config };
        }).content
      );
    },

    setChannels: async (channels) => {
      settings ??= await getSettings<NotifyChannelsSettingsContent>(SETTINGS_NAMES.NOTIFY_CHANNELS);
      settings = await saveSettings<NotifyChannelsSettingsContent>({
        ...settings,
        content: channels,
      });

      set(
        produce((state: NotifyChannelsState) => {
          state.channels = settings.content;
          state.loadedAtOnce = true;
        })
      );
    },
  };
});
