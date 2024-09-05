import { z } from "zod";

export const accessTypeMap: Map<string, [string, string]> = new Map([
  ["tencent", ["腾讯云", "/imgs/providers/tencent.svg"]],
  ["aliyun", ["阿里云", "/imgs/providers/aliyun.svg"]],
  ["cloudflare", ["Cloudflare", "/imgs/providers/cloudflare.svg"]],
  ["qiniu", ["七牛云", "/imgs/providers/qiniu.svg"]],
  ["ssh", ["SSH部署", "/imgs/providers/ssh.svg"]],
  ["webhook", ["Webhook", "/imgs/providers/webhook.svg"]],
]);

export const accessFormType = z.union(
  [
    z.literal("aliyun"),
    z.literal("tencent"),
    z.literal("ssh"),
    z.literal("webhook"),
    z.literal("cloudflare"),
    z.literal("qiniu"),
  ],
  { message: "请选择云服务商" }
);

export type Access = {
  id: string;
  name: string;
  configType: string;
  config:
    | TencentConfig
    | AliyunConfig
    | SSHConfig
    | WebhookConfig
    | CloudflareConfig
    | QiniuConfig;
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

export type SSHConfig = {
  host: string;
  port: string;
  command: string;
  username: string;
  password?: string;
  key?: string;
  keyFile?: string;
  certPath: string;
  keyPath: string;
};
