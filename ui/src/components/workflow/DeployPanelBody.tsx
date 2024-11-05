import { accessProviders } from "@/domain/access";
import { WorkflowNode } from "@/domain/workflow";
import { memo } from "react";
import { useTranslation } from "react-i18next";

type DeployPanelBodyProps = {
  data: WorkflowNode;
};
const DeployPanelBody = ({ data }: DeployPanelBodyProps) => {
  const { t } = useTranslation();
  return (
    <>
      {/* 默认展示服务商列表 */}
      <div className="text-lg font-semibold text-gray-700">选择服务商</div>
      {accessProviders
        .filter((provider) => provider[3] === "apply" || provider[3] === "all")
        .reduce((acc: string[][][], provider, index) => {
          if (index % 2 === 0) {
            acc.push([provider]);
          } else {
            acc[acc.length - 1].push(provider);
          }
          return acc;
        }, [])
        .map((providerRow, rowIndex) => (
          <div key={rowIndex} className="flex space-x-5">
            {providerRow.map((provider, index) => (
              <div key={index} className="flex space-x-2 w-1/3 items-center cursor-pointer hover:bg-slate-100 p-2 rounded-sm">
                <img src={provider[2]} alt={provider[1]} className="w-8 h-8" />
                <div className="text-muted-foreground">{t(provider[1])}</div>
              </div>
            ))}
          </div>
        ))}
    </>
  );
};

export default memo(DeployPanelBody);
