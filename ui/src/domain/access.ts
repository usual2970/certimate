import { z } from "zod";

export const accessTypeMap: Map<string, [string, string]> = new Map([
  ["aliyun", ["阿里云", "/imgs/providers/aliyun.svg"]],
  ["tencent", ["腾讯云", "/imgs/providers/tencent.svg"]],
  ["aws", ["AWS", "/imgs/providers/aws.svg"]],
  ["cloudflare", ["Cloudflare", "/imgs/providers/cloudflare.svg"]],
  ["namesilo", ["Namesilo", "/imgs/providers/namesilo.svg"]],
  ["godaddy", ["GoDaddy", "/imgs/providers/godaddy.svg"]],
  ["qiniu", ["七牛云", "/imgs/providers/qiniu.svg"]],
  ["ssh", ["SSH部署", "/imgs/providers/ssh.svg"]],
  ["webhook", ["Webhook", "/imgs/providers/webhook.svg"]],
]);

export const getProviderInfo = (t: string) => {
  return accessTypeMap.get(t);
};

export const accessFormType = z.union(
  [
    z.literal("aliyun"),
    z.literal("tencent"),
    z.literal("aws"),
    z.literal("ssh"),
    z.literal("webhook"),
    z.literal("cloudflare"),
    z.literal("qiniu"),
    z.literal("namesilo"),
    z.literal("godaddy"),
  ],
  { message: "请选择云服务商" }
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
    | AWSConfig
    | SSHConfig
    | WebhookConfig
    | CloudflareConfig
    | QiniuConfig
    | NamesiloConfig
    | GodaddyConfig;

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

export type AWSConfig = {
  accessKeyId: string;
  secretAccessKey: string;
  sslProvider: string; //sslprovider
};

export type NamesiloConfig = {
  apiKey: string;
};
export type GodaddyConfig = {
  apiKey: string;
  apiSecret: string;
};

export type SSHConfig = {
  host: string[];
  port: string;
  command: string;
  username: string;
  password?: string;
  key?: string;
  keyFile?: string;
  certPath: string;
  keyPath: string;
};

export const getUsageByConfigType = (configType: string): AccessUsage => {
  switch (configType) {
    case "aliyun":
    case "tencent":
    case "aws":
      return "all";
    case "ssh":
    case "webhook":
    case "qiniu":
      return "deploy";

    case "cloudflare":
    case "namesilo":
    case "godaddy":
      return "apply";
    default:
      return "all";
  }
};