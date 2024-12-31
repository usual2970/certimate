import {
  ACCESS_PROVIDER_ACMEHTTPREQ,
  ACCESS_PROVIDER_ALIYUN,
  ACCESS_PROVIDER_AWS,
  ACCESS_PROVIDER_BAIDUCLOUD,
  ACCESS_PROVIDER_BYTEPLUS,
  ACCESS_PROVIDER_CLOUDFLARE,
  ACCESS_PROVIDER_DOGECLOUD,
  ACCESS_PROVIDER_HUAWEICLOUD,
  ACCESS_PROVIDER_KUBERNETES,
  ACCESS_PROVIDER_LOCAL,
  ACCESS_PROVIDER_NAMEDOTCOM,
  ACCESS_PROVIDER_NAMESILO,
  ACCESS_PROVIDER_GODADDY,
  ACCESS_PROVIDER_POWERDNS,
  ACCESS_PROVIDER_QINIU,
  ACCESS_PROVIDER_SSH,
  ACCESS_PROVIDER_TENCENTCLOUD,
  ACCESS_PROVIDER_VOLCENGINE,
  ACCESS_PROVIDER_WEBHOOK,
  type AccessUsageType,
} from "./access";

export type AccessProvider = {
  type: string;
  name: string;
  icon: string;
  usage: AccessUsageType;
};

export const accessProvidersMap: Map<AccessProvider["type"], AccessProvider> = new Map(
  /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
  [
    [ACCESS_PROVIDER_LOCAL, "common.provider.local", "/imgs/providers/local.svg", "deploy"],
    [ACCESS_PROVIDER_SSH, "common.provider.ssh", "/imgs/providers/ssh.svg", "deploy"],
    [ACCESS_PROVIDER_WEBHOOK, "common.provider.webhook", "/imgs/providers/webhook.svg", "deploy"],
    [ACCESS_PROVIDER_KUBERNETES, "common.provider.kubernetes", "/imgs/providers/kubernetes.svg", "deploy"],
    [ACCESS_PROVIDER_ALIYUN, "common.provider.aliyun", "/imgs/providers/aliyun.svg", "all"],
    [ACCESS_PROVIDER_TENCENTCLOUD, "common.provider.tencentcloud", "/imgs/providers/tencentcloud.svg", "all"],
    [ACCESS_PROVIDER_HUAWEICLOUD, "common.provider.huaweicloud", "/imgs/providers/huaweicloud.svg", "all"],
    [ACCESS_PROVIDER_BAIDUCLOUD, "common.provider.baiducloud", "/imgs/providers/baiducloud.svg", "all"],
    [ACCESS_PROVIDER_QINIU, "common.provider.qiniu", "/imgs/providers/qiniu.svg", "deploy"],
    [ACCESS_PROVIDER_DOGECLOUD, "common.provider.dogecloud", "/imgs/providers/dogecloud.svg", "deploy"],
    [ACCESS_PROVIDER_VOLCENGINE, "common.provider.volcengine", "/imgs/providers/volcengine.svg", "all"],
    [ACCESS_PROVIDER_BYTEPLUS, "common.provider.byteplus", "/imgs/providers/byteplus.svg", "all"],
    [ACCESS_PROVIDER_AWS, "common.provider.aws", "/imgs/providers/aws.svg", "apply"],
    [ACCESS_PROVIDER_CLOUDFLARE, "common.provider.cloudflare", "/imgs/providers/cloudflare.svg", "apply"],
    [ACCESS_PROVIDER_NAMEDOTCOM, "common.provider.namedotcom", "/imgs/providers/namedotcom.svg", "apply"],
    [ACCESS_PROVIDER_NAMESILO, "common.provider.namesilo", "/imgs/providers/namesilo.svg", "apply"],
    [ACCESS_PROVIDER_GODADDY, "common.provider.godaddy", "/imgs/providers/godaddy.svg", "apply"],
    [ACCESS_PROVIDER_POWERDNS, "common.provider.powerdns", "/imgs/providers/powerdns.svg", "apply"],
    [ACCESS_PROVIDER_ACMEHTTPREQ, "common.provider.acmehttpreq", "/imgs/providers/acmehttpreq.svg", "apply"],
  ].map(([type, name, icon, usage]) => [type, { type, name, icon, usage: usage as AccessUsageType }])
);

export type DeployProvider = {
  type: string;
  name: string;
  icon: string;
  provider: AccessProvider["type"];
};

export const deployProvidersMap: Map<DeployProvider["type"], DeployProvider> = new Map(
  [
    /*
   注意：此处的顺序决定显示在前端的顺序。
   NOTICE: The following order determines the order displayed at the frontend.
  */
    [`${ACCESS_PROVIDER_LOCAL}`, "common.provider.local"],
    [`${ACCESS_PROVIDER_SSH}`, "common.provider.ssh"],
    [`${ACCESS_PROVIDER_WEBHOOK}`, "common.provider.webhook"],
    [`${ACCESS_PROVIDER_KUBERNETES}-secret`, "common.provider.kubernetes.secret"],
    [`${ACCESS_PROVIDER_ALIYUN}-oss`, "common.provider.aliyun.oss"],
    [`${ACCESS_PROVIDER_ALIYUN}-cdn`, "common.provider.aliyun.cdn"],
    [`${ACCESS_PROVIDER_ALIYUN}-dcdn`, "common.provider.aliyun.dcdn"],
    [`${ACCESS_PROVIDER_ALIYUN}-clb`, "common.provider.aliyun.clb"],
    [`${ACCESS_PROVIDER_ALIYUN}-alb`, "common.provider.aliyun.alb"],
    [`${ACCESS_PROVIDER_ALIYUN}-nlb`, "common.provider.aliyun.nlb"],
    [`${ACCESS_PROVIDER_TENCENTCLOUD}-cdn`, "common.provider.tencentcloud.cdn"],
    [`${ACCESS_PROVIDER_TENCENTCLOUD}-ecdn`, "common.provider.tencentcloud.ecdn"],
    [`${ACCESS_PROVIDER_TENCENTCLOUD}-clb`, "common.provider.tencentcloud.clb"],
    [`${ACCESS_PROVIDER_TENCENTCLOUD}-cos`, "common.provider.tencentcloud.cos"],
    [`${ACCESS_PROVIDER_TENCENTCLOUD}-eo`, "common.provider.tencentcloud.eo"],
    [`${ACCESS_PROVIDER_HUAWEICLOUD}-cdn`, "common.provider.huaweicloud.cdn"],
    [`${ACCESS_PROVIDER_HUAWEICLOUD}-elb`, "common.provider.huaweicloud.elb"],
    [`${ACCESS_PROVIDER_BAIDUCLOUD}-cdn`, "common.provider.baiducloud.cdn"],
    [`${ACCESS_PROVIDER_VOLCENGINE}-cdn`, "common.provider.volcengine.cdn"],
    [`${ACCESS_PROVIDER_VOLCENGINE}-live`, "common.provider.volcengine.live"],
    [`${ACCESS_PROVIDER_QINIU}-cdn`, "common.provider.qiniu.cdn"],
    [`${ACCESS_PROVIDER_DOGECLOUD}-cdn`, "common.provider.dogecloud.cdn"],
    [`${ACCESS_PROVIDER_BYTEPLUS}-cdn`, "common.provider.byteplus.cdn"],
  ].map(([type, name]) => [type, { type, name, icon: accessProvidersMap.get(type.split("-")[0])!.icon, provider: type.split("-")[0] }])
);
