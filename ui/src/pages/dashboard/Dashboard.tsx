import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useRequest } from "ahooks";
import { Card, Col, Divider, notification, Row, Space, Statistic, theme, Typography } from "antd";
import { PageHeader } from "@ant-design/pro-components";
import {
  CalendarClock as CalendarClockIcon,
  CalendarX2 as CalendarX2Icon,
  FolderCheck as FolderCheckIcon,
  SquareSigma as SquareSigmaIcon,
  Workflow as WorkflowIcon,
} from "lucide-react";
import { ClientResponseError } from "pocketbase";

import { type Statistics } from "@/domain/statistics";
import { get as getStatistics } from "@/api/statistics";
import { getErrMsg } from "@/utils/error";

const Dashboard = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

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

      <div>TODO: {t("dashboard.latest_workflow_run")}</div>
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
