import { useTranslation } from "react-i18next";
import { DeleteOutlined as DeleteOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { Dropdown } from "antd";

import Show from "@/components/Show";
import { deployProvidersMap } from "@/domain/provider";
import { notifyChannelsMap } from "@/domain/settings";
import { type WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import AddNode from "./AddNode";
import PanelBody from "./PanelBody";
import { usePanel } from "./PanelProvider";

type NodeProps = {
  data: WorkflowNode;
};

const i18nPrefix = "workflow.node";

const Node = ({ data }: NodeProps) => {
  const { updateNode, removeNode } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeNode"]));
  const handleNameBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    updateNode({ ...data, name: e.target.innerText });
  };

  const { showPanel } = usePanel();

  const { t } = useTranslation();

  const handleNodeSettingClick = () => {
    showPanel({
      name: data.name,
      children: <PanelBody data={data} />,
    });
  };

  const getSetting = () => {
    if (!data.validated) {
      return <>{t(`${i18nPrefix}.setting.label`)}</>;
    }

    switch (data.type) {
      case WorkflowNodeType.Start:
        return (
          <div className="flex space-x-2 items-baseline">
            <div className="text-stone-700">
              <Show when={data.config?.executionMethod == "auto"} fallback={<>{t(`workflow.props.trigger.manual`)}</>}>
                {t(`workflow.props.trigger.auto`) + ":"}
              </Show>
            </div>
            <Show when={data.config?.executionMethod == "auto"}>
              <div className="text-muted-foreground">{data.config?.crontab as string}</div>
            </Show>
          </div>
        );
      case WorkflowNodeType.Apply:
        return <div className="text-muted-foreground truncate">{data.config?.domain as string}</div>;
      case WorkflowNodeType.Deploy: {
        const provider = deployProvidersMap.get(data.config?.providerType as string);
        return (
          <div className="flex space-x-2 items-center text-muted-foreground">
            <img src={provider?.icon} className="w-6 h-6" />
            <div>{t(provider?.name ?? "")}</div>
          </div>
        );
      }
      case WorkflowNodeType.Notify: {
        const channelLabel = notifyChannelsMap.get(data.config?.channel as string);
        return (
          <div className="flex space-x-2 items-center justify-between">
            <div className="text-stone-700 truncate">{t(channelLabel?.name ?? "")}</div>
            <div className="text-muted-foreground truncate">{(data.config?.subject as string) ?? ""}</div>
          </div>
        );
      }

      default:
        return <>{t(`${i18nPrefix}.setting.label`)}</>;
    }
  };

  return (
    <>
      <div className="rounded-md shadow-md w-[260px] relative">
        {data.type != WorkflowNodeType.Start && (
          <>
            <Dropdown
              menu={{
                items: [
                  {
                    key: "delete",
                    label: t(`${i18nPrefix}.delete.label`),
                    icon: <DeleteOutlinedIcon />,
                    danger: true,
                    onClick: () => {
                      removeNode(data.id);
                    },
                  },
                ],
              }}
              trigger={["click"]}
            >
              <div className="absolute right-2 top-1 cursor-pointer">
                <EllipsisOutlinedIcon className="text-white" size={17} />
              </div>
            </Dropdown>
          </>
        )}

        <div className="w-[260px] h-[60px] flex flex-col justify-center items-center bg-primary text-white rounded-t-md px-5">
          <div
            contentEditable
            suppressContentEditableWarning
            onBlur={handleNameBlur}
            className="w-full text-center outline-none focus:bg-white focus:text-stone-600 focus:rounded-sm"
          >
            {data.name}
          </div>
        </div>
        <div className="p-2 text-sm text-primary flex flex-col justify-center bg-white">
          <div className="leading-7 text-primary cursor-pointer" onClick={handleNodeSettingClick}>
            {getSetting()}
          </div>
        </div>
      </div>
      <AddNode data={data} />
    </>
  );
};

export default Node;
