import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import {
  ApiOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  LockOutlined,
  PlusOutlined,
  SelectOutlined,
  SendOutlined,
  SyncOutlined,
} from "@ant-design/icons";
import { PageHeader } from "@ant-design/pro-components";
import { useRequest } from "ahooks";
import type { TableProps } from "antd";
import { Button, Card, Col, Divider, Empty, Flex, Grid, Row, Space, Statistic, Table, Tag, Typography, notification, theme } from "antd";
import dayjs from "dayjs";
import {
  CalendarClock as CalendarClockIcon,
  CalendarX2 as CalendarX2Icon,
  FolderCheck as FolderCheckIcon,
  SquareSigma as SquareSigmaIcon,
  Workflow as WorkflowIcon,
} from "lucide-react";
import { ClientResponseError } from "pocketbase";

import { get as getStatistics } from "@/api/statistics";
import WorkflowRunDetailDrawer from "@/components/workflow/WorkflowRunDetailDrawer";
import { type Statistics } from "@/domain/statistics";
import { WORKFLOW_TRIGGERS } from "@/domain/workflow";
import { WORKFLOW_RUN_STATUSES, type WorkflowRunModel } from "@/domain/workflowRun";
import { list as listWorkflowRuns } from "@/repository/workflowRun";
import { getErrMsg } from "@/utils/error";

const { useBreakpoint } = Grid;

const Dashboard = () => {
  const navigate = useNavigate();

  const screens = useBreakpoint();

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
      key: "name",
      title: t("workflow.props.name"),
      ellipsis: true,
      render: (_, record) => (
        <Space className="max-w-full" direction="vertical" size={4}>
          <Typography.Text ellipsis>{record.expand?.workflowId?.name}</Typography.Text>
          <Typography.Text type="secondary" ellipsis>
            {record.expand?.workflowId?.description}
          </Typography.Text>
        </Space>
      ),
    },
    {
      key: "status",
      title: t("workflow_run.props.status"),
      ellipsis: true,
      render: (_, record) => {
        if (record.status === WORKFLOW_RUN_STATUSES.PENDING) {
          return <Tag icon={<ClockCircleOutlined />}>{t("workflow_run.props.status.pending")}</Tag>;
        } else if (record.status === WORKFLOW_RUN_STATUSES.RUNNING) {
          return (
            <Tag icon={<SyncOutlined spin />} color="processing">
              {t("workflow_run.props.status.running")}
            </Tag>
          );
        } else if (record.status === WORKFLOW_RUN_STATUSES.SUCCEEDED) {
          return (
            <Tag icon={<CheckCircleOutlined />} color="success">
              {t("workflow_run.props.status.succeeded")}
            </Tag>
          );
        } else if (record.status === WORKFLOW_RUN_STATUSES.FAILED) {
          return (
            <Tag icon={<CloseCircleOutlined />} color="error">
              {t("workflow_run.props.status.failed")}
            </Tag>
          );
        }

        return <></>;
      },
    },
    {
      key: "trigger",
      title: t("workflow_run.props.trigger"),
      ellipsis: true,
      render: (_, record) => {
        if (record.trigger === WORKFLOW_TRIGGERS.AUTO) {
          return t("workflow_run.props.trigger.auto");
        } else if (record.trigger === WORKFLOW_TRIGGERS.MANUAL) {
          return t("workflow_run.props.trigger.manual");
        }

        return <></>;
      },
    },
    {
      key: "startedAt",
      title: t("workflow_run.props.started_at"),
      ellipsis: true,
      render: (_, record) => {
        if (record.startedAt) {
          return dayjs(record.startedAt).format("YYYY-MM-DD HH:mm:ss");
        }

        return <></>;
      },
    },
    {
      key: "endedAt",
      title: t("workflow_run.props.ended_at"),
      ellipsis: true,
      render: (_, record) => {
        if (record.endedAt) {
          return dayjs(record.endedAt).format("YYYY-MM-DD HH:mm:ss");
        }

        return <></>;
      },
    },
    {
      key: "$action",
      align: "end",
      fixed: "right",
      width: 120,
      render: (_, record) => (
        <Button.Group>
          <WorkflowRunDetailDrawer data={record} trigger={<Button color="primary" icon={<SelectOutlined />} variant="text" />} />
        </Button.Group>
      ),
    },
  ];
  const [tableData, setTableData] = useState<WorkflowRunModel[]>([]);
  const [_tableTotal, setTableTotal] = useState<number>(0);

  const [page, _setPage] = useState<number>(1);
  const [pageSize, _setPageSize] = useState<number>(3);

  const { loading: loadingWorkflowRun } = useRequest(
    () => {
      return listWorkflowRuns({
        page: page,
        perPage: pageSize,
        expand: true,
      });
    },
    {
      refreshDeps: [page, pageSize],
      onSuccess: (data) => {
        setTableData(data.items);
        setTableTotal(data.totalItems > 3 ? 3 : data.totalItems);
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

  const statisticsGridSpans = {
    xs: { flex: "50%" },
    md: { flex: "50%" },
    lg: { flex: "33.3333%" },
    xl: { flex: "33.3333%" },
    xxl: { flex: "20%" },
  };
  const [statistics, setStatistics] = useState<Statistics>();

  const { loading } = useRequest(
    () => {
      return getStatistics();
    },
    {
      onSuccess: (data) => {
        setStatistics(data);
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
    <div className="p-4">
      {NotificationContextHolder}

      <PageHeader title={t("dashboard.page.title")} />

      <Row className="justify-stretch" gutter={[16, 16]}>
        <Col {...statisticsGridSpans}>
          <StatisticCard
            icon={<SquareSigmaIcon size={48} strokeWidth={1} color={themeToken.colorInfo} />}
            label={t("dashboard.statistics.all_certificates")}
            loading={loading}
            value={statistics?.certificateTotal ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/certificates")}
          />
        </Col>
        <Col {...statisticsGridSpans}>
          <StatisticCard
            icon={<CalendarClockIcon size={48} strokeWidth={1} color={themeToken.colorWarning} />}
            label={t("dashboard.statistics.expire_soon_certificates")}
            loading={loading}
            value={statistics?.certificateExpireSoon ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/certificates?state=expireSoon")}
          />
        </Col>
        <Col {...statisticsGridSpans}>
          <StatisticCard
            icon={<CalendarX2Icon size={48} strokeWidth={1} color={themeToken.colorError} />}
            label={t("dashboard.statistics.expired_certificates")}
            loading={loading}
            value={statistics?.certificateExpired ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/certificates?state=expired")}
          />
        </Col>
        <Col {...statisticsGridSpans}>
          <StatisticCard
            icon={<WorkflowIcon size={48} strokeWidth={1} color={themeToken.colorInfo} />}
            label={t("dashboard.statistics.all_workflows")}
            loading={loading}
            value={statistics?.workflowTotal ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/workflows")}
          />
        </Col>
        <Col {...statisticsGridSpans}>
          <StatisticCard
            icon={<FolderCheckIcon size={48} strokeWidth={1} color={themeToken.colorSuccess} />}
            label={t("dashboard.statistics.enabled_workflows")}
            loading={loading}
            value={statistics?.workflowEnabled ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/workflows?state=enabled")}
          />
        </Col>
      </Row>

      <Divider />

      <Flex vertical={!screens.md} gap={16}>
        <Card className="sm:h-full sm:w-[500px] sm:pb-32">
          <div className="text-lg font-semibold">{t("dashboard.quick_actions")}</div>
          <div className="mt-9">
            <Button className="w-full" type="primary" size="large" icon={<PlusOutlined />} onClick={() => navigate("/workflows/new")}>
              {t("dashboard.quick_actions.create_workflow")}
            </Button>
            <Button className="mt-5 w-full" size="large" icon={<LockOutlined />} onClick={() => navigate("/settings/password")}>
              {t("dashboard.quick_actions.change_login_password")}
            </Button>
            <Button className="mt-5 w-full" size="large" icon={<SendOutlined />} onClick={() => navigate("/settings/notification")}>
              {t("dashboard.quick_actions.notification_settings")}
            </Button>
            <Button className="mt-5 w-full" size="large" icon={<ApiOutlined />} onClick={() => navigate("/settings/ssl-provider")}>
              {t("dashboard.quick_actions.certificate_authority_configuration")}
            </Button>
          </div>
        </Card>
        <Card className="size-full">
          <div className="text-lg font-semibold">{t("dashboard.latest_workflow_run")} </div>
          <Table<WorkflowRunModel>
            className="mt-5"
            columns={tableColumns}
            dataSource={tableData}
            loading={loadingWorkflowRun}
            locale={{
              emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />,
            }}
            rowKey={(record: WorkflowRunModel) => record.id}
            scroll={{ x: "max(100%, 960px)" }}
          />
        </Card>
      </Flex>
    </div>
  );
};

const StatisticCard = ({
  label,
  loading,
  icon,
  value,
  suffix,
  onClick,
}: {
  label: React.ReactNode;
  loading?: boolean;
  icon: React.ReactNode;
  value?: string | number | React.ReactNode;
  suffix?: React.ReactNode;
  onClick?: () => void;
}) => {
  return (
    <Card className="size-full overflow-hidden" bordered={false} hoverable loading={loading} onClick={onClick}>
      <Space size="middle">
        {icon}
        <Statistic
          title={label}
          valueRender={() => {
            return <Typography.Text className="text-4xl">{value}</Typography.Text>;
          }}
          suffix={<Typography.Text className="text-sm">{suffix}</Typography.Text>}
        />
      </Space>
    </Card>
  );
};

export default Dashboard;
