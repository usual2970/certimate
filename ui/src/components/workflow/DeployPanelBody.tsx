import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import Show from "@/components/Show";
import { deployProvidersMap } from "@/domain/provider";
import { type WorkflowNode } from "@/domain/workflow";

import DeployNodeForm from "./node/DeployNodeForm";

type DeployPanelBodyProps = {
  data: WorkflowNode;
};

const DeployPanelBody = ({ data }: DeployPanelBodyProps) => {
  const { t } = useTranslation();

  const [providerType, setProviderType] = useState("");

  useEffect(() => {
    if (data.config?.providerType) {
      setProviderType(data.config.providerType as string);
    }
  }, [data]);
  return (
    <>
      {/* 默认展示服务商列表 */}
      <Show when={!providerType} fallback={<DeployNodeForm data={data} defaultProivderType={providerType} />}>
        <div className="text-lg font-semibold text-gray-700">选择服务商</div>
        {Array.from(deployProvidersMap.values()).map((provider, index) => {
          return (
            <div
              key={index}
              className="flex space-x-2 w-[47%] items-center cursor-pointer hover:bg-slate-100 p-2 rounded-sm"
              onClick={() => {
                setProviderType(provider.type);
              }}
            >
              <img src={provider.icon} alt={provider.type} className="w-8 h-8" />
              <div className="text-muted-foreground text-sm">{t(provider.name)}</div>
            </div>
          );
        })}
      </Show>
    </>
  );
};

export default memo(DeployPanelBody);
