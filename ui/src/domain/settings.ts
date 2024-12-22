export const SETTINGS_NAME_EMAILS = "emails" as const;
export const SETTINGS_NAME_NOTIFYTEMPLATES = "notifyTemplates" as const;
export const SETTINGS_NAME_NOTIFYCHANNELS = "notifyChannels" as const;
export const SETTINGS_NAME_SSLPROVIDER = "sslProvider" as const;
export const SETTINGS_NAMES = Object.freeze({
  EMAILS: SETTINGS_NAME_EMAILS,
  NOTIFY_TEMPLATES: SETTINGS_NAME_NOTIFYTEMPLATES,
  NOTIFY_CHANNELS: SETTINGS_NAME_NOTIFYCHANNELS,
  SSL_PROVIDER: SETTINGS_NAME_SSLPROVIDER,
} as const);

export type SettingsNames = (typeof SETTINGS_NAMES)[keyof typeof SETTINGS_NAMES];

export interface SettingsModel<T extends NonNullable<unknown> = NonNullable<unknown>> extends BaseModel {
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
  subject: "您有 {COUNT} 张证书即将过期",
  message: "有 {COUNT} 张证书即将过期，域名分别为 {DOMAINS}，请保持关注！",
};
// #endregion

// #region Settings: NotifyChannels
export const NOTIFY_CHANNEL_BARK = "bark" as const;
export const NOTIFY_CHANNEL_DINGTALK = "dingtalk" as const;
export const NOTIFY_CHANNEL_EMAIL = "email" as const;
export const NOTIFY_CHANNEL_LARK = "lark" as const;
export const NOTIFY_CHANNEL_SERVERCHAN = "serverchan" as const;
export const NOTIFY_CHANNEL_TELEGRAM = "telegram" as const;
export const NOTIFY_CHANNEL_WEBHOOK = "webhook" as const;
export const NOTIFY_CHANNEL_WECOM = "wecom" as const;
export const NOTIFY_CHANNELS = Object.freeze({
  BARK: NOTIFY_CHANNEL_BARK,
  DINGTALK: NOTIFY_CHANNEL_DINGTALK,
  EMAIL: NOTIFY_CHANNEL_EMAIL,
  LARK: NOTIFY_CHANNEL_LARK,
  SERVERCHAN: NOTIFY_CHANNEL_SERVERCHAN,
  TELEGRAM: NOTIFY_CHANNEL_TELEGRAM,
  WEBHOOK: NOTIFY_CHANNEL_WEBHOOK,
  WECOM: NOTIFY_CHANNEL_WECOM,
} as const);

export type NotifyChannels = (typeof NOTIFY_CHANNELS)[keyof typeof NOTIFY_CHANNELS];

export type NotifyChannelsSettingsContent = {
  /*
    注意：如果追加新的类型，请保持以 ASCII 排序。
    NOTICE: If you add new type, please keep ASCII order.
  */
  [key: string]: ({ enabled?: boolean } & Record<string, unknown>) | undefined;
  [NOTIFY_CHANNEL_BARK]?: BarkNotifyChannelConfig;
  [NOTIFY_CHANNEL_DINGTALK]?: DingTalkNotifyChannelConfig;
  [NOTIFY_CHANNEL_EMAIL]?: EmailNotifyChannelConfig;
  [NOTIFY_CHANNEL_LARK]?: LarkNotifyChannelConfig;
  [NOTIFY_CHANNEL_SERVERCHAN]?: ServerChanNotifyChannelConfig;
  [NOTIFY_CHANNEL_TELEGRAM]?: TelegramNotifyChannelConfig;
  [NOTIFY_CHANNEL_WEBHOOK]?: WebhookNotifyChannelConfig;
  [NOTIFY_CHANNEL_WECOM]?: WeComNotifyChannelConfig;
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

export type LarkNotifyChannelConfig = {
  webhookUrl: string;
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
    ["email", "common.notifier.email"],
    ["dingtalk", "common.notifier.dingtalk"],
    ["lark", "common.notifier.lark"],
    ["wecom", "common.notifier.wecom"],
    ["telegram", "common.notifier.telegram"],
    ["serverchan", "common.notifier.serverchan"],
    ["bark", "common.notifier.bark"],
    ["webhook", "common.notifier.webhook"],
  ].map(([type, name]) => [type, { type, name }])
);
// #endregion

// #region Settings: SSLProvider
export const SSLPROVIDER_LETSENCRYPT = "letsencrypt" as const;
export const SSLPROVIDER_ZEROSSL = "zerossl" as const;
export const SSLPROVIDER_GOOGLETRUSTSERVICES = "gts" as const;
export const SSLPROVIDERS = Object.freeze({
  LETS_ENCRYPT: SSLPROVIDER_LETSENCRYPT,
  ZERO_SSL: SSLPROVIDER_ZEROSSL,
  GOOGLE_TRUST_SERVICES: SSLPROVIDER_GOOGLETRUSTSERVICES,
} as const);

export type SSLProviders = (typeof SSLPROVIDERS)[keyof typeof SSLPROVIDERS];

export type SSLProviderSettingsContent = {
  provider: (typeof SSLPROVIDERS)[keyof typeof SSLPROVIDERS];
  config: {
    [key: string]: Record<string, unknown> | undefined;
    letsencrypt?: SSLProviderLetsEncryptConfig;
    zerossl?: SSLProviderZeroSSLConfig;
    gts?: SSLProviderGoogleTrustServicesConfig;
  };
};

export type SSLProviderLetsEncryptConfig = NonNullable<unknown>;

export type SSLProviderZeroSSLConfig = {
  eabKid: string;
  eabHmacKey: string;
};

export type SSLProviderGoogleTrustServicesConfig = {
  eabKid: string;
  eabHmacKey: string;
};
// #endregion
