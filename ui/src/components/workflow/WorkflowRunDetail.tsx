import { useState } from "react";
import { useTranslation } from "react-i18next";
import { RightOutlined as RightOutlinedIcon, SelectOutlined as SelectOutlinedIcon } from "@ant-design/icons";
import { useRequest } from "ahooks";
import { Alert, Button, Collapse, Divider, Empty, Skeleton, Space, Spin, Table, type TableProps, Tooltip, Typography, notification } from "antd";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import CertificateDetailDrawer from "@/components/certificate/CertificateDetailDrawer";
import Show from "@/components/Show";
import { type CertificateModel } from "@/domain/certificate";
import type { WorkflowLogModel } from "@/domain/workflowLog";
import { WORKFLOW_RUN_STATUSES, type WorkflowRunModel } from "@/domain/workflowRun";
import { listByWorkflowRunId as listCertificatesByWorkflowRunId } from "@/repository/certificate";
import { listByWorkflowRunId as listLogsByWorkflowRunId } from "@/repository/workflowLog";
import { mergeCls } from "@/utils/css";
import { getErrMsg } from "@/utils/error";

export type WorkflowRunDetailProps = {
  className?: string;
  style?: React.CSSProperties;
  data: WorkflowRunModel;
};

const WorkflowRunDetail = ({ data, ...props }: WorkflowRunDetailProps) => {
  const { t } = useTranslation();

  return (
    <div {...props}>
      <Show when={data.status === WORKFLOW_RUN_STATUSES.SUCCEEDED}>
        <Alert showIcon type="success" message={<Typography.Text type="success">{t("workflow_run.props.status.succeeded")}</Typography.Text>} />
      </Show>

      <Show when={data.status === WORKFLOW_RUN_STATUSES.FAILED}>
        <Alert showIcon type="error" message={<Typography.Text type="danger">{t("workflow_run.props.status.failed")}</Typography.Text>} />
      </Show>

      <div className="my-4">
        <WorkflowRunLogs runId={data.id} runStatus={data.status} />
      </div>

      <Show when={data.status === WORKFLOW_RUN_STATUSES.SUCCEEDED}>
        <Divider />

        <WorkflowRunArtifacts runId={data.id} />
      </Show>
    </div>
  );
};

const WorkflowRunLogs = ({ runId, runStatus }: { runId: string; runStatus: string }) => {
  const { t } = useTranslation();

  type Log = Pick<WorkflowLogModel, "level" | "message" | "data" | "created">;
  type LogGroup = { id: string; name: string; records: Log[] };

  const [listData, setListData] = useState<LogGroup[]>([]);
  const { loading } = useRequest(
    () => {
      return listLogsByWorkflowRunId(runId);
    },
    {
      refreshDeps: [runId, runStatus],
      pollingInterval: runStatus === WORKFLOW_RUN_STATUSES.PENDING || runStatus === WORKFLOW_RUN_STATUSES.RUNNING ? 5000 : 0,
      pollingWhenHidden: false,
      throttleWait: 500,
      onBefore: () => {
        setListData([]);
      },
      onSuccess: (res) => {
        setListData(
          res.items.reduce((acc, e) => {
            let group = acc.at(-1);
            if (!group || group.id !== e.nodeId) {
              group = { id: e.nodeId, name: e.nodeName, records: [] };
              acc.push(group);
            }
            group.records.push({ level: e.level, message: e.message, data: e.data, created: e.created });
            return acc;
          }, [] as LogGroup[])
        );
      },
      onError: (err) => {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);

        throw err;
      },
    }
  );

  const renderLogRecord = (record: Log) => {
    let message = <>{record.message}</>;
    if (record.data != null && Object.keys(record.data).length > 0) {
      message = (
        <details>
          <summary>{record.message}</summary>
          {Object.entries(record.data).map(([key, value]) => (
            <div key={key} className="flex space-x-2 " style={{ wordBreak: "break-word" }}>
              <div>{key}:</div>
              <div>{JSON.stringify(value)}</div>
            </div>
          ))}
        </details>
      );
    }

    return (
      <div className="flex space-x-2 text-xs" style={{ wordBreak: "break-word" }}>
        <div className="whitespace-nowrap text-stone-400">[{dayjs(record.created).format("YYYY-MM-DD HH:mm:ss")}]</div>
        <div
          className={mergeCls(
            "font-mono",
            record.level === "DEBUG" ? "text-stone-400" : "",
            record.level === "WARN" ? "text-yellow-600" : "",
            record.level === "ERROR" ? "text-red-600" : ""
          )}
        >
          {message}
        </div>
      </div>
    );
  };

  return (
    <>
      <Typography.Title level={5}>{t("workflow_run.logs")}</Typography.Title>
      <div className="rounded-md bg-black text-stone-200">
        <Show
          when={!loading || listData.length > 0}
          fallback={
            <Spin spinning>
              <Skeleton />
            </Spin>
          }
        >
          <div className=" py-2">
            <Collapse
              style={{ color: "inherit" }}
              bordered={false}
              defaultActiveKey={listData.map((group) => group.id)}
              expandIcon={({ isActive }) => <RightOutlinedIcon rotate={isActive ? 90 : 0} />}
              items={listData.map((group) => {
                return {
                  key: group.id,
                  classNames: {
                    header: "text-sm text-stone-200",
                    body: "text-stone-200",
                  },
                  style: { color: "inherit", border: "none" },
                  styles: {
                    header: { color: "inherit" },
                  },
                  label: group.name,
                  children: <div className="flex flex-col space-y-1">{group.records.map((record) => renderLogRecord(record))}</div>,
                };
              })}
            />
          </div>
        </Show>
      </div>
    </>
  );
};

const WorkflowRunArtifacts = ({ runId }: { runId: string }) => {
  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const tableColumns: TableProps<CertificateModel>["columns"] = [
    {
      key: "$index",
      align: "center",
      fixed: "left",
      width: 50,
      render: (_, __, index) => index + 1,
    },
    {
      key: "type",
      title: t("workflow_run_artifact.props.type"),
      render: () => t("workflow_run_artifact.props.type.certificate"),
    },
    {
      key: "name",
      title: t("workflow_run_artifact.props.name"),
      ellipsis: true,
      render: (_, record) => {
        return (
          <Typography.Text delete={!!record.deleted} ellipsis>
            {record.subjectAltNames}
          </Typography.Text>
        );
      },
    },
    {
      key: "$action",
      align: "end",
      width: 120,
      render: (_, record) => (
        <Space.Compact>
          <CertificateDetailDrawer
            data={record}
            trigger={
              <Tooltip title={t("certificate.action.view")}>
                <Button color="primary" disabled={!!record.deleted} icon={<SelectOutlinedIcon />} variant="text" />
              </Tooltip>
            }
          />
        </Space.Compact>
      ),
    },
  ];
  const [tableData, setTableData] = useState<CertificateModel[]>([]);
  const { loading: tableLoading } = useRequest(
    () => {
      return listCertificatesByWorkflowRunId(runId);
    },
    {
      refreshDeps: [runId],
      onBefore: () => {
        setTableData([]);
      },
      onSuccess: (res) => {
        setTableData(res.items);
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

  return (
    <>
      {NotificationContextHolder}

      <Typography.Title level={5}>{t("workflow_run.artifacts")}</Typography.Title>
      <Table<CertificateModel>
        columns={tableColumns}
        dataSource={tableData}
        loading={tableLoading}
        locale={{
          emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />,
        }}
        pagination={false}
        rowKey={(record) => record.id}
        size="small"
      />
    </>
  );
};

export default WorkflowRunDetail;
