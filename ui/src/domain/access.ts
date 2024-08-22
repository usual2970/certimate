import { z } from "zod";

export const accessTypeMap: Map<string, [string, string]> = new Map([
  ["tencent", ["腾讯云", "/imgs/providers/tencent.svg"]],
  ["aliyun", ["阿里云", "/imgs/providers/aliyun.svg"]],
]);

export const accessFormType = z.union(
  [z.literal("aliyun"), z.literal("tencent")],
  { message: "请选择云服务商" }
);

export type Access = {
  id: string;
  name: string;
  configType: string;
  config: TencentConfig | AliyunConfig;
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
