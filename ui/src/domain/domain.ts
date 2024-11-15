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
  disableFollowCNAME?: boolean;
};

export type Statistic = {
  total: number;
  expired: number;
  enabled: number;
  disabled: number;
};

export type DeployTarget = {
  type: string;
  provider: string;
  name: string;
  icon: string;
};

export const deployTargetList: string[][] = [
  ["aliyun-oss", "common.provider.aliyun.oss", "/imgs/providers/aliyun.svg"],
  ["aliyun-cdn", "common.provider.aliyun.cdn", "/imgs/providers/aliyun.svg"],
  ["aliyun-dcdn", "common.provider.aliyun.dcdn", "/imgs/providers/aliyun.svg"],
  ["aliyun-clb", "common.provider.aliyun.clb", "/imgs/providers/aliyun.svg"],
  ["aliyun-alb", "common.provider.aliyun.alb", "/imgs/providers/aliyun.svg"],
  ["aliyun-nlb", "common.provider.aliyun.nlb", "/imgs/providers/aliyun.svg"],
  ["tencent-cdn", "common.provider.tencent.cdn", "/imgs/providers/tencent.svg"],
  ["tencent-ecdn", "common.provider.tencent.ecdn", "/imgs/providers/tencent.svg"],
  ["tencent-clb", "common.provider.tencent.clb", "/imgs/providers/tencent.svg"],
  ["tencent-cos", "common.provider.tencent.cos", "/imgs/providers/tencent.svg"],
  ["tencent-teo", "common.provider.tencent.teo", "/imgs/providers/tencent.svg"],
  ["huaweicloud-cdn", "common.provider.huaweicloud.cdn", "/imgs/providers/huaweicloud.svg"],
  ["huaweicloud-elb", "common.provider.huaweicloud.elb", "/imgs/providers/huaweicloud.svg"],
  ["baiducloud-cdn", "common.provider.baiducloud.cdn", "/imgs/providers/baiducloud.svg"],
  ["qiniu-cdn", "common.provider.qiniu.cdn", "/imgs/providers/qiniu.svg"],
  ["dogecloud-cdn", "common.provider.dogecloud.cdn", "/imgs/providers/dogecloud.svg"],
  ["local", "common.provider.local", "/imgs/providers/local.svg"],
  ["ssh", "common.provider.ssh", "/imgs/providers/ssh.svg"],
  ["webhook", "common.provider.webhook", "/imgs/providers/webhook.svg"],
  ["k8s-secret", "common.provider.kubernetes.secret", "/imgs/providers/k8s.svg"],
  ["volcengine-live", "common.provider.volcengine.live", "/imgs/providers/volcengine.svg"],
  ["volcengine-cdn", "common.provider.volcengine.cdn", "/imgs/providers/volcengine.svg"],
];

export const deployTargetsMap: Map<DeployTarget["type"], DeployTarget> = new Map(
  deployTargetList.map(([type, name, icon]) => [type, { type, provider: type.split("-")[0], name, icon }])
);

export const deployTargets = deployTargetList.map(([type, name, icon]) => ({ type, provider: type.split("-")[0], name, icon }));
