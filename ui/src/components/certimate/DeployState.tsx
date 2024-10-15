import { CircleCheck, CircleX } from "lucide-react";

import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Deployment } from "@/domain/deployment";

type DeployStateProps = {
  deployment: Deployment;
};

const DeployState = ({ deployment }: DeployStateProps) => {
  // 获取指定阶段的错误信息
  const error = (state: "check" | "apply" | "deploy") => {
    if (!deployment.log[state]) {
      return "";
    }
    return deployment.log[state][deployment.log[state].length - 1].error;
  };

  return (
    <>
      {(deployment.phase === "deploy" && deployment.phaseSuccess) || deployment.wholeSuccess ? (
        <CircleCheck size={16} className="text-green-700" />
      ) : (
        <>
          {error(deployment.phase).length ? (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild className="cursor-pointer">
                  <CircleX size={16} className="text-red-700" />
                </TooltipTrigger>
                <TooltipContent className="max-w-[35em]">{error(deployment.phase)}</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          ) : (
            <CircleX size={16} className="text-red-700" />
          )}
        </>
      )}
    </>
  );
};

export default DeployState;
