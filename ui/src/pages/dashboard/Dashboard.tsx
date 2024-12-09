import React, { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
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

import { type Statistic as StatisticType } from "@/domain/domain";
import { get as getStatistics } from "@/api/statistics";

const Dashboard = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const [loading, setLoading] = useState<boolean>(false);

  const statisticGridSpans = {
    xs: { flex: "50%" },
    md: { flex: "50%" },
    lg: { flex: "33.3333%" },
    xl: { flex: "33.3333%" },
    xxl: { flex: "20%" },
  };

  const [statistic, setStatistic] = useState<StatisticType>();

  const fetchStatistic = useCallback(async () => {
    if (loading) return;
    setLoading(true);

    try {
      const data = await getStatistics();
      setStatistic(data);
    } catch (err) {
      if (err instanceof ClientResponseError && err.isAbort) {
        return;
      }

      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStatistic();
  }, []);

  return (
    <>
      {NotificationContextHolder}

      <PageHeader title={t("dashboard.page.title")} />

      <Row className="justify-stretch" gutter={[16, 16]}>
        <Col {...statisticGridSpans}>
          <StatisticCard
            icon={<SquareSigmaIcon size={48} strokeWidth={1} color={themeToken.colorInfo} />}
            label={t("dashboard.statistics.all_certificates")}
            value={statistic?.certificateTotal ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/certificates")}
          />
        </Col>
        <Col {...statisticGridSpans}>
          <StatisticCard
            icon={<CalendarClockIcon size={48} strokeWidth={1} color={themeToken.colorWarning} />}
            label={t("dashboard.statistics.expire_soon_certificates")}
            value={statistic?.certificateExpireSoon ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/certificates?state=expireSoon")}
          />
        </Col>
        <Col {...statisticGridSpans}>
          <StatisticCard
            icon={<CalendarX2Icon size={48} strokeWidth={1} color={themeToken.colorError} />}
            label={t("dashboard.statistics.expired_certificates")}
            value={statistic?.certificateExpired ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/certificates?state=expired")}
          />
        </Col>
        <Col {...statisticGridSpans}>
          <StatisticCard
            icon={<WorkflowIcon size={48} strokeWidth={1} color={themeToken.colorInfo} />}
            label={t("dashboard.statistics.all_workflows")}
            value={statistic?.workflowTotal ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/workflows")}
          />
        </Col>
        <Col {...statisticGridSpans}>
          <StatisticCard
            icon={<FolderCheckIcon size={48} strokeWidth={1} color={themeToken.colorSuccess} />}
            label={t("dashboard.statistics.enabled_workflows")}
            value={statistic?.workflowEnabled ?? "-"}
            suffix={t("dashboard.statistics.unit")}
            onClick={() => navigate("/workflows?state=enabled")}
          />
        </Col>
      </Row>

      <Divider />

      <div>TODO: 最近执行的工作流 LatestWorkflowRun</div>
    </>
  );
};

const StatisticCard = ({
  label,
  icon,
  value,
  suffix,
  onClick,
}: {
  label: React.ReactNode;
  icon: React.ReactNode;
  value?: string | number | React.ReactNode;
  suffix?: React.ReactNode;
  onClick?: () => void;
}) => {
  return (
    <Card className="size-full overflow-hidden" bordered={false} hoverable onClick={onClick}>
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
