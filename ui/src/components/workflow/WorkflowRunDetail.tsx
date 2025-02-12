import { useState } from "react";
import { useTranslation } from "react-i18next";
import { SelectOutlined as SelectOutlinedIcon } from "@ant-design/icons";
import { useRequest } from "ahooks";
import { Alert, Button, Divider, Empty, Table, type TableProps, Tooltip, Typography, notification } from "antd";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import CertificateDetailDrawer from "@/components/certificate/CertificateDetailDrawer";
import Show from "@/components/Show";
import { type CertificateModel } from "@/domain/certificate";
import { WORKFLOW_RUN_STATUSES, type WorkflowRunModel } from "@/domain/workflowRun";
import { listByWorkflowRunId as listCertificateByWorkflowRunId } from "@/repository/certificate";
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
        <Typography.Title level={5}>{t("workflow_run.logs")}</Typography.Title>
        <div className="rounded-md bg-black p-4 text-stone-200">
          <div className="flex flex-col space-y-4">
            {data.logs?.map((item, i) => {
              return (
                <div key={i} className="flex flex-col space-y-2">
                  <div className="font-semibold">{item.nodeName}</div>
                  <div className="flex flex-col space-y-1">
                    {item.records?.map((output, j) => {
                      return (
                        <div key={j} className="flex space-x-2 text-sm" style={{ wordBreak: "break-word" }}>
                          <div className="whitespace-nowrap">[{dayjs(output.time).format("YYYY-MM-DD HH:mm:ss")}]</div>
                          {output.error ? <div className="text-red-500">{output.error}</div> : <div>{output.content}</div>}
                        </div>
                      );
                    })}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      <Show when={data.status === WORKFLOW_RUN_STATUSES.SUCCEEDED}>
        <Divider />

        <WorkflowRunArtifacts runId={data.id} />
      </Show>
    </div>
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
        <Button.Group>
          <CertificateDetailDrawer
            data={record}
            trigger={
              <Tooltip title={t("certificate.action.view")}>
                <Button color="primary" disabled={!!record.deleted} icon={<SelectOutlinedIcon />} variant="text" />
              </Tooltip>
            }
          />
        </Button.Group>
      ),
    },
  ];
  const [tableData, setTableData] = useState<CertificateModel[]>([]);
  const { loading: tableLoading } = useRequest(
    () => {
      return listCertificateByWorkflowRunId(runId);
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
