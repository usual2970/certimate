import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useRequest } from "ahooks";
import {
  Button,
  Divider,
  Empty,
  Menu,
  Modal,
  notification,
  Radio,
  Space,
  Switch,
  Table,
  theme,
  Tooltip,
  Typography,
  type MenuProps,
  type TableProps,
} from "antd";
import { PageHeader } from "@ant-design/pro-components";
import { Filter as FilterIcon, Pencil as PencilIcon, Plus as PlusIcon, Trash2 as Trash2Icon } from "lucide-react";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import { WorkflowModel } from "@/domain/workflow";
import { list as listWorkflow, remove as removeWorkflow, save as saveWorkflow } from "@/repository/workflow";
import { getErrMsg } from "@/utils/error";

const WorkflowList = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [modalApi, ModelContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const tableColumns: TableProps<WorkflowModel>["columns"] = [
    {
      key: "$index",
      align: "center",
      fixed: "left",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "name",
      title: t("workflow.props.name"),
      ellipsis: true,
      render: (_, record) => (
        <Space className="max-w-full" direction="vertical" size={4}>
          <Typography.Text ellipsis>{record.name}</Typography.Text>
          <Typography.Text type="secondary" ellipsis>
            {record.description}
          </Typography.Text>
        </Space>
      ),
    },
    {
      key: "trigger",
      title: t("workflow.props.trigger"),
      ellipsis: true,
      render: (_, record) => {
        const trigger = record.type;
        if (!trigger) {
          return "-";
        } else if (trigger === "manual") {
          return <Typography.Text>{t("workflow.props.trigger.manual")}</Typography.Text>;
        } else if (trigger === "auto") {
          return (
            <Space className="max-w-full" direction="vertical" size={4}>
              <Typography.Text>{t("workflow.props.trigger.auto")}</Typography.Text>
              <Typography.Text type="secondary">{record.crontab ?? ""}</Typography.Text>
            </Space>
          );
        }
      },
    },
    {
      key: "state",
      title: t("workflow.props.state"),
      defaultFilteredValue: searchParams.has("state") ? [searchParams.get("state") as string] : undefined,
      filterDropdown: ({ setSelectedKeys, confirm, clearFilters }) => {
        const items: Required<MenuProps>["items"] = [
          ["enabled", "workflow.props.state.filter.enabled"],
          ["disabled", "workflow.props.state.filter.disabled"],
        ].map(([key, label]) => {
          return {
            key,
            label: <Radio checked={filters["state"] === key}>{t(label)}</Radio>,
            onClick: () => {
              if (filters["state"] !== key) {
                setPage(1);
                setFilters((prev) => ({ ...prev, state: key }));
                setSelectedKeys([key]);
              }

              confirm({ closeDropdown: true });
            },
          };
        });

        const handleResetClick = () => {
          setPage(1);
          setFilters((prev) => ({ ...prev, state: undefined }));
          setSelectedKeys([]);
          clearFilters?.();
          confirm();
        };

        const handleConfirmClick = () => {
          confirm();
        };

        return (
          <div style={{ padding: 0 }}>
            <Menu items={items} selectable={false} />
            <Divider style={{ margin: 0 }} />
            <Space className="justify-end w-full" style={{ padding: themeToken.paddingSM }}>
              <Button size="small" disabled={!filters.state} onClick={handleResetClick}>
                {t("common.button.reset")}
              </Button>
              <Button type="primary" size="small" onClick={handleConfirmClick}>
                {t("common.button.ok")}
              </Button>
            </Space>
          </div>
        );
      },
      filterIcon: () => <FilterIcon size={14} />,
      render: (_, record) => {
        const enabled = record.enabled;
        return (
          <Switch
            checked={enabled}
            onChange={() => {
              handleEnabledChange(record);
            }}
          />
        );
      },
    },
    {
      key: "lastExecutedAt",
      title: "最近执行状态",
      render: () => {
        // TODO: 最近执行状态
        return <>TODO</>;
      },
    },
    {
      key: "createdAt",
      title: t("workflow.props.created_at"),
      ellipsis: true,
      render: (_, record) => {
        return dayjs(record.created!).format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      key: "updatedAt",
      title: t("workflow.props.updated_at"),
      ellipsis: true,
      render: (_, record) => {
        return dayjs(record.updated!).format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      key: "$action",
      align: "end",
      fixed: "right",
      width: 120,
      render: (_, record) => (
        <Space size={0}>
          <Tooltip title={t("workflow.action.edit")}>
            <Button
              type="link"
              icon={<PencilIcon size={16} />}
              onClick={() => {
                navigate(`/workflows/${record.id}`);
              }}
            />
          </Tooltip>
          <Tooltip title={t("workflow.action.delete")}>
            <Button
              type="link"
              danger={true}
              icon={<Trash2Icon size={16} />}
              onClick={() => {
                handleDeleteClick(record);
              }}
            />
          </Tooltip>
        </Space>
      ),
    },
  ];
  const [tableData, setTableData] = useState<WorkflowModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [filters, setFilters] = useState<Record<string, unknown>>(() => {
    return {
      state: searchParams.get("state"),
    };
  });

  const [page, setPage] = useState<number>(() => parseInt(+searchParams.get("page")! + "") || 1);
  const [pageSize, setPageSize] = useState<number>(() => parseInt(+searchParams.get("perPage")! + "") || 10);

  const { loading } = useRequest(
    () => {
      return listWorkflow({
        page: page,
        perPage: pageSize,
        enabled: (filters["state"] as string) === "enabled" ? true : (filters["state"] as string) === "disabled" ? false : undefined,
      });
    },
    {
      refreshDeps: [filters, page, pageSize],
      onSuccess: (data) => {
        setTableData(data.items);
        setTableTotal(data.totalItems);
      },
      onError: (err) => {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);
        notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
      },
    }
  );

  const handleEnabledChange = async (workflow: WorkflowModel) => {
    try {
      const resp = await saveWorkflow({
        id: workflow.id,
        enabled: !tableData.find((item) => item.id === workflow.id)?.enabled,
      });
      if (resp) {
        setTableData((prev) => {
          return prev.map((item) => {
            if (item.id === workflow.id) {
              return resp;
            }
            return item;
          });
        });
      }
    } catch (err) {
      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    }
  };

  const handleDeleteClick = (workflow: WorkflowModel) => {
    modalApi.confirm({
      title: t("workflow.action.delete"),
      content: t("workflow.action.delete.confirm"),
      onOk: async () => {
        try {
          const resp: boolean = await removeWorkflow(workflow);
          if (resp) {
            setTableData((prev) => prev.filter((item) => item.id !== workflow.id));
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  const handleCreateClick = () => {
    navigate("/workflows/");
  };

  return (
    <div className="p-4">
      {ModelContextHolder}
      {NotificationContextHolder}

      <PageHeader
        title={t("workflow.page.title")}
        extra={[
          <Button
            key="create"
            type="primary"
            icon={<PlusIcon size={16} />}
            onClick={() => {
              handleCreateClick();
            }}
          >
            {t("workflow.action.create")}
          </Button>,
        ]}
      />

      <Table<WorkflowModel>
        columns={tableColumns}
        dataSource={tableData}
        loading={loading}
        locale={{
          emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t("workflow.nodata")} />,
        }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: tableTotal,
          showSizeChanger: true,
          onChange: (page: number, pageSize: number) => {
            setPage(page);
            setPageSize(pageSize);
          },
          onShowSizeChange: (page: number, pageSize: number) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
        rowKey={(record: WorkflowModel) => record.id}
        scroll={{ x: "max(100%, 960px)" }}
      />
    </div>
  );
};

export default WorkflowList;
