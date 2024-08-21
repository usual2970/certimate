import { Separator } from "../ui/separator";

type DeployProgressProps = {
  phase?: "check" | "apply" | "deploy";
  phaseSuccess?: boolean;
};

const DeployProgress = ({ phase, phaseSuccess }: DeployProgressProps) => {
  let rs = <> </>;
  if (phase === "check") {
    if (phaseSuccess) {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">检查 </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap  text-muted-foreground">获取</div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">部署</div>
        </div>
      );
    } else {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-red-600">检查 </div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap  text-muted-foreground">获取</div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">部署</div>
        </div>
      );
    }
  }

  if (phase === "apply") {
    if (phaseSuccess) {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">检查 </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-green-600">获取</div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">部署</div>
        </div>
      );
    } else {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">检查 </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-red-600">获取</div>
          <Separator className="h-1 grow" />
          <div className="text-xs text-nowrap text-muted-foreground">部署</div>
        </div>
      );
    }
  }

  if (phase === "deploy") {
    if (phaseSuccess) {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">检查 </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-green-600">获取</div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap text-green-600">部署</div>
        </div>
      );
    } else {
      rs = (
        <div className="flex items-center">
          <div className="text-xs text-nowrap text-green-600">检查 </div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap  text-green-600">获取</div>
          <Separator className="h-1 grow bg-green-600" />
          <div className="text-xs text-nowrap text-red-600">部署</div>
        </div>
      );
    }
  }

  return rs;
};

export default DeployProgress;
