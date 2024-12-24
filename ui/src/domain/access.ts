/*
  注意：如果追加新的常量值，请保持以 ASCII 排序。
  NOTICE: If you add new constant, please keep ASCII order.
 */
export const ACCESS_PROVIDER_TYPE_ACMEHTTPREQ = "acmehttpreq" as const;
export const ACCESS_PROVIDER_TYPE_ALIYUN = "aliyun" as const;
export const ACCESS_PROVIDER_TYPE_AWS = "aws" as const;
export const ACCESS_PROVIDER_TYPE_BAIDUCLOUD = "baiducloud" as const;
export const ACCESS_PROVIDER_TYPE_BYTEPLUS = "byteplus" as const;
export const ACCESS_PROVIDER_TYPE_CLOUDFLARE = "cloudflare" as const;
export const ACCESS_PROVIDER_TYPE_DOGECLOUD = "dogecloud" as const;
export const ACCESS_PROVIDER_TYPE_GODADDY = "godaddy" as const;
export const ACCESS_PROVIDER_TYPE_HUAWEICLOUD = "huaweicloud" as const;
export const ACCESS_PROVIDER_TYPE_KUBERNETES = "k8s" as const;
export const ACCESS_PROVIDER_TYPE_LOCAL = "local" as const;
export const ACCESS_PROVIDER_TYPE_NAMEDOTCOM = "namedotcom" as const;
export const ACCESS_PROVIDER_TYPE_NAMESILO = "namesilo" as const;
export const ACCESS_PROVIDER_TYPE_POWERDNS = "powerdns" as const;
export const ACCESS_PROVIDER_TYPE_QINIU = "qiniu" as const;
export const ACCESS_PROVIDER_TYPE_SSH = "ssh" as const;
export const ACCESS_PROVIDER_TYPE_TENCENTCLOUD = "tencentcloud" as const;
export const ACCESS_PROVIDER_TYPE_VOLCENGINE = "volcengine" as const;
export const ACCESS_PROVIDER_TYPE_WEBHOOK = "webhook" as const;
export const ACCESS_PROVIDER_TYPES = Object.freeze({
  ACMEHTTPREQ: ACCESS_PROVIDER_TYPE_ACMEHTTPREQ,
  ALIYUN: ACCESS_PROVIDER_TYPE_ALIYUN,
  AWS: ACCESS_PROVIDER_TYPE_AWS,
  BAIDUCLOUD: ACCESS_PROVIDER_TYPE_BAIDUCLOUD,
  BYTEPLUS: ACCESS_PROVIDER_TYPE_BYTEPLUS,
  CLOUDFLARE: ACCESS_PROVIDER_TYPE_CLOUDFLARE,
  DOGECLOUD: ACCESS_PROVIDER_TYPE_DOGECLOUD,
  GODADDY: ACCESS_PROVIDER_TYPE_GODADDY,
  HUAWEICLOUD: ACCESS_PROVIDER_TYPE_HUAWEICLOUD,
  KUBERNETES: ACCESS_PROVIDER_TYPE_KUBERNETES,
  LOCAL: ACCESS_PROVIDER_TYPE_LOCAL,
  NAMEDOTCOM: ACCESS_PROVIDER_TYPE_NAMEDOTCOM,
  NAMESILO: ACCESS_PROVIDER_TYPE_NAMESILO,
  POWERDNS: ACCESS_PROVIDER_TYPE_POWERDNS,
  QINIU: ACCESS_PROVIDER_TYPE_QINIU,
  SSH: ACCESS_PROVIDER_TYPE_SSH,
  TENCENTCLOUD: ACCESS_PROVIDER_TYPE_TENCENTCLOUD,
  VOLCENGINE: ACCESS_PROVIDER_TYPE_VOLCENGINE,
  WEBHOOK: ACCESS_PROVIDER_TYPE_WEBHOOK,
} as const);

export type AccessProviderTypes = (typeof ACCESS_PROVIDER_TYPES)[keyof typeof ACCESS_PROVIDER_TYPES];

export const ACCESS_PROVIDER_USAGE_ALL = "all" as const;
export const ACCESS_PROVIDER_USAGE_APPLY = "apply" as const;
export const ACCESS_PROVIDER_USAGE_DEPLOY = "deploy" as const;
export const ACCESS_PROVIDER_USAGES = Object.freeze({
  ALL: ACCESS_PROVIDER_USAGE_ALL,
  APPLY: ACCESS_PROVIDER_USAGE_APPLY,
  DEPLOY: ACCESS_PROVIDER_USAGE_DEPLOY,
} as const);

export type AccessProviderUsages = (typeof ACCESS_PROVIDER_USAGES)[keyof typeof ACCESS_PROVIDER_USAGES];

// #region AccessModel
export interface AccessModel extends BaseModel {
  name: string;
  configType: string;
  usage: AccessProviderUsages;
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

// #region AccessProvider
export type AccessProvider = {
  type: string;
  name: string;
  icon: string;
  usage: AccessProviderUsages;
};

export const accessProvidersMap: Map<AccessProvider["type"], AccessProvider> = new Map(
  /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
  [
    [ACCESS_PROVIDER_TYPE_ALIYUN, "common.provider.aliyun", "/imgs/providers/aliyun.svg", "all"],
    [ACCESS_PROVIDER_TYPE_TENCENTCLOUD, "common.provider.tencentcloud", "/imgs/providers/tencentcloud.svg", "all"],
    [ACCESS_PROVIDER_TYPE_HUAWEICLOUD, "common.provider.huaweicloud", "/imgs/providers/huaweicloud.svg", "all"],
    [ACCESS_PROVIDER_TYPE_BAIDUCLOUD, "common.provider.baiducloud", "/imgs/providers/baiducloud.svg", "all"],
    [ACCESS_PROVIDER_TYPE_QINIU, "common.provider.qiniu", "/imgs/providers/qiniu.svg", "deploy"],
    [ACCESS_PROVIDER_TYPE_DOGECLOUD, "common.provider.dogecloud", "/imgs/providers/dogecloud.svg", "deploy"],
    [ACCESS_PROVIDER_TYPE_VOLCENGINE, "common.provider.volcengine", "/imgs/providers/volcengine.svg", "all"],
    [ACCESS_PROVIDER_TYPE_BYTEPLUS, "common.provider.byteplus", "/imgs/providers/byteplus.svg", "all"],
    [ACCESS_PROVIDER_TYPE_AWS, "common.provider.aws", "/imgs/providers/aws.svg", "apply"],
    [ACCESS_PROVIDER_TYPE_CLOUDFLARE, "common.provider.cloudflare", "/imgs/providers/cloudflare.svg", "apply"],
    [ACCESS_PROVIDER_TYPE_NAMEDOTCOM, "common.provider.namedotcom", "/imgs/providers/namedotcom.svg", "apply"],
    [ACCESS_PROVIDER_TYPE_NAMESILO, "common.provider.namesilo", "/imgs/providers/namesilo.svg", "apply"],
    [ACCESS_PROVIDER_TYPE_GODADDY, "common.provider.godaddy", "/imgs/providers/godaddy.svg", "apply"],
    [ACCESS_PROVIDER_TYPE_POWERDNS, "common.provider.powerdns", "/imgs/providers/powerdns.svg", "apply"],
    [ACCESS_PROVIDER_TYPE_LOCAL, "common.provider.local", "/imgs/providers/local.svg", "deploy"],
    [ACCESS_PROVIDER_TYPE_SSH, "common.provider.ssh", "/imgs/providers/ssh.svg", "deploy"],
    [ACCESS_PROVIDER_TYPE_WEBHOOK, "common.provider.webhook", "/imgs/providers/webhook.svg", "deploy"],
    [ACCESS_PROVIDER_TYPE_KUBERNETES, "common.provider.kubernetes", "/imgs/providers/kubernetes.svg", "deploy"],
    [ACCESS_PROVIDER_TYPE_ACMEHTTPREQ, "common.provider.acmehttpreq", "/imgs/providers/acmehttpreq.svg", "apply"],
  ].map(([type, name, icon, usage]) => [type, { type, name, icon, usage: usage as AccessProviderUsages }])
);
// #endregion
