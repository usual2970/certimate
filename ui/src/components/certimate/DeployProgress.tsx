import { useTranslation } from "react-i18next";

import { cn } from "@/lib/utils"


import { Separator } from "../ui/separator";

type DeployProgressProps = {
  phase?: "check" | "apply" | "deploy";
  phaseSuccess?: boolean;
};

const DeployProgress = ({ phase, phaseSuccess }: DeployProgressProps) => {
  const { t } = useTranslation();

  let step = 0;

  if (phase === "check") {
    step = 1
  } else if (phase === "apply") {
    step = 2
  } else if (phase === "deploy") {
    step = 3
  }

  return (
    <div className="flex items-center">
      <div className={
        cn(
          "text-xs text-nowrap",
          step === 1 ? phaseSuccess ? "text-green-600" : "text-red-600" : "",
          step > 1 ? "text-green-600" : "",
        )
      }>
        {t('deploy.progress.check')}
      </div>
      <Separator className={
        cn(
          "h-1 grow max-w-[60px]",
          step > 1 ? "bg-green-600" : "",
        )
      } />
      <div className={
        cn(
          "text-xs text-nowrap",
          step < 2 ? "text-muted-foreground" : "",
          step === 2 ? phaseSuccess ? "text-green-600" : "text-red-600" : "",
          step > 2 ? "text-green-600" : "",
        )
      }>
        {t('deploy.progress.apply')}
      </div>
      <Separator className={
        cn(
          "h-1 grow max-w-[60px]",
          step > 2 ? "bg-green-600" : "",
        )
      } />
      <div className={
        cn(
          "text-xs text-nowrap",
          step < 3 ? "text-muted-foreground" : "",
          step === 3 ? phaseSuccess ? "text-green-600" : "text-red-600" : "",
          step > 3 ? "text-green-600" : "",
        )
      }>
        {t('deploy.progress.deploy')}
      </div>
    </div>
  )
};

export default DeployProgress;
