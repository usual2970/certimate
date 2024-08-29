import { Deployment, Pahse } from "./deployment";

export type Domain = {
  id: string;
  domain: string;
  crontab: string;
  access: string;
  targetAccess: string;
  targetType: string;
  expiredAt?: string;
  phase?: Pahse;
  phaseSuccess?: boolean;
  lastDeployedAt?: string;
  enabled?: boolean;
  created?: string;
  updated?: string;
  deleted?: string;
  rightnow?: boolean;
  certificate?: string;
  privateKey?: string;
  expand?: {
    lastDeployment?: Deployment;
  };
};

export const getLastDeployment = (domain: Domain): Deployment | undefined => {
  return domain.expand?.lastDeployment;
};

export const targetTypeMap: Map<string, [string, string]> = new Map([
  ["aliyun-cdn", ["阿里云-CDN", "/imgs/providers/aliyun.svg"]],
  ["aliyun-oss", ["阿里云-OSS", "/imgs/providers/aliyun.svg"]],
  ["tencent-cdn", ["腾讯云-CDN", "/imgs/providers/tencent.svg"]],
  ["ssh", ["SSH部署", "/imgs/providers/ssh.png"]],
  ["webhook", ["Webhook", "/imgs/providers/webhook.svg"]],
]);

export const targetTypeKeys = Array.from(targetTypeMap.keys());
