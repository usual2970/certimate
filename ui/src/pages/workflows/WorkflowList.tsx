import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router-dom";
import {
  CheckCircleOutlined as CheckCircleOutlinedIcon,
  ClockCircleOutlined as ClockCircleOutlinedIcon,
  CloseCircleOutlined as CloseCircleOutlinedIcon,
  DeleteOutlined as DeleteOutlinedIcon,
  EditOutlined as EditOutlinedIcon,
  PlusOutlined as PlusOutlinedIcon,
  ReloadOutlined as ReloadOutlinedIcon,
  StopOutlined as StopOutlinedIcon,
  SyncOutlined as SyncOutlinedIcon,
} from "@ant-design/icons";

import { PageHeader } from "@ant-design/pro-components";
import { useRequest } from "ahooks";
import {
  Button,
  Card,
  Divider,
  Empty,
  Flex,
  Input,
  Menu,
  type MenuProps,
  Modal,
  Radio,
  Space,
  Switch,
  Table,
  type TableProps,
  Tooltip,
  Typography,
  message,
  notification,
  theme,
} from "antd";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import { WORKFLOW_TRIGGERS, type WorkflowModel, isAllNodesValidated } from "@/domain/workflow";
import { WORKFLOW_RUN_STATUSES } from "@/domain/workflowRun";
import { list as listWorkflows, remove as removeWorkflow, save as saveWorkflow } from "@/repository/workflow";
import { getErrMsg } from "@/utils/error";

const WorkflowList = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [messageApi, MessageContextHolder] = message.useMessage();
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
        const trigger = record.trigger;
        if (!trigger) {
          return "-";
        } else if (trigger === WORKFLOW_TRIGGERS.MANUAL) {
          return <Typography.Text>{t("workflow.props.trigger.manual")}</Typography.Text>;
        } else if (trigger === WORKFLOW_TRIGGERS.AUTO) {
          return (
            <Space className="max-w-full" direction="vertical" size={4}>
              <Typography.Text>{t("workflow.props.trigger.auto")}</Typography.Text>
              <Typography.Text type="secondary">{record.triggerCron ?? ""}</Typography.Text>
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
            <Space className="w-full justify-end" style={{ padding: themeToken.paddingSM }}>
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
      key: "lastRun",
      title: t("workflow.props.last_run_at"),
      render: (_, record) => {
        let icon = <></>;
        if (record.lastRunStatus === WORKFLOW_RUN_STATUSES.PENDING) {
          icon = <ClockCircleOutlinedIcon style={{ color: themeToken.colorTextSecondary }} />;
        } else if (record.lastRunStatus === WORKFLOW_RUN_STATUSES.RUNNING) {
          icon = <SyncOutlinedIcon style={{ color: themeToken.colorInfo }} spin />;
        } else if (record.lastRunStatus === WORKFLOW_RUN_STATUSES.SUCCEEDED) {
          icon = <CheckCircleOutlinedIcon style={{ color: themeToken.colorSuccess }} />;
        } else if (record.lastRunStatus === WORKFLOW_RUN_STATUSES.FAILED) {
          icon = <CloseCircleOutlinedIcon style={{ color: themeToken.colorError }} />;
        } else if (record.lastRunStatus === WORKFLOW_RUN_STATUSES.CANCELED) {
          icon = <StopOutlinedIcon style={{ color: themeToken.colorWarning }} />;
        }

        return (
          <Space>
            {icon}
            <Typography.Text>{record.lastRunTime ? dayjs(record.lastRunTime!).format("YYYY-MM-DD HH:mm:ss") : ""}</Typography.Text>
          </Space>
        );
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
        <Space.Compact>
          <Tooltip title={t("workflow.action.edit")}>
            <Button
              color="primary"
              icon={<EditOutlinedIcon />}
              variant="text"
              onClick={() => {
                navigate(`/workflows/${record.id}`);
              }}
            />
          </Tooltip>

          <Tooltip title={t("workflow.action.delete")}>
            <Button
              color="danger"
              danger
              icon={<DeleteOutlinedIcon />}
              variant="text"
              onClick={() => {
                handleDeleteClick(record);
              }}
            />
          </Tooltip>
        </Space.Compact>
      ),
    },
  ];
  const [tableData, setTableData] = useState<WorkflowModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [filters, setFilters] = useState<Record<string, unknown>>(() => {
    return {
      keyword: searchParams.get("keyword"),
      state: searchParams.get("state"),
    };
  });

  const [page, setPage] = useState<number>(() => parseInt(+searchParams.get("page")! + "") || 1);
  const [pageSize, setPageSize] = useState<number>(() => parseInt(+searchParams.get("perPage")! + "") || 10);

  const {
    loading,
    error: loadedError,
    run: refreshData,
  } = useRequest(
    () => {
      return listWorkflows({
        keyword: filters["keyword"] as string,
        enabled: (filters["state"] as string) === "enabled" ? true : (filters["state"] as string) === "disabled" ? false : undefined,
        page: page,
        perPage: pageSize,
      });
    },
    {
      refreshDeps: [filters, page, pageSize],
      onSuccess: (res) => {
        setTableData(res.items);
        setTableTotal(res.totalItems);
      },
      onError: (err) => {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);
        notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

        throw err;
      },
    }
  );

  const handleSearch = (value: string) => {
    setFilters((prev) => ({ ...prev, keyword: value.trim() }));
    setPage(1);
  };

  const handleCreateClick = () => {
    navigate("/workflows/new");
  };

  const handleReloadClick = () => {
    if (loading) return;

    refreshData();
  };

  const handleEnabledChange = async (workflow: WorkflowModel) => {
    try {
      if (!workflow.enabled && (!workflow.content || !isAllNodesValidated(workflow.content))) {
        messageApi.warning(t("workflow.action.enable.failed.uncompleted"));
        return;
      }

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
          const resp = await removeWorkflow(workflow);
          if (resp) {
            setTableData((prev) => prev.filter((item) => item.id !== workflow.id));
            refreshData();
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  return (
    <div className="p-4">
      {MessageContextHolder}
      {ModelContextHolder}
      {NotificationContextHolder}

      <PageHeader
        title={t("workflow.page.title")}
        extra={[
          <Button
            key="create"
            type="primary"
            icon={<PlusOutlinedIcon />}
            onClick={() => {
              handleCreateClick();
            }}
          >
            {t("workflow.action.create")}
          </Button>,
        ]}
      />

      <Card size="small">
        <div className="mb-4">
          <Flex gap="small">
            <div className="flex-1">
              <Input.Search allowClear defaultValue={filters["keyword"] as string} placeholder={t("workflow.search.placeholder")} onSearch={handleSearch} />
            </div>
            <div>
              <Button icon={<ReloadOutlinedIcon spin={loading} />} onClick={handleReloadClick} />
            </div>
          </Flex>
        </div>

        <Table<WorkflowModel>
          columns={tableColumns}
          dataSource={tableData}
          loading={loading}
          locale={{
            emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={getErrMsg(loadedError ?? t("workflow.nodata"))} />,
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
          rowKey={(record) => record.id}
          scroll={{ x: "max(100%, 960px)" }}
        />
      </Card>
    </div>
  );
};

export default WorkflowList;
