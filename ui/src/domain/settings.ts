export type Setting<T> = {
  id?: string;
  name?: string;
  content?: T;
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
  email?: NotifyChannelEmail;
  webhook?: NotifyChannel;
  dingtalk?: NotifyChannel;
  lark?: NotifyChannel;
  telegram?: NotifyChannel;
  serverchan?: NotifyChannel;
  bark?: NotifyChannelBark;
};

export type NotifyChannel =
  | NotifyChannelEmail
  | NotifyChannelWebhook
  | NotifyChannelDingTalk
  | NotifyChannelLark
  | NotifyChannelTelegram
  | NotifyChannelServerChan
  | NotifyChannelBark;

export type NotifyChannelEmail = {
  smtpHost: string;
  smtpPort: number;
  smtpTLS: boolean;
  username: string;
  password: string;
  senderAddress: string;
  receiverAddress: string;
  enabled: boolean;
};

export type NotifyChannelWebhook = {
  url: string;
  enabled: boolean;
};

export type NotifyChannelDingTalk = {
  accessToken: string;
  secret: string;
  enabled: boolean;
};

export type NotifyChannelLark = {
  webhookUrl: string;
  enabled: boolean;
};

export type NotifyChannelTelegram = {
  apiToken: string;
  chatId: string;
  enabled: boolean;
};

export type NotifyChannelServerChan = {
  url: string;
  enabled: boolean;
};

export type NotifyChannelBark = {
  deviceKey: string;
  serverUrl: string;
  enabled: boolean;
};

export const defaultNotifyTemplate: NotifyTemplate = {
  title: "您有 {COUNT} 张证书即将过期",
  content: "有 {COUNT} 张证书即将过期，域名分别为 {DOMAINS}，请保持关注！",
};

export type SSLProvider = "letsencrypt" | "zerossl" | "gts";

export type SSLProviderSetting = {
  provider: SSLProvider;
  config: {
    [key: string]: {
      [key: string]: string;
    };
  };
};
