import { z } from "zod";

export const accessTypeMap: Map<string, [string, string]> = new Map([
  ["aliyun", ["aliyun", "/imgs/providers/aliyun.svg"]],
  ["tencent", ["tencent", "/imgs/providers/tencent.svg"]],
  ["huaweicloud", ["huaweicloud", "/imgs/providers/huaweicloud.svg"]],
  ["qiniu", ["qiniu", "/imgs/providers/qiniu.svg"]],
  ["cloudflare", ["cloudflare", "/imgs/providers/cloudflare.svg"]],
  ["namesilo", ["namesilo", "/imgs/providers/namesilo.svg"]],
  ["godaddy", ["go.daddy", "/imgs/providers/godaddy.svg"]],
  ["local", ["local.deployment", "/imgs/providers/local.svg"]],
  ["ssh", ["ssh", "/imgs/providers/ssh.svg"]],
  ["webhook", ["webhook", "/imgs/providers/webhook.svg"]],
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
    z.literal("cloudflare"),
    z.literal("namesilo"),
    z.literal("godaddy"),
    z.literal("local"),
    z.literal("ssh"),
    z.literal("webhook"),
  ],
  { message: "access.not.empty" }
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
    | CloudflareConfig
    | NamesiloConfig
    | GodaddyConfig
    | LocalConfig
    | SSHConfig
    | WebhookConfig;
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

export type LocalConfig = {
  command: string;
  certPath: string;
  keyPath: string;
};

export type SSHConfig = {
  host: string;
  port: string;
  preCommand?: string;
  command: string;
  username: string;
  password?: string;
  key?: string;
  keyFile?: string;
  certPath: string;
  keyPath: string;
};

export type WebhookConfig = {
  url: string;
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
      return "deploy";

    case "cloudflare":
    case "namesilo":
    case "godaddy":
      return "apply";

    default:
      return "all";
  }
};
