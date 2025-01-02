import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  CheckCircleOutlined as CheckCircleOutlinedIcon,
  CloseCircleOutlined as CloseCircleOutlinedIcon,
  SelectOutlined as SelectOutlinedIcon,
} from "@ant-design/icons";
import { useRequest } from "ahooks";
import { Button, Empty, notification, Space, Table, theme, Typography, type TableProps } from "antd";
import { ClientResponseError } from "pocketbase";

import { type WorkflowRunModel } from "@/domain/workflowRun";
import { list as listWorkflowRuns } from "@/repository/workflowRun";
import { getErrMsg } from "@/utils/error";
import WorkflowRunDetailDrawer from "./WorkflowRunDetailDrawer";

export type WorkflowRunsProps = {
  className?: string;
  style?: React.CSSProperties;
  workflowId: string;
};

const WorkflowRuns = ({ className, style, workflowId }: WorkflowRunsProps) => {
  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const tableColumns: TableProps<WorkflowRunModel>["columns"] = [
    {
      key: "$index",
      align: "center",
      fixed: "left",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "id",
      title: t("workflow_run.props.id"),
      ellipsis: true,
      render: (_, record) => <span className="font-mono">{record.id}</span>,
    },
    {
      key: "status",
      title: t("workflow_run.props.status"),
      ellipsis: true,
      render: (_, record) => {
        if (record.succeed) {
          return (
            <Space>
              <CheckCircleOutlinedIcon style={{ color: themeToken.colorSuccess }} />
              <Typography.Text type="success">{t("workflow_run.props.status.succeeded")}</Typography.Text>
            </Space>
          );
        } else {
          return (
            <Space>
              <CloseCircleOutlinedIcon style={{ color: themeToken.colorError }} />
              <Typography.Text type="danger">{t("workflow_run.props.status.failed")}</Typography.Text>
            </Space>
          );
        }
      },
    },
    {
      key: "trigger",
      title: t("workflow_run.props.trigger"),
      ellipsis: true,
      render: () => {
        return "TODO";
      },
    },
    {
      key: "startedAt",
      title: t("workflow_run.props.started_at"),
      ellipsis: true,
      render: () => {
        return "TODO";
      },
    },
    {
      key: "completedAt",
      title: t("workflow_run.props.completed_at"),
      ellipsis: true,
      render: () => {
        return "TODO";
      },
    },
    {
      key: "$action",
      align: "end",
      fixed: "right",
      width: 120,
      render: (_, record) => (
        <Button.Group>
          <WorkflowRunDetailDrawer data={record} trigger={<Button color="primary" icon={<SelectOutlinedIcon />} variant="text" />} />
        </Button.Group>
      ),
    },
  ];
  const [tableData, setTableData] = useState<WorkflowRunModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);

  const { loading } = useRequest(
    () => {
      return listWorkflowRuns({
        workflowId: workflowId,
        page: page,
        perPage: pageSize,
      });
    },
    {
      refreshDeps: [workflowId, page, pageSize],
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

  return (
    <>
      {NotificationContextHolder}

      <div className={className} style={style}>
        <Table<WorkflowRunModel>
          columns={tableColumns}
          dataSource={tableData}
          loading={loading}
          locale={{
            emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />,
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
          rowKey={(record: WorkflowRunModel) => record.id}
          scroll={{ x: "max(100%, 960px)" }}
        />
      </div>
    </>
  );
};

export default WorkflowRuns;
