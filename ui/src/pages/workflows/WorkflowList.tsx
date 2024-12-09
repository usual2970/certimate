import { useCallback, useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
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

import { Workflow as WorkflowType } from "@/domain/workflow";
import { list as listWorkflow, remove as removeWorkflow, save as saveWorkflow } from "@/repository/workflow";

const WorkflowList = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [modalApi, ModelContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const [loading, setLoading] = useState<boolean>(false);

  const tableColumns: TableProps<WorkflowType>["columns"] = [
    {
      key: "$index",
      align: "center",
      fixed: "left",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "name",
      title: t("common.text.name"),
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
      key: "type",
      title: t("workflow.props.executionMethod"),
      ellipsis: true,
      render: (_, record) => {
        const method = record.type;
        if (!method) {
          return "-";
        } else if (method === "manual") {
          return <Typography.Text>{t("workflow.node.start.form.executionMethod.options.manual")}</Typography.Text>;
        } else if (method === "auto") {
          return (
            <Space className="max-w-full" direction="vertical" size={4}>
              <Typography.Text>{t("workflow.node.start.form.executionMethod.options.auto")}</Typography.Text>
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
                setFilters((prev) => ({ ...prev, state: key }));
                setSelectedKeys([key]);
              }

              confirm({ closeDropdown: true });
            },
          };
        });

        const handleResetClick = () => {
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
                {t("common.reset")}
              </Button>
              <Button type="primary" size="small" onClick={handleConfirmClick}>
                {t("common.confirm")}
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
      title: t("common.text.created_at"),
      ellipsis: true,
      render: (_, record) => {
        return dayjs(record.created!).format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      key: "updatedAt",
      title: t("common.text.updated_at"),
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
          <Tooltip title={t("common.edit")}>
            <Button
              type="link"
              icon={<PencilIcon size={16} />}
              onClick={() => {
                navigate(`/workflows/detail?id=${record.id}`);
              }}
            />
          </Tooltip>
          <Tooltip title={t("common.delete")}>
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
  const [tableData, setTableData] = useState<WorkflowType[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [filters, setFilters] = useState<Record<string, unknown>>(() => {
    return {
      state: searchParams.get("state"),
    };
  });

  const [page, setPage] = useState<number>(() => parseInt(+searchParams.get("page")! + "") || 1);
  const [pageSize, setPageSize] = useState<number>(() => parseInt(+searchParams.get("perPage")! + "") || 10);

  const fetchTableData = useCallback(async () => {
    if (loading) return;
    setLoading(true);

    try {
      const resp = await listWorkflow({
        page: page,
        perPage: pageSize,
        enabled: (filters["state"] as string) === "enabled" ? true : (filters["state"] as string) === "disabled" ? false : undefined,
      });

      setTableData(resp.items);
      setTableTotal(resp.totalItems);
    } catch (err) {
      if (err instanceof ClientResponseError && err.isAbort) {
        return;
      }

      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
    } finally {
      setLoading(false);
    }
  }, [filters, page, pageSize]);

  useEffect(() => {
    fetchTableData();
  }, [fetchTableData]);

  const handleEnabledChange = async (workflow: WorkflowType) => {
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
      notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
    }
  };

  const handleDeleteClick = (workflow: WorkflowType) => {
    modalApi.confirm({
      title: t("workflow.action.delete.alert.title"),
      content: t("workflow.action.delete.alert.content"),
      onOk: async () => {
        try {
          const resp = await removeWorkflow(workflow);
          if (resp) {
            setTableData((prev) => prev.filter((item) => item.id !== workflow.id));
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
        }
      },
    });
  };

  const handleCreateClick = () => {
    navigate("/workflows/detail");
  };

  return (
    <>
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

      <Table<WorkflowType>
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
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
          onShowSizeChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
        rowKey={(record) => record.id}
        scroll={{ x: "max(100%, 960px)" }}
      />
    </>
  );
};

export default WorkflowList;
