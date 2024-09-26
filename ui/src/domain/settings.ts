export type Setting = {
  id?: string;
  name?: string;
  content?:
    | EmailsSetting
    | NotifyTemplates
    | NotifyChannels
    | SSLProviderSetting;
};

export type EmailsSetting = {
  emails: string[];
};

export type NotifyTemplates = {
  notifyTemplates: NotifyTemplate[];
};

export type NotifyTemplate = {
  title: string;
  content: string;
};

export type NotifyChannels = {
  dingtalk?: NotifyChannel;
  telegram?: NotifyChannel;
  webhook?: NotifyChannel;
};

export type NotifyChannel =
  | NotifyChannelDingTalk
  | NotifyChannelTelegram
  | NotifyChannelWebhook;

export type NotifyChannelDingTalk = {
  accessToken: string;
  secret: string;
  enabled: boolean;
};

export type NotifyChannelTelegram = {
  apiToken: string;
  chatId: string;
  enabled: boolean;
};

export type NotifyChannelWebhook = {
  url: string;
  enabled: boolean;
};

export const defaultNotifyTemplate: NotifyTemplate = {
  title: "您有{COUNT}张证书即将过期",
  content: "有{COUNT}张证书即将过期,域名分别为{DOMAINS},请保持关注！",
};

export type SSLProvider = "letsencrypt" | "zerossl";

export type SSLProviderSetting = {
  provider: SSLProvider;
  config: {
    [key: string]: {
      [key: string]: string;
    };
  };
};
