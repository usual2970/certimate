import { z } from "zod";
import { type BaseModel } from "pocketbase";

type AccessUsages = "apply" | "deploy" | "all";

type AccessProvider = {
  type: string;
  name: string;
  icon: string;
  usage: AccessUsages;
};

export const accessProvidersMap: Map<AccessProvider["type"], AccessProvider> = new Map(
  [
    ["aliyun", "common.provider.aliyun", "/imgs/providers/aliyun.svg", "all"],
    ["tencent", "common.provider.tencent", "/imgs/providers/tencent.svg", "all"],
    ["huaweicloud", "common.provider.huaweicloud", "/imgs/providers/huaweicloud.svg", "all"],
    ["baiducloud", "common.provider.baiducloud", "/imgs/providers/baiducloud.svg", "all"],
    ["qiniu", "common.provider.qiniu", "/imgs/providers/qiniu.svg", "deploy"],
    ["dogecloud", "common.provider.dogecloud", "/imgs/providers/dogecloud.svg", "deploy"],
    ["volcengine", "common.provider.volcengine", "/imgs/providers/volcengine.svg", "all"],
    ["byteplus", "common.provider.byteplus", "/imgs/providers/byteplus.svg", "all"],
    ["aws", "common.provider.aws", "/imgs/providers/aws.svg", "apply"],
    ["cloudflare", "common.provider.cloudflare", "/imgs/providers/cloudflare.svg", "apply"],
    ["namesilo", "common.provider.namesilo", "/imgs/providers/namesilo.svg", "apply"],
    ["godaddy", "common.provider.godaddy", "/imgs/providers/godaddy.svg", "apply"],
    ["pdns", "common.provider.pdns", "/imgs/providers/pdns.svg", "apply"],
    ["httpreq", "common.provider.httpreq", "/imgs/providers/httpreq.svg", "apply"],
    ["local", "common.provider.local", "/imgs/providers/local.svg", "deploy"],
    ["ssh", "common.provider.ssh", "/imgs/providers/ssh.svg", "deploy"],
    ["webhook", "common.provider.webhook", "/imgs/providers/webhook.svg", "deploy"],
    ["k8s", "common.provider.kubernetes", "/imgs/providers/k8s.svg", "deploy"],
  ].map(([type, name, icon, usage]) => [type, { type, name, icon, usage: usage as AccessUsages }])
);

export const accessTypeFormSchema = z.union(
  [
    z.literal("aliyun"),
    z.literal("tencent"),
    z.literal("huaweicloud"),
    z.literal("baiducloud"),
    z.literal("qiniu"),
    z.literal("dogecloud"),
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
    z.literal("volcengine"),
    z.literal("byteplus"),
  ],
  { message: "access.authorization.form.type.placeholder" }
);

export interface AccessModel extends Omit<BaseModel, "created" | "updated"> {
  name: string;
  configType: string;
  usage: AccessUsages;
  group?: string;
  config:
    | AliyunConfig
    | TencentConfig
    | HuaweiCloudConfig
    | QiniuConfig
    | DogeCloudConfig
    | AwsConfig
    | CloudflareConfig
    | NamesiloConfig
    | GodaddyConfig
    | PdnsConfig
    | HttpreqConfig
    | LocalConfig
    | SSHConfig
    | WebhookConfig
    | KubernetesConfig
    | VolcengineConfig
    | ByteplusConfig;
}

export type AliyunConfig = {
  accessKeyId: string;
  accessKeySecret: string;
};

export type TencentConfig = {
  secretId: string;
  secretKey: string;
};

export type HuaweiCloudConfig = {
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
};

export type BaiduCloudConfig = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type QiniuConfig = {
  accessKey: string;
  secretKey: string;
};

export type DogeCloudConfig = {
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

export type VolcengineConfig = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type ByteplusConfig = {
  accessKey: string;
  secretKey: string;
};
