import { z } from "zod";

export const accessTypeMap: Map<string, [string, string]> = new Map([
  ["aliyun", ["common.provider.aliyun", "/imgs/providers/aliyun.svg"]],
  ["tencent", ["common.provider.tencent", "/imgs/providers/tencent.svg"]],
  ["huaweicloud", ["common.provider.huaweicloud", "/imgs/providers/huaweicloud.svg"]],
  ["qiniu", ["common.provider.qiniu", "/imgs/providers/qiniu.svg"]],
  ["aws", ["common.provider.aws", "/imgs/providers/aws.svg"]],
  ["cloudflare", ["common.provider.cloudflare", "/imgs/providers/cloudflare.svg"]],
  ["namesilo", ["common.provider.namesilo", "/imgs/providers/namesilo.svg"]],
  ["godaddy", ["common.provider.godaddy", "/imgs/providers/godaddy.svg"]],
  ["pdns", ["common.provider.pdns", "/imgs/providers/pdns.svg"]],
  ["httpreq", ["common.provider.httpreq", "/imgs/providers/httpreq.svg"]],
  ["local", ["common.provider.local", "/imgs/providers/local.svg"]],
  ["ssh", ["common.provider.ssh", "/imgs/providers/ssh.svg"]],
  ["webhook", ["common.provider.webhook", "/imgs/providers/webhook.svg"]],
  ["k8s", ["common.provider.kubernetes", "/imgs/providers/k8s.svg"]],
]);

export const getProviderInfo = (t: string) => {
  return accessTypeMap.get(t);
};

export const accessFormType = z.union(
  [
    z.literal("aliyun"),
    z.literal("tencent"),
    z.literal("huaweicloud"),
    z.literal("qiniu"),
    z.literal("aws"),
    z.literal("cloudflare"),
    z.literal("namesilo"),
    z.literal("godaddy"),
    z.literal("pdns"),
    z.literal("httpreq"),
    z.literal("local"),
    z.literal("ssh"),
    z.literal("webhook"),
    z.literal("k8s"),
  ],
  { message: "access.authorization.form.type.placeholder" }
);

type AccessUsage = "apply" | "deploy" | "all";

export type Access = {
  id: string;
  name: string;
  configType: string;
  usage: AccessUsage;
  group?: string;
  config:
    | AliyunConfig
    | TencentConfig
    | HuaweicloudConfig
    | QiniuConfig
    | AwsConfig
    | CloudflareConfig
    | NamesiloConfig
    | GodaddyConfig
    | PdnsConfig
    | HttpreqConfig
    | LocalConfig
    | SSHConfig
    | WebhookConfig
    | KubernetesConfig;
  deleted?: string;
  created?: string;
  updated?: string;
};

export type AliyunConfig = {
  accessKeyId: string;
  accessKeySecret: string;
};

export type TencentConfig = {
  secretId: string;
  secretKey: string;
};

export type HuaweicloudConfig = {
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
};

export type QiniuConfig = {
  accessKey: string;
  secretKey: string;
};

export type AwsConfig = {
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
  hostedZoneId?: string;
};

export type CloudflareConfig = {
  dnsApiToken: string;
};

export type NamesiloConfig = {
  apiKey: string;
};

export type GodaddyConfig = {
  apiKey: string;
  apiSecret: string;
};

export type PdnsConfig = {
  apiUrl: string;
  apiKey: string;
};

export type HttpreqConfig = {
  endpoint: string;
  mode: string;
  username: string;
  password: string;
};

export type LocalConfig = Record<string, string>;

export type SSHConfig = {
  host: string;
  port: string;
  username: string;
  password?: string;
  key?: string;
  keyFile?: string;
  keyPassphrase?: string;
};

export type WebhookConfig = {
  url: string;
};

export type KubernetesConfig = {
  kubeConfig: string;
};

export const getUsageByConfigType = (configType: string): AccessUsage => {
  switch (configType) {
    case "aliyun":
    case "tencent":
    case "huaweicloud":
      return "all";

    case "qiniu":
    case "local":
    case "ssh":
    case "webhook":
    case "k8s":
      return "deploy";

    case "aws":
    case "cloudflare":
    case "namesilo":
    case "godaddy":
    case "pdns":
    case "httpreq":
      return "apply";

    default:
      return "all";
  }
};
