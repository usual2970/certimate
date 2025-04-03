import { type ApplyCAProviderType } from "./provider";

export const SETTINGS_NAMES = Object.freeze({
  EMAILS: "emails",
  NOTIFY_TEMPLATES: "notifyTemplates",
  NOTIFY_CHANNELS: "notifyChannels",
  SSL_PROVIDER: "sslProvider",
  PERSISTENCE: "persistence",
} as const);

export type SettingsNames = (typeof SETTINGS_NAMES)[keyof typeof SETTINGS_NAMES];

export interface SettingsModel<T extends NonNullable<unknown> = any> extends BaseModel {
  name: string;
  content: T;
}

// #region Settings: Emails
export type EmailsSettingsContent = {
  emails: string[];
};
// #endregion

// #region Settings: NotifyTemplates
export type NotifyTemplatesSettingsContent = {
  notifyTemplates: NotifyTemplate[];
};

export type NotifyTemplate = {
  subject: string;
  message: string;
};

export const defaultNotifyTemplate: NotifyTemplate = {
  subject: "有 ${COUNT} 张证书即将过期",
  message: "有 ${COUNT} 张证书即将过期，域名分别为 ${DOMAINS}，请保持关注！",
};
// #endregion

// #region Settings: NotifyChannels
export const NOTIFY_CHANNELS = Object.freeze({
  BARK: "bark",
  DINGTALK: "dingtalk",
  EMAIL: "email",
  GOTIFY: "gotify",
  LARK: "lark",
  PUSHPLUS: "pushplus",
  SERVERCHAN: "serverchan",
  TELEGRAM: "telegram",
  WEBHOOK: "webhook",
  WECOM: "wecom",
} as const);

export type NotifyChannels = (typeof NOTIFY_CHANNELS)[keyof typeof NOTIFY_CHANNELS];

export type NotifyChannelsSettingsContent = {
  /*
    注意：如果追加新的类型，请保持以 ASCII 排序。
    NOTICE: If you add new type, please keep ASCII order.
  */
  [key: string]: ({ enabled?: boolean } & Record<string, unknown>) | undefined;
  [NOTIFY_CHANNELS.BARK]?: BarkNotifyChannelConfig;
  [NOTIFY_CHANNELS.DINGTALK]?: DingTalkNotifyChannelConfig;
  [NOTIFY_CHANNELS.EMAIL]?: EmailNotifyChannelConfig;
  [NOTIFY_CHANNELS.GOTIFY]?: GotifyNotifyChannelConfig;
  [NOTIFY_CHANNELS.LARK]?: LarkNotifyChannelConfig;
  [NOTIFY_CHANNELS.PUSHPLUS]?: PushPlusNotifyChannelConfig;
  [NOTIFY_CHANNELS.SERVERCHAN]?: ServerChanNotifyChannelConfig;
  [NOTIFY_CHANNELS.TELEGRAM]?: TelegramNotifyChannelConfig;
  [NOTIFY_CHANNELS.WEBHOOK]?: WebhookNotifyChannelConfig;
  [NOTIFY_CHANNELS.WECOM]?: WeComNotifyChannelConfig;
};

export type BarkNotifyChannelConfig = {
  deviceKey: string;
  serverUrl: string;
  enabled?: boolean;
};

export type EmailNotifyChannelConfig = {
  smtpHost: string;
  smtpPort: number;
  smtpTLS: boolean;
  username: string;
  password: string;
  senderAddress: string;
  receiverAddress: string;
  enabled?: boolean;
};

export type DingTalkNotifyChannelConfig = {
  accessToken: string;
  secret: string;
  enabled?: boolean;
};

export type GotifyNotifyChannelConfig = {
  url: string;
  token: string;
  priority: string;
  enabled?: boolean;
};

export type LarkNotifyChannelConfig = {
  webhookUrl: string;
  enabled?: boolean;
};

export type PushPlusNotifyChannelConfig = {
  token: string;
  enabled?: boolean;
};

export type ServerChanNotifyChannelConfig = {
  url: string;
  enabled?: boolean;
};

export type TelegramNotifyChannelConfig = {
  apiToken: string;
  chatId: string;
  enabled?: boolean;
};

export type WebhookNotifyChannelConfig = {
  url: string;
  enabled?: boolean;
};

export type WeComNotifyChannelConfig = {
  webhookUrl: string;
  enabled?: boolean;
};

export type NotifyChannel = {
  type: string;
  name: string;
};

export const notifyChannelsMap: Map<NotifyChannel["type"], NotifyChannel> = new Map(
  [
    [NOTIFY_CHANNELS.EMAIL, "common.notifier.email"],
    [NOTIFY_CHANNELS.DINGTALK, "common.notifier.dingtalk"],
    [NOTIFY_CHANNELS.GOTIFY, "common.notifier.gotify"],
    [NOTIFY_CHANNELS.LARK, "common.notifier.lark"],
    [NOTIFY_CHANNELS.PUSHPLUS, "common.notifier.pushplus"],
    [NOTIFY_CHANNELS.WECOM, "common.notifier.wecom"],
    [NOTIFY_CHANNELS.TELEGRAM, "common.notifier.telegram"],
    [NOTIFY_CHANNELS.SERVERCHAN, "common.notifier.serverchan"],
    [NOTIFY_CHANNELS.BARK, "common.notifier.bark"],
    [NOTIFY_CHANNELS.WEBHOOK, "common.notifier.webhook"],
  ].map(([type, name]) => [type, { type, name }])
);
// #endregion

// #region Settings: SSLProvider
export type SSLProviderSettingsContent = {
  provider: ApplyCAProviderType;
  config: {
    [key: string]: Record<string, unknown> | undefined;
  };
};
// #endregion

// #region Settings: Persistence
export type PersistenceSettingsContent = {
  workflowRunsMaxDaysRetention?: number;
  expiredCertificatesMaxDaysRetention?: number;
};
// #endregion
