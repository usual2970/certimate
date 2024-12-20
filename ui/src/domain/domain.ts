import {
  ACCESS_PROVIDER_TYPE_ALIYUN,
  ACCESS_PROVIDER_TYPE_BAIDUCLOUD,
  ACCESS_PROVIDER_TYPE_BYTEPLUS,
  ACCESS_PROVIDER_TYPE_DOGECLOUD,
  ACCESS_PROVIDER_TYPE_HUAWEICLOUD,
  ACCESS_PROVIDER_TYPE_KUBERNETES,
  ACCESS_PROVIDER_TYPE_LOCAL,
  ACCESS_PROVIDER_TYPE_QINIU,
  ACCESS_PROVIDER_TYPE_SSH,
  ACCESS_PROVIDER_TYPE_TENCENTCLOUD,
  ACCESS_PROVIDER_TYPE_VOLCENGINE,
  ACCESS_PROVIDER_TYPE_WEBHOOK,
} from "./access";

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

export type DeployTarget = {
  type: string;
  provider: string;
  name: string;
  icon: string;
};

export const deployTargetList: string[][] = [
  /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
  [`${ACCESS_PROVIDER_TYPE_LOCAL}`, "common.provider.local", "/imgs/providers/local.svg"],
  [`${ACCESS_PROVIDER_TYPE_SSH}`, "common.provider.ssh", "/imgs/providers/ssh.svg"],
  [`${ACCESS_PROVIDER_TYPE_WEBHOOK}`, "common.provider.webhook", "/imgs/providers/webhook.svg"],
  [`${ACCESS_PROVIDER_TYPE_ALIYUN}-oss`, "common.provider.aliyun.oss", "/imgs/providers/aliyun.svg"],
  [`${ACCESS_PROVIDER_TYPE_ALIYUN}-cdn`, "common.provider.aliyun.cdn", "/imgs/providers/aliyun.svg"],
  [`${ACCESS_PROVIDER_TYPE_ALIYUN}-dcdn`, "common.provider.aliyun.dcdn", "/imgs/providers/aliyun.svg"],
  [`${ACCESS_PROVIDER_TYPE_ALIYUN}-clb`, "common.provider.aliyun.clb", "/imgs/providers/aliyun.svg"],
  [`${ACCESS_PROVIDER_TYPE_ALIYUN}-alb`, "common.provider.aliyun.alb", "/imgs/providers/aliyun.svg"],
  [`${ACCESS_PROVIDER_TYPE_ALIYUN}-nlb`, "common.provider.aliyun.nlb", "/imgs/providers/aliyun.svg"],
  [`${ACCESS_PROVIDER_TYPE_TENCENTCLOUD}-cdn`, "common.provider.tencentcloud.cdn", "/imgs/providers/tencentcloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_TENCENTCLOUD}-ecdn`, "common.provider.tencentcloud.ecdn", "/imgs/providers/tencentcloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_TENCENTCLOUD}-clb`, "common.provider.tencentcloud.clb", "/imgs/providers/tencentcloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_TENCENTCLOUD}-cos`, "common.provider.tencentcloud.cos", "/imgs/providers/tencentcloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_TENCENTCLOUD}-eo`, "common.provider.tencentcloud.eo", "/imgs/providers/tencentcloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_HUAWEICLOUD}-cdn`, "common.provider.huaweicloud.cdn", "/imgs/providers/huaweicloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_HUAWEICLOUD}-elb`, "common.provider.huaweicloud.elb", "/imgs/providers/huaweicloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_BAIDUCLOUD}-cdn`, "common.provider.baiducloud.cdn", "/imgs/providers/baiducloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_VOLCENGINE}-cdn`, "common.provider.volcengine.cdn", "/imgs/providers/volcengine.svg"],
  [`${ACCESS_PROVIDER_TYPE_VOLCENGINE}-live`, "common.provider.volcengine.live", "/imgs/providers/volcengine.svg"],
  [`${ACCESS_PROVIDER_TYPE_QINIU}-cdn`, "common.provider.qiniu.cdn", "/imgs/providers/qiniu.svg"],
  [`${ACCESS_PROVIDER_TYPE_DOGECLOUD}-cdn`, "common.provider.dogecloud.cdn", "/imgs/providers/dogecloud.svg"],
  [`${ACCESS_PROVIDER_TYPE_BYTEPLUS}-cdn`, "common.provider.byteplus.cdn", "/imgs/providers/byteplus.svg"],
  [`${ACCESS_PROVIDER_TYPE_KUBERNETES}-secret`, "common.provider.kubernetes.secret", "/imgs/providers/kubernetes.svg"],
];

export const deployTargetsMap: Map<DeployTarget["type"], DeployTarget> = new Map(
  deployTargetList.map(([type, name, icon]) => [type, { type, provider: type.split("-")[0], name, icon }])
);

export const deployTargets = deployTargetList.map(([type, name, icon]) => ({ type, provider: type.split("-")[0], name, icon }));
