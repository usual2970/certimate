import { useTranslation } from "react-i18next";

import { Separator } from "../ui/separator";

type DeployProgressProps = {
  phase?: "check" | "apply" | "deploy";
  phaseSuccess?: boolean;
};

const DeployProgress = ({ phase, phaseSuccess }: DeployProgressProps) => {
  const { t } = useTranslation();

  let rs = <> </>;
  if (phase === "check") {
    if (phaseSuccess) {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">
            {t('deploy.progress.check')}
          </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap  text-muted-foreground">
            {t('deploy.progress.apply')}
          </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">
            {t('deploy.progress.deploy')}
          </div>
        </div>
      );
    } else {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-red-600">
            {t('deploy.progress.check')}
          </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap  text-muted-foreground">
            {t('deploy.progress.apply')}
          </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">
            {t('deploy.progress.deploy')}
          </div>
        </div>
      );
    }
  }

  if (phase === "apply") {
    if (phaseSuccess) {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">
            {t('deploy.progress.check')}
          </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-green-600">
            {t('deploy.progress.apply')}
          </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">
            {t('deploy.progress.deploy')}
          </div>
        </div>
      );
    } else {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">
            {t('deploy.progress.check')}
          </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-red-600">
            {t('deploy.progress.apply')}
          </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">
            {t('deploy.progress.deploy')}
          </div>
        </div>
      );
    }
  }

  if (phase === "deploy") {
    if (phaseSuccess) {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">
            {t('deploy.progress.check')}
          </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-green-600">
            {t('deploy.progress.apply')}
          </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap text-green-600">
            {t('deploy.progress.deploy')}
          </div>
        </div>
      );
    } else {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">
            {t('deploy.progress.check')}
          </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-green-600">
            {t('deploy.progress.apply')}
          </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap text-red-600">
            {t('deploy.progress.deploy')}
          </div>
        </div>
      );
    }
  }

  return rs;
};

export default DeployProgress;
