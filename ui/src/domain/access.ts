import { z } from "zod";

export const accessTypeMap: Map<string, [string, string]> = new Map([
  ["tencent", ["腾讯云", "/imgs/providers/tencent.svg"]],
  ["aliyun", ["阿里云", "/imgs/providers/aliyun.svg"]],
  ["ssh", ["SSH部署", "/imgs/providers/ssh.png"]],
]);

export const accessFormType = z.union(
  [z.literal("aliyun"), z.literal("tencent"), z.literal("ssh")],
  { message: "请选择云服务商" }
);

export type Access = {
  id: string;
  name: string;
  configType: string;
  config: TencentConfig | AliyunConfig | SSHConfig;
  deleted?: string;
  created?: string;
  updated?: string;
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
