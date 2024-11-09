import { WorkflowNode } from "@/domain/workflow";
import { memo } from "react";
import DeployToAliyunOSS from "./DeployToAliyunOss";

export type DeployFormProps = {
  data: WorkflowNode;
  defaultProivder?: string;
};
const DeployForm = ({ data, defaultProivder }: DeployFormProps) => {
  return getForm(data, defaultProivder);
};

export default memo(DeployForm);

const getForm = (data: WorkflowNode, defaultProivder?: string) => {
  const provider = defaultProivder || data.config?.providerType;
  switch (provider) {
    case "aliyun-oss":
      return <DeployToAliyunOSS data={data} />;
    case "tencent":
      return <TencentForm data={data} />;
    case "aws":
      return <AwsForm data={data} />;
    default:
      return <></>;
  }
};

