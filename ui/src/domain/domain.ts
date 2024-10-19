import { Deployment, Pahse } from "./deployment";

export type Domain = {
  id?: string;
  domain: string;
  email?: string;
  crontab: string;
  access: string;
  targetAccess?: string;
  targetType?: string;
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

  applyConfig?: ApplyConfig;
  deployConfig?: DeployConfig[];
};

export type KVType = {
  key: string;
  value: string;
};

export type DeployConfig = {
  id?: string;
  access: string;
  type: string;
  config?: {
    [key: string]: string;
  } & {
    variables?: KVType[];
  };
};

export type ApplyConfig = {
  access: string;
  email: string;
  keyAlgorithm?: string;
  nameservers?: string;
  timeout?: number;
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
  ["aliyun-oss", ["common.provider.aliyun.oss", "/imgs/providers/aliyun.svg"]],
  ["aliyun-cdn", ["common.provider.aliyun.cdn", "/imgs/providers/aliyun.svg"]],
  ["aliyun-dcdn", ["common.provider.aliyun.dcdn", "/imgs/providers/aliyun.svg"]],
  ["tencent-cdn", ["common.provider.tencent.cdn", "/imgs/providers/tencent.svg"]],
  ["huaweicloud-cdn", ["common.provider.huaweicloud.cdn", "/imgs/providers/huaweicloud.svg"]],
  ["qiniu-cdn", ["common.provider.qiniu.cdn", "/imgs/providers/qiniu.svg"]],
  ["local", ["common.provider.local", "/imgs/providers/local.svg"]],
  ["ssh", ["common.provider.ssh", "/imgs/providers/ssh.svg"]],
  ["webhook", ["common.provider.webhook", "/imgs/providers/webhook.svg"]],
  ["k8s-secret", ["common.provider.kubernetes.secret", "/imgs/providers/k8s.svg"]],
]);

export const targetTypeKeys = Array.from(targetTypeMap.keys());
