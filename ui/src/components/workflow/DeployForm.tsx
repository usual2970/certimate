import { memo } from "react";

import { type WorkflowNode } from "@/domain/workflow";
import DeployToAliyunALB from "./DeployToAliyunALB";
import DeployToAliyunCDN from "./DeployToAliyunCDN";
import DeployToAliyunCLB from "./DeployToAliyunCLB";
import DeployToAliyunNLB from "./DeployToAliyunNLB";
import DeployToAliyunOSS from "./DeployToAliyunOss";
import DeployToBaiduCloudCDN from "./DeployToBaiduCloudCDN";
import DeployToBytePlusCDN from "./DeployToByteplusCDN";
import DeployToDogeCloudCDN from "./DeployToDogeCloudCDN";
import DeployToHuaweiCloudCDN from "./DeployToHuaweiCloudCDN";
import DeployToHuaweiCloudELB from "./DeployToHuaweiCloudELB";
import DeployToKubernetesSecret from "./DeployToKubernetesSecret";
import DeployToLocal from "./DeployToLocal";
import DeployToQiniuCDN from "./DeployToQiniuCDN";
import DeployToSSH from "./DeployToSSH";
import DeployToTencentCDN from "./DeployToTencentCDN";
import DeployToTencentCLB from "./DeployToTencentCLB";
import DeployToTencentCOS from "./DeployToTencentCOS";
import DeployToTencentEO from "./DeployToTencentTEO";
import DeployToVolcEngineCDN from "./DeployToVolcengineCDN";
import DeployToVolcEngineLive from "./DeployToVolcengineLive";
import DeployToWebhook from "./DeployToWebhook";

export type DeployFormProps = {
  data: WorkflowNode;
  defaultProivder?: string;
};

const DeployForm = ({ data, defaultProivder }: DeployFormProps) => {
  return <div className="dark:text-stone-200">{getForm(data, defaultProivder)}</div>;
};

export default memo(DeployForm);

const getForm = (data: WorkflowNode, defaultProivder?: string) => {
  const provider = defaultProivder || data.config?.providerType;
  switch (provider) {
    case "aliyun-oss":
      return <DeployToAliyunOSS data={data} />;
    case "aliyun-alb":
      return <DeployToAliyunALB data={data} />;
    case "aliyun-cdn":
    case "aliyun-dcdn":
      return <DeployToAliyunCDN data={data} />;
    case "aliyun-clb":
      return <DeployToAliyunCLB data={data} />;
    case "aliyun-nlb":
      return <DeployToAliyunNLB data={data} />;
    case "baiducloud-cdn":
      return <DeployToBaiduCloudCDN data={data} />;
    case "dogecloud-cdn":
      return <DeployToDogeCloudCDN data={data} />;
    case "huaweicloud-cdn":
      return <DeployToHuaweiCloudCDN data={data} />;
    case "huaweicloud-elb":
      return <DeployToHuaweiCloudELB data={data} />;
    case "k8s-secret":
      return <DeployToKubernetesSecret data={data} />;
    case "qiniu-cdn":
      return <DeployToQiniuCDN data={data} />;
    case "webhook":
      return <DeployToWebhook data={data} />;
    case "tencentcloud-cdn":
    case "tencentcloud-ecdn":
      return <DeployToTencentCDN data={data} />;
    case "tencentcloud-clb":
      return <DeployToTencentCLB data={data} />;
    case "tencentcloud-cos":
      return <DeployToTencentCOS data={data} />;
    case "tencentcloud-eo":
      return <DeployToTencentEO data={data} />;
    case "ssh":
      return <DeployToSSH data={data} />;
    case "local":
      return <DeployToLocal data={data} />;
    case "byteplus-cdn":
      return <DeployToBytePlusCDN data={data} />;
    case "volcengine-cdn":
      return <DeployToVolcEngineCDN data={data} />;
    case "volcengine-live":
      return <DeployToVolcEngineLive data={data} />;
    default:
      return <></>;
  }
};
