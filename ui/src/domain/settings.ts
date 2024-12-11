import { type BaseModel } from "pocketbase";

export interface SettingsModel<T> extends BaseModel {
  name: string;
  content: T;
}

export type EmailsSettingsContent = {
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
  [key: string]: NotifyChannel;
};

export type NotifyChannel =
  | NotifyChannelEmail
  | NotifyChannelWebhook
  | NotifyChannelDingTalk
  | NotifyChannelLark
  | NotifyChannelTelegram
  | NotifyChannelServerChan
  | NotifyChannelBark;

type ChannelLabel = {
  name: string;
  label: string;
};
export const channels: ChannelLabel[] = [
  {
    name: "dingtalk",
    label: "common.provider.dingtalk",
  },
  {
    name: "lark",
    label: "common.provider.lark",
  },
  {
    name: "telegram",
    label: "common.provider.telegram",
  },
  {
    name: "webhook",
    label: "common.provider.webhook",
  },
  {
    name: "serverchan",
    label: "common.provider.serverchan",
  },
  {
    name: "email",
    label: "common.provider.email",
  },
  {
    name: "bark",
    label: "common.provider.bark",
  },
];

export const channelLabelMap: Map<string, ChannelLabel> = new Map(channels.map((item) => [item.name, item]));
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
