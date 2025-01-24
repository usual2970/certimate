import { type AccessUsageType } from "./provider";

export interface AccessModel extends BaseModel {
  name: string;
  provider: string;
  config: /*
    注意：如果追加新的类型，请保持以 ASCII 排序。
    NOTICE: If you add new type, please keep ASCII order.
  */ Record<string, unknown> &
    (
      | AccessConfigForACMEHttpReq
      | AccessConfigForAliyun
      | AccessConfigForAWS
      | AccessConfigForAzure
      | AccessConfigForBaiduCloud
      | AccessConfigForBytePlus
      | AccessConfigForCloudflare
      | AccessConfigForClouDNS
      | AccessConfigForDogeCloud
      | AccessConfigForEdgio
      | AccessConfigForGname
      | AccessConfigForGoDaddy
      | AccessConfigForHuaweiCloud
      | AccessConfigForKubernetes
      | AccessConfigForLocal
      | AccessConfigForNameDotCom
      | AccessConfigForNameSilo
      | AccessConfigForPowerDNS
      | AccessConfigForQiniu
      | AccessConfigForRainYun
      | AccessConfigForSSH
      | AccessConfigForTencentCloud
      | AccessConfigForUCloud
      | AccessConfigForVolcEngine
      | AccessConfigForWebhook
      | AccessConfigForWestcn
    );
  usage: AccessUsageType;
}

// #region AccessConfig
export type AccessConfigForACMEHttpReq = {
  endpoint: string;
  mode?: string;
  username?: string;
  password?: string;
};

export type AccessConfigForAliyun = {
  accessKeyId: string;
  accessKeySecret: string;
};

export type AccessConfigForAWS = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type AccessConfigForAzure = {
  tenantId: string;
  clientId: string;
  clientSecret: string;
  environment?: string;
};

export type AccessConfigForBaiduCloud = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type AccessConfigForBytePlus = {
  accessKey: string;
  secretKey: string;
};

export type AccessConfigForCloudflare = {
  dnsApiToken: string;
};

export type AccessConfigForClouDNS = {
  authId: string;
  authPassword: string;
};

export type AccessConfigForDogeCloud = {
  accessKey: string;
  secretKey: string;
};

export type AccessConfigForEdgio = {
  clientId: string;
  clientSecret: string;
};

export type AccessConfigForGname = {
  appId: string;
  appKey: string;
};

export type AccessConfigForGoDaddy = {
  apiKey: string;
  apiSecret: string;
};

export type AccessConfigForHuaweiCloud = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type AccessConfigForKubernetes = {
  kubeConfig?: string;
};

export type AccessConfigForLocal = NonNullable<unknown>;

export type AccessConfigForNameDotCom = {
  username: string;
  apiToken: string;
};

export type AccessConfigForNameSilo = {
  apiKey: string;
};

export type AccessConfigForNS1 = {
  apiKey: string;
};

export type AccessConfigForPowerDNS = {
  apiUrl: string;
  apiKey: string;
};

export type AccessConfigForQiniu = {
  accessKey: string;
  secretKey: string;
};

export type AccessConfigForRainYun = {
  apiKey: string;
};

export type AccessConfigForSSH = {
  host: string;
  port: number;
  username: string;
  password?: string;
  key?: string;
  keyPassphrase?: string;
};

export type AccessConfigForTencentCloud = {
  secretId: string;
  secretKey: string;
};

export type AccessConfigForUCloud = {
  privateKey: string;
  publicKey: string;
  projectId?: string;
};

export type AccessConfigForVolcEngine = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type AccessConfigForWebhook = {
  url: string;
};

export type AccessConfigForWestcn = {
  username: string;
  apiPassword: string;
};
// #endregion
