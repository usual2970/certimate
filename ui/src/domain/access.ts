/*
  注意：如果追加新的常量值，请保持以 ASCII 排序。
  NOTICE: If you add new constant, please keep ASCII order.
 */
export const ACCESS_PROVIDER_ACMEHTTPREQ = "acmehttpreq" as const;
export const ACCESS_PROVIDER_ALIYUN = "aliyun" as const;
export const ACCESS_PROVIDER_AWS = "aws" as const;
export const ACCESS_PROVIDER_BAIDUCLOUD = "baiducloud" as const;
export const ACCESS_PROVIDER_BYTEPLUS = "byteplus" as const;
export const ACCESS_PROVIDER_CLOUDFLARE = "cloudflare" as const;
export const ACCESS_PROVIDER_DOGECLOUD = "dogecloud" as const;
export const ACCESS_PROVIDER_GODADDY = "godaddy" as const;
export const ACCESS_PROVIDER_HUAWEICLOUD = "huaweicloud" as const;
export const ACCESS_PROVIDER_KUBERNETES = "k8s" as const;
export const ACCESS_PROVIDER_LOCAL = "local" as const;
export const ACCESS_PROVIDER_NAMEDOTCOM = "namedotcom" as const;
export const ACCESS_PROVIDER_NAMESILO = "namesilo" as const;
export const ACCESS_PROVIDER_POWERDNS = "powerdns" as const;
export const ACCESS_PROVIDER_QINIU = "qiniu" as const;
export const ACCESS_PROVIDER_SSH = "ssh" as const;
export const ACCESS_PROVIDER_TENCENTCLOUD = "tencentcloud" as const;
export const ACCESS_PROVIDER_VOLCENGINE = "volcengine" as const;
export const ACCESS_PROVIDER_WEBHOOK = "webhook" as const;
export const ACCESS_PROVIDERS = Object.freeze({
  ACMEHTTPREQ: ACCESS_PROVIDER_ACMEHTTPREQ,
  ALIYUN: ACCESS_PROVIDER_ALIYUN,
  AWS: ACCESS_PROVIDER_AWS,
  BAIDUCLOUD: ACCESS_PROVIDER_BAIDUCLOUD,
  BYTEPLUS: ACCESS_PROVIDER_BYTEPLUS,
  CLOUDFLARE: ACCESS_PROVIDER_CLOUDFLARE,
  DOGECLOUD: ACCESS_PROVIDER_DOGECLOUD,
  GODADDY: ACCESS_PROVIDER_GODADDY,
  HUAWEICLOUD: ACCESS_PROVIDER_HUAWEICLOUD,
  KUBERNETES: ACCESS_PROVIDER_KUBERNETES,
  LOCAL: ACCESS_PROVIDER_LOCAL,
  NAMEDOTCOM: ACCESS_PROVIDER_NAMEDOTCOM,
  NAMESILO: ACCESS_PROVIDER_NAMESILO,
  POWERDNS: ACCESS_PROVIDER_POWERDNS,
  QINIU: ACCESS_PROVIDER_QINIU,
  SSH: ACCESS_PROVIDER_SSH,
  TENCENTCLOUD: ACCESS_PROVIDER_TENCENTCLOUD,
  VOLCENGINE: ACCESS_PROVIDER_VOLCENGINE,
  WEBHOOK: ACCESS_PROVIDER_WEBHOOK,
} as const);

export type AccessProviderType = (typeof ACCESS_PROVIDERS)[keyof typeof ACCESS_PROVIDERS];

export const ACCESS_USAGE_ALL = "all" as const;
export const ACCESS_USAGE_APPLY = "apply" as const;
export const ACCESS_USAGE_DEPLOY = "deploy" as const;
export const ACCESS_USAGES = Object.freeze({
  ALL: ACCESS_USAGE_ALL,
  APPLY: ACCESS_USAGE_APPLY,
  DEPLOY: ACCESS_USAGE_DEPLOY,
} as const);

export type AccessUsageType = (typeof ACCESS_USAGES)[keyof typeof ACCESS_USAGES];

// #region AccessModel
export interface AccessModel extends BaseModel {
  name: string;
  configType: string;
  usage: AccessUsageType;
  config: /*
    注意：如果追加新的类型，请保持以 ASCII 排序。
    NOTICE: If you add new type, please keep ASCII order.
  */ Record<string, unknown> &
    (
      | ACMEHttpReqAccessConfig
      | AliyunAccessConfig
      | AWSAccessConfig
      | BaiduCloudAccessConfig
      | BytePlusAccessConfig
      | CloudflareAccessConfig
      | DogeCloudAccessConfig
      | GoDaddyAccessConfig
      | HuaweiCloudAccessConfig
      | KubernetesAccessConfig
      | LocalAccessConfig
      | NameDotComAccessConfig
      | NameSiloAccessConfig
      | PowerDNSAccessConfig
      | QiniuAccessConfig
      | SSHAccessConfig
      | TencentCloudAccessConfig
      | VolcEngineAccessConfig
      | WebhookAccessConfig
    );
}
// #endregion

// #region AccessConfig
export type ACMEHttpReqAccessConfig = {
  endpoint: string;
  mode?: string;
  username?: string;
  password?: string;
};

export type AliyunAccessConfig = {
  accessKeyId: string;
  accessKeySecret: string;
};

export type AWSAccessConfig = {
  accessKeyId: string;
  secretAccessKey: string;
  region?: string;
  hostedZoneId?: string;
};

export type BaiduCloudAccessConfig = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type BytePlusAccessConfig = {
  accessKey: string;
  secretKey: string;
};

export type CloudflareAccessConfig = {
  dnsApiToken: string;
};

export type DogeCloudAccessConfig = {
  accessKey: string;
  secretKey: string;
};

export type GoDaddyAccessConfig = {
  apiKey: string;
  apiSecret: string;
};

export type HuaweiCloudAccessConfig = {
  accessKeyId: string;
  secretAccessKey: string;
  region?: string;
};

export type KubernetesAccessConfig = {
  kubeConfig?: string;
};

export type LocalAccessConfig = NonNullable<unknown>;

export type NameDotComAccessConfig = {
  username: string;
  apiToken: string;
};

export type NameSiloAccessConfig = {
  apiKey: string;
};

export type PowerDNSAccessConfig = {
  apiUrl: string;
  apiKey: string;
};

export type QiniuAccessConfig = {
  accessKey: string;
  secretKey: string;
};

export type SSHAccessConfig = {
  host: string;
  port: number;
  username: string;
  password?: string;
  key?: string;
  keyPassphrase?: string;
};

export type TencentCloudAccessConfig = {
  secretId: string;
  secretKey: string;
};

export type VolcEngineAccessConfig = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type WebhookAccessConfig = {
  url: string;
};
// #endregion
