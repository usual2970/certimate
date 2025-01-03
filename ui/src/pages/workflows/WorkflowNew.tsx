import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { PageHeader } from "@ant-design/pro-components";
import { Card, Col, Row, Spin, Typography, notification } from "antd";
import { sleep } from "radash";

import { type WorkflowModel, initWorkflow } from "@/domain/workflow";
import { save as saveWorkflow } from "@/repository/workflow";
import { getErrMsg } from "@/utils/error";

const TEMPLATE_KEY_BLANK = "blank" as const;
const TEMPLATE_KEY_STANDARD = "standard" as const;
type TemplateKeys = typeof TEMPLATE_KEY_BLANK | typeof TEMPLATE_KEY_STANDARD;

const WorkflowNew = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const templateGridSpans = {
    xs: { flex: "100%" },
    md: { flex: "100%" },
    lg: { flex: "50%" },
    xl: { flex: "50%" },
    xxl: { flex: "50%" },
  };
  const [templateSelectKey, setTemplateSelectKey] = useState<TemplateKeys>();

  const handleTemplateSelect = async (key: TemplateKeys) => {
    if (templateSelectKey) return;

    setTemplateSelectKey(key);

    try {
      let workflow: WorkflowModel;

      switch (key) {
        case TEMPLATE_KEY_BLANK:
          workflow = initWorkflow();
          break;

        case TEMPLATE_KEY_STANDARD:
          workflow = initWorkflow({ template: "standard" });
          break;

        default:
          throw "Invalid args: `key`";
      }

      workflow = await saveWorkflow(workflow);
      await sleep(500);

      await navigate(`/workflows/${workflow.id}`, { replace: true });
    } catch (err) {
      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      setTemplateSelectKey(undefined);
    }
  };

  return (
    <div>
      {NotificationContextHolder}

      <Card styles={{ body: { padding: "0.5rem", paddingBottom: 0 } }}>
        <PageHeader title={t("workflow.new.title")}>
          <Typography.Paragraph type="secondary">{t("workflow.new.subtitle")}</Typography.Paragraph>
        </PageHeader>
      </Card>

      <div className="p-4">
        <div className="mx-auto max-w-[960px] px-2">
          <Typography.Text type="secondary">
            <div className="mb-8 mt-4 text-xl">{t("workflow.new.templates.title")}</div>
          </Typography.Text>

          <Row className="justify-stretch" gutter={[16, 16]}>
            <Col {...templateGridSpans}>
              <Card
                className="size-full"
                cover={<img className="min-h-[120px] object-contain" src="/imgs/workflow/tpl-standard.png" />}
                hoverable
                onClick={() => handleTemplateSelect(TEMPLATE_KEY_STANDARD)}
              >
                <div className="flex w-full items-center gap-4">
                  <Card.Meta
                    className="grow"
                    title={t("workflow.new.templates.template.standard.title")}
                    description={t("workflow.new.templates.template.standard.description")}
                  />
                  <Spin spinning={templateSelectKey === TEMPLATE_KEY_STANDARD} />
                </div>
              </Card>
            </Col>
            <Col {...templateGridSpans}>
              <Card
                className="size-full"
                cover={<img className="min-h-[120px] object-contain" src="/imgs/workflow/tpl-blank.png" />}
                hoverable
                onClick={() => handleTemplateSelect(TEMPLATE_KEY_BLANK)}
              >
                <div className="flex w-full items-center gap-4">
                  <Card.Meta
                    className="grow"
                    title={t("workflow.new.templates.template.blank.title")}
                    description={t("workflow.new.templates.template.blank.description")}
                  />
                  <Spin spinning={templateSelectKey === TEMPLATE_KEY_BLANK} />
                </div>
              </Card>
            </Col>
          </Row>
        </div>
      </div>
    </div>
  );
};

export default WorkflowNew;
