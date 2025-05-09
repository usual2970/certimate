export interface AccessModel extends BaseModel {
  name: string;
  provider: string;
  config: /*
    注意：如果追加新的类型，请保持以 ASCII 排序。
    NOTICE: If you add new type, please keep ASCII order.
  */ Record<string, unknown> &
    (
      | AccessConfigFor1Panel
      | AccessConfigForACMEHttpReq
      | AccessConfigForAliyun
      | AccessConfigForAWS
      | AccessConfigForAzure
      | AccessConfigForBaiduCloud
      | AccessConfigForBaishan
      | AccessConfigForBaotaPanel
      | AccessConfigForBunny
      | AccessConfigForBytePlus
      | AccessConfigForCacheFly
      | AccessConfigForCdnfly
      | AccessConfigForCloudflare
      | AccessConfigForClouDNS
      | AccessConfigForCMCCCloud
      | AccessConfigForDeSEC
      | AccessConfigForDingTalkBot
      | AccessConfigForDNSLA
      | AccessConfigForDogeCloud
      | AccessConfigForDynv6
      | AccessConfigForEdgio
      | AccessConfigForEmail
      | AccessConfigForGcore
      | AccessConfigForGname
      | AccessConfigForGoDaddy
      | AccessConfigForGoEdge
      | AccessConfigForGoogleTrustServices
      | AccessConfigForHuaweiCloud
      | AccessConfigForJDCloud
      | AccessConfigForKubernetes
      | AccessConfigForLarkBot
      | AccessConfigForMattermost
      | AccessConfigForNamecheap
      | AccessConfigForNameDotCom
      | AccessConfigForNameSilo
      | AccessConfigForPorkbun
      | AccessConfigForPowerDNS
      | AccessConfigForProxmoxVE
      | AccessConfigForQiniu
      | AccessConfigForRainYun
      | AccessConfigForSafeLine
      | AccessConfigForSSH
      | AccessConfigForSSLCom
      | AccessConfigForTelegram
      | AccessConfigForTencentCloud
      | AccessConfigForUCloud
      | AccessConfigForUpyun
      | AccessConfigForVercel
      | AccessConfigForVolcEngine
      | AccessConfigForWangsu
      | AccessConfigForWebhook
      | AccessConfigForWeComBot
      | AccessConfigForWestcn
      | AccessConfigForZeroSSL
    );
  reserve?: "ca" | "notification";
}

// #region AccessConfig
export type AccessConfigFor1Panel = {
  apiUrl: string;
  apiKey: string;
  allowInsecureConnections?: boolean;
};

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

export type AccessConfigForBaishan = {
  apiToken: string;
};

export type AccessConfigForBaotaPanel = {
  apiUrl: string;
  apiKey: string;
  allowInsecureConnections?: boolean;
};

export type AccessConfigForBunny = {
  apiKey: string;
};

export type AccessConfigForBytePlus = {
  accessKey: string;
  secretKey: string;
};

export type AccessConfigForCacheFly = {
  apiToken: string;
};

export type AccessConfigForCdnfly = {
  apiUrl: string;
  apiKey: string;
  apiSecret: string;
  allowInsecureConnections?: boolean;
};

export type AccessConfigForCloudflare = {
  dnsApiToken: string;
  zoneApiToken?: string;
};

export type AccessConfigForClouDNS = {
  authId: string;
  authPassword: string;
};

export type AccessConfigForCMCCCloud = {
  accessKeyId: string;
  accessKeySecret: string;
};

export type AccessConfigForDeSEC = {
  token: string;
};

export type AccessConfigForDingTalkBot = {
  webhookUrl: string;
  secret?: string;
};

export type AccessConfigForDNSLA = {
  apiId: string;
  apiSecret: string;
};

export type AccessConfigForDogeCloud = {
  accessKey: string;
  secretKey: string;
};

export type AccessConfigForDynv6 = {
  httpToken: string;
};

export type AccessConfigForEdgio = {
  clientId: string;
  clientSecret: string;
};

export type AccessConfigForEmail = {
  smtpHost: string;
  smtpPort: number;
  smtpTls: boolean;
  username: string;
  password: string;
  defaultSenderAddress?: string;
  defaultReceiverAddress?: string;
};

export type AccessConfigForGcore = {
  apiToken: string;
};

export type AccessConfigForGname = {
  appId: string;
  appKey: string;
};

export type AccessConfigForGoDaddy = {
  apiKey: string;
  apiSecret: string;
};

export type AccessConfigForGoEdge = {
  apiUrl: string;
  accessKeyId: string;
  accessKey: string;
  allowInsecureConnections?: boolean;
};

export type AccessConfigForGoogleTrustServices = {
  eabKid: string;
  eabHmacKey: string;
};

export type AccessConfigForHuaweiCloud = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type AccessConfigForJDCloud = {
  accessKeyId: string;
  accessKeySecret: string;
};

export type AccessConfigForKubernetes = {
  kubeConfig?: string;
};

export type AccessConfigForLarkBot = {
  webhookUrl: string;
};

export type AccessConfigForMattermost = {
  serverUrl: string;
  username: string;
  password: string;
  defaultChannelId?: string;
};

export type AccessConfigForNamecheap = {
  username: string;
  apiKey: string;
};

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

export type AccessConfigForPorkbun = {
  apiKey: string;
  secretApiKey: string;
};

export type AccessConfigForPowerDNS = {
  apiUrl: string;
  apiKey: string;
  allowInsecureConnections?: boolean;
};

export type AccessConfigForProxmoxVE = {
  apiUrl: string;
  apiToken: string;
  apiTokenSecret?: string;
  allowInsecureConnections?: boolean;
};

export type AccessConfigForQiniu = {
  accessKey: string;
  secretKey: string;
};

export type AccessConfigForRainYun = {
  apiKey: string;
};

export type AccessConfigForSafeLine = {
  apiUrl: string;
  apiToken: string;
  allowInsecureConnections?: boolean;
};

export type AccessConfigForSSH = {
  host: string;
  port: number;
  username: string;
  password?: string;
  key?: string;
  keyPassphrase?: string;
};

export type AccessConfigForSSLCom = {
  eabKid: string;
  eabHmacKey: string;
};

export type AccessConfigForTelegram = {
  botToken: string;
  defaultChatId?: number;
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

export type AccessConfigForUpyun = {
  username: string;
  password: string;
};

export type AccessConfigForVercel = {
  apiAccessToken: string;
  teamId?: string;
};

export type AccessConfigForVolcEngine = {
  accessKeyId: string;
  secretAccessKey: string;
};

export type AccessConfigForWangsu = {
  accessKeyId: string;
  accessKeySecret: string;
  apiKey: string;
};

export type AccessConfigForWebhook = {
  url: string;
  method: string;
  headers?: string;
  allowInsecureConnections?: boolean;
  defaultDataForDeployment?: string;
  defaultDataForNotification?: string;
};

export type AccessConfigForWeComBot = {
  webhookUrl: string;
};

export type AccessConfigForWestcn = {
  username: string;
  apiPassword: string;
};

export type AccessConfigForZeroSSL = {
  eabKid: string;
  eabHmacKey: string;
};
// #endregion
