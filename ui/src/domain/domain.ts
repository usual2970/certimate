import { Deployment, Pahse } from "./deployment";

export type Domain = {
  id: string;
  domain: string;
  email?: string;
  crontab: string;
  access: string;
  targetAccess?: string;
  targetType: string;
  expiredAt?: string;
  phase?: Pahse;
  phaseSuccess?: boolean;
  lastDeployedAt?: string;
  variables?: string;
  nameservers?: string;
  group?: string;
  enabled?: boolean;
  deployed?: boolean;
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

export type Statistic = {
  total: number;
  expired: number;
  enabled: number;
  disabled: number;
};

export const getLastDeployment = (domain: Domain): Deployment | undefined => {
  return domain.expand?.lastDeployment;
};

export const targetTypeMap: Map<string, [string, string]> = new Map([
  ["aliyun-cdn", ["aliyun.cdn", "/imgs/providers/aliyun.svg"]],
  ["aliyun-oss", ["aliyun.oss", "/imgs/providers/aliyun.svg"]],
  ["aliyun-dcdn", ["aliyun.dcdn", "/imgs/providers/aliyun.svg"]],
  ["tencent-cdn", ["tencent.cdn", "/imgs/providers/tencent.svg"]],
  ["ssh", ["ssh", "/imgs/providers/ssh.svg"]],
  ["qiniu-cdn", ["qiniu.cdn", "/imgs/providers/qiniu.svg"]],
  ["webhook", ["webhook", "/imgs/providers/webhook.svg"]],
  ["local", ["local", "/imgs/providers/local.svg"]],
]);

export const targetTypeKeys = Array.from(targetTypeMap.keys());
