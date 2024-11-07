import { WorkflowNode } from "@/domain/workflow";
import { memo } from "react";
import DeployToAliyunOSS from "./DeployToAliyunOss";

export type DeployFormProps = {
  data: WorkflowNode;
};
const DeployForm = ({ data }: DeployFormProps) => {
  return getForm(data);
};

export default memo(DeployForm);

const getForm = (data: WorkflowNode) => {
  switch (data.config?.providerType) {
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

