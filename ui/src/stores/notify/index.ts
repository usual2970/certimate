import { create } from "zustand";
import { produce } from "immer";

import { SETTINGS_NAMES, type NotifyChannelsSettingsContent, type SettingsModel } from "@/domain/settings";
import { get as getSettings, save as saveSettings } from "@/repository/settings";

export interface NotifyChannelState {
  initialized: boolean;
  channels: NotifyChannelsSettingsContent;
  setChannel: (channel: keyof NotifyChannelsSettingsContent, config: NotifyChannelsSettingsContent[keyof NotifyChannelsSettingsContent]) => void;
  setChannels: (channels: NotifyChannelsSettingsContent) => void;
  fetchChannels: () => Promise<void>;
}

export const useNotifyChannelStore = create<NotifyChannelState>((set, get) => {
  let fetcher: Promise<SettingsModel<NotifyChannelsSettingsContent>> | null = null; // 防止多次重复请求
  let settings: SettingsModel<NotifyChannelsSettingsContent>; // 记录当前设置的其他字段，保存回数据库时用

  return {
    initialized: false,
    channels: {},

    setChannel: async (channel, config) => {
      settings ??= await getSettings<NotifyChannelsSettingsContent>(SETTINGS_NAMES.NOTIFY_CHANNELS);
      return get().setChannels({
        ...settings.content,
        [channel]: { ...settings.content[channel], ...config },
      });
    },

    setChannels: async (channels) => {
      settings ??= await getSettings<NotifyChannelsSettingsContent>(SETTINGS_NAMES.NOTIFY_CHANNELS);
      settings = await saveSettings<NotifyChannelsSettingsContent>({
        ...settings,
        content: channels,
      });

      set(
        produce((state: NotifyChannelState) => {
          state.channels = settings.content;
          state.initialized = true;
        })
      );
    },

    fetchChannels: async () => {
      fetcher ??= getSettings<NotifyChannelsSettingsContent>(SETTINGS_NAMES.NOTIFY_CHANNELS);

      try {
        settings = await fetcher;
        set({ channels: settings.content ?? {}, initialized: true });
      } finally {
        fetcher = null;
      }
    },
  };
});
