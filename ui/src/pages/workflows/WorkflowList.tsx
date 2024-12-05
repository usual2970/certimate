import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Modal, notification, Space, Switch, Table, Tooltip, Typography, type TableProps } from "antd";
import { PageHeader } from "@ant-design/pro-components";
import { Pencil as PencilIcon, Plus as PlusIcon, Trash2 as Trash2Icon } from "lucide-react";

import { Workflow as WorkflowType } from "@/domain/workflow";
import { list as listWorkflow, remove as removeWorkflow, save as saveWorkflow, type WorkflowListReq } from "@/repository/workflow";

const WorkflowList = () => {
  const { t } = useTranslation();

  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [modalApi, ModelContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const [loading, setLoading] = useState<boolean>(false);

  const tableColumns: TableProps<WorkflowType>["columns"] = [
    {
      key: "$index",
      align: "center",
      title: "",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "name",
      title: t("common.text.name"),
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
      key: "enabled",
      title: t("workflow.props.enabled"),
      render: (_, record) => {
        const enabled = record.enabled;
        return (
          <>
            <Switch
              checked={enabled}
              onChange={() => {
                handleEnabledChange(record.id);
              }}
            />
          </>
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
        return new Date(record.created!).toLocaleString();
      },
    },
    {
      key: "updatedAt",
      title: t("common.text.updated_at"),
      ellipsis: true,
      render: (_, record) => {
        return new Date(record.updated!).toLocaleString();
      },
    },
    {
      key: "$operations",
      align: "end",
      width: 100,
      render: (_, record) => (
        <Space>
          <Tooltip title={t("common.edit")}>
            <Button
              type="link"
              icon={<PencilIcon size={16} />}
              onClick={() => {
                navigate(`/workflow/detail?id=${record.id}`);
              }}
            />
          </Tooltip>
          <Tooltip title={t("common.delete")}>
            <Button
              type="link"
              danger={true}
              icon={<Trash2Icon size={16} />}
              onClick={() => {
                handleDeleteClick(record.id);
              }}
            />
          </Tooltip>
        </Space>
      ),
    },
  ];
  const [tableData, setTableData] = useState<WorkflowType[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);

  // TODO: 表头筛选
  const fetchTableData = async () => {
    if (loading) return;
    setLoading(true);

    const state = searchParams.get("state");
    const req: WorkflowListReq = { page: page, perPage: pageSize };
    if (state == "enabled") {
      req.enabled = true;
    }

    try {
      const resp = await listWorkflow(req);

      setTableData(resp.items);
      setTableTotal(resp.totalItems);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTableData();
  }, [page, pageSize]);

  const handleEnabledChange = async (id: string) => {
    try {
      const resp = await saveWorkflow({ id, enabled: !tableData.find((item) => item.id === id)?.enabled });
      if (resp) {
        setTableData((prev) => {
          return prev.map((item) => {
            if (item.id === id) {
              return resp;
            }
            return item;
          });
        });
      }
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
    }
  };

  const handleDeleteClick = (id: string) => {
    modalApi.confirm({
      title: t("workflow.action.delete.alert.title"),
      content: t("workflow.action.delete.alert.content"),
      onOk: async () => {
        try {
          const resp = await removeWorkflow(id);
          if (resp) {
            setTableData((prev) => prev.filter((item) => item.id !== id));
          }
        } catch (err) {
          notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
        }
      },
    });
  };

  const handleCreateClick = () => {
    navigate("/workflow/detail");
  };

  // TODO: Empty 样式

  return (
    <>
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
        rowKey={(record) => record.id}
        loading={loading}
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
      />

      {ModelContextHolder}
      {NotificationContextHolder}
    </>
  );
};

export default WorkflowList;
