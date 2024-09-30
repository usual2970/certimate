import { z } from "zod";

export const accessTypeMap: Map<string, [string, string]> = new Map([
  ["tencent", ["tencent", "/imgs/providers/tencent.svg"]],
  ["aliyun", ["aliyun", "/imgs/providers/aliyun.svg"]],
  ["cloudflare", ["cloudflare", "/imgs/providers/cloudflare.svg"]],
  ["namesilo", ["namesilo", "/imgs/providers/namesilo.svg"]],
  ["godaddy", ["go.daddy", "/imgs/providers/godaddy.svg"]],
  ["qiniu", ["qiniu", "/imgs/providers/qiniu.svg"]],
  ["ssh", ["ssh", "/imgs/providers/ssh.svg"]],
  ["webhook", ["webhook", "/imgs/providers/webhook.svg"]],
  ["local", ["local.deployment", "/imgs/providers/local.svg"]],
]);

export const getProviderInfo = (t: string) => {
  return accessTypeMap.get(t);
};

export const accessFormType = z.union(
  [
    z.literal("aliyun"),
    z.literal("tencent"),
    z.literal("ssh"),
    z.literal("webhook"),
    z.literal("cloudflare"),
    z.literal("qiniu"),
    z.literal("namesilo"),
    z.literal("godaddy"),
    z.literal("local"),
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
    | TencentConfig
    | AliyunConfig
    | SSHConfig
    | WebhookConfig
    | CloudflareConfig
    | QiniuConfig
    | NamesiloConfig
    | GodaddyConfig
    | LocalConfig;

  deleted?: string;
  created?: string;
  updated?: string;
};

export type QiniuConfig = {
  accessKey: string;
  secretKey: string;
};

export type WebhookConfig = {
  url: string;
};

export type CloudflareConfig = {
  dnsApiToken: string;
};

export type TencentConfig = {
  secretId: string;
  secretKey: string;
};

export type AliyunConfig = {
  accessKeyId: string;
  accessKeySecret: string;
};

export type NamesiloConfig = {
  apiKey: string;
};
export type GodaddyConfig = {
  apiKey: string;
  apiSecret: string;
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

export type LocalConfig = {
  command: string;
  certPath: string;
  keyPath: string;
};

export const getUsageByConfigType = (configType: string): AccessUsage => {
  switch (configType) {
    case "aliyun":
    case "tencent":
      return "all";
    case "ssh":
    case "webhook":
    case "qiniu":
    case "local":
      return "deploy";

    case "cloudflare":
    case "namesilo":
    case "godaddy":
      return "apply";
    default:
      return "all";
  }
};
