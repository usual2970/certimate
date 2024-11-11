import { WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import AddNode from "./AddNode";
import { useWorkflowStore, WorkflowState } from "@/providers/workflow";
import { useShallow } from "zustand/shallow";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";
import { Ellipsis, Trash2 } from "lucide-react";
import { usePanel } from "./PanelProvider";
import PanelBody from "./PanelBody";
import { useTranslation } from "react-i18next";
import Show from "../Show";
import { deployTargetsMap } from "@/domain/domain";
import { channelLabelMap } from "@/domain/settings";

type NodeProps = {
  data: WorkflowNode;
};

const i18nPrefix = "workflow.node";

const selectState = (state: WorkflowState) => ({
  updateNode: state.updateNode,
  removeNode: state.removeNode,
});
const Node = ({ data }: NodeProps) => {
  const { updateNode, removeNode } = useWorkflowStore(useShallow(selectState));
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
    console.log(data);
    if (!data.validated) {
      return <>{t(`${i18nPrefix}.setting.label`)}</>;
    }

    switch (data.type) {
      case WorkflowNodeType.Start:
        return (
          <div className="flex space-x-2 items-baseline">
            <div className="text-stone-700">
              <Show when={data.config?.executionMethod == "auto"} fallback={<>{t(`workflow.node.start.form.executionMethod.options.manual`)}</>}>
                {t(`workflow.node.start.form.executionMethod.options.auto`) + ":"}
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
        const provider = deployTargetsMap.get(data.config?.providerType as string);
        return (
          <div className="flex space-x-2 items-center text-muted-foreground">
            <img src={provider?.icon} className="w-6 h-6" />
            <div>{t(provider?.name ?? "")}</div>
          </div>
        );
      }
      case WorkflowNodeType.Notify: {
        const channelLabel = channelLabelMap.get(data.config?.channel as string);
        return (
          <div className="flex space-x-2 items-baseline">
            <div className="text-stone-700">{t(channelLabel?.label ?? "")}</div>
            <div className="text-muted-foreground truncate">{(data.config?.title as string) ?? ""}</div>
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
            <DropdownMenu>
              <DropdownMenuTrigger className="absolute right-2 top-1">
                <Ellipsis className="text-white" size={17} />
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem
                  className="flex space-x-2 text-red-600"
                  onClick={() => {
                    removeNode(data.id);
                  }}
                >
                  <Trash2 size={16} /> <div>{t(`${i18nPrefix}.delete.label`)}</div>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
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
