import { type AccessUsageType } from "./provider";

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
