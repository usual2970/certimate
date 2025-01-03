import { useTranslation } from "react-i18next";
import { CloseCircleOutlined as CloseCircleOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { Avatar, Button, Card, Dropdown, Popover, Space, Typography } from "antd";
import { produce } from "immer";

import Show from "@/components/Show";
import { deployProvidersMap } from "@/domain/provider";
import { notifyChannelsMap } from "@/domain/settings";
import { type WorkflowNode, WorkflowNodeType } from "@/domain/workflow";
import { useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import PanelBody from "./PanelBody";
import { usePanel } from "./PanelProvider";
import AddNode from "./node/AddNode";

export type NodeProps = {
  node: WorkflowNode;
};

const WorkflowElement = ({ node }: NodeProps) => {
  const { t } = useTranslation();

  const { updateNode, removeNode } = useWorkflowStore(useZustandShallowSelector(["updateNode", "removeNode"]));
  const { showPanel } = usePanel();

  const renderNodeContent = () => {
    if (!node.validated) {
      return <Typography.Link>{t("workflow_node.action.configure_node")}</Typography.Link>;
    }

    switch (node.type) {
      case WorkflowNodeType.Start: {
        return (
          <div className="flex items-center justify-between space-x-2">
            <Typography.Text className="truncate">
              {node.config?.executionMethod === "auto"
                ? t("workflow.props.trigger.auto")
                : node.config?.executionMethod === "manual"
                  ? t("workflow.props.trigger.manual")
                  : ""}
            </Typography.Text>
            <Typography.Text className="truncate" type="secondary">
              {node.config?.executionMethod === "auto" ? (node.config?.crontab as string) : ""}
            </Typography.Text>
          </div>
        );
      }

      case WorkflowNodeType.Apply: {
        return <Typography.Text className="truncate">{node.config?.domain as string}</Typography.Text>;
      }

      case WorkflowNodeType.Deploy: {
        const provider = deployProvidersMap.get(node.config?.providerType as string);
        return (
          <Space>
            <Avatar src={provider?.icon} size="small" />
            <Typography.Text className="truncate">{t(provider?.name ?? "")}</Typography.Text>
          </Space>
        );
      }

      case WorkflowNodeType.Notify: {
        const channel = notifyChannelsMap.get(node.config?.channel as string);
        return (
          <div className="flex items-center justify-between space-x-2">
            <Typography.Text className="truncate">{t(channel?.name ?? "")}</Typography.Text>
            <Typography.Text className="truncate" type="secondary">
              {(node.config?.subject as string) ?? ""}
            </Typography.Text>
          </div>
        );
      }

      default: {
        return <></>;
      }
    }
  };

  const handleNodeNameBlur = (e: React.FocusEvent<HTMLDivElement>) => {
    const oldName = node.name;
    const newName = e.target.innerText.trim();
    if (oldName === newName) {
      return;
    }

    updateNode(
      produce(node, (draft) => {
        draft.name = newName;
      })
    );
  };

  const handleNodeClick = () => {
    showPanel({
      name: node.name,
      children: <PanelBody data={node} />,
    });
  };

  return (
    <>
      <Popover
        arrow={false}
        content={
          <Show when={node.type !== WorkflowNodeType.Start}>
            <Dropdown
              menu={{
                items: [
                  {
                    key: "delete",
                    label: t("workflow_node.action.delete_node"),
                    icon: <CloseCircleOutlinedIcon />,
                    danger: true,
                    onClick: () => {
                      removeNode(node.id);
                    },
                  },
                ],
              }}
              trigger={["click"]}
            >
              <Button color="primary" icon={<EllipsisOutlinedIcon />} variant="text" />
            </Dropdown>
          </Show>
        }
        overlayClassName="shadow-md"
        overlayInnerStyle={{ padding: 0 }}
        placement="rightTop"
      >
        <Card className="relative w-[256px] overflow-hidden shadow-md" styles={{ body: { padding: 0 } }} hoverable>
          <div className="bg-primary flex h-[48px] flex-col items-center justify-center truncate px-4 py-2 text-white">
            <div
              className="focus:bg-background focus:text-foreground w-full overflow-hidden text-center outline-none focus:rounded-sm"
              contentEditable
              suppressContentEditableWarning
              onBlur={handleNodeNameBlur}
            >
              {node.name}
            </div>
          </div>

          <div className="flex flex-col justify-center px-4 py-2">
            <div className="cursor-pointer text-sm" onClick={handleNodeClick}>
              {renderNodeContent()}
            </div>
          </div>
        </Card>
      </Popover>

      <AddNode node={node} />
    </>
  );
};

export default WorkflowElement;
