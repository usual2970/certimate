import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router-dom";
import {
  ApartmentOutlined as ApartmentOutlinedIcon,
  CaretRightOutlined as CaretRightOutlinedIcon,
  DeleteOutlined as DeleteOutlinedIcon,
  DownOutlined as DownOutlinedIcon,
  EllipsisOutlined as EllipsisOutlinedIcon,
  HistoryOutlined as HistoryOutlinedIcon,
  MinusOutlined,
  PlusCircleOutlined,
  ReloadOutlined,
  UndoOutlined as UndoOutlinedIcon,
} from "@ant-design/icons";
import { PageHeader } from "@ant-design/pro-components";
import { Alert, Button, Card, Dropdown, Form, Input, Modal, Space, Tabs, Typography, message, notification } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { isEqual } from "radash";
import { z } from "zod";

import { run as runWorkflow } from "@/api/workflows";
import ModalForm from "@/components/ModalForm";
import Show from "@/components/Show";
import WorkflowElements from "@/components/workflow/WorkflowElements";
import WorkflowRuns from "@/components/workflow/WorkflowRuns";
import { isAllNodesValidated } from "@/domain/workflow";
import { WORKFLOW_RUN_STATUSES } from "@/domain/workflowRun";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { remove as removeWorkflow, subscribe as subscribeWorkflow, unsubscribe as unsubscribeWorkflow } from "@/repository/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { getErrMsg } from "@/utils/error";

const WorkflowDetail = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [modalApi, ModalContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const [scale, setScale] = useState(1);

  const { id: workflowId } = useParams();
  const { workflow, initialized, ...workflowState } = useWorkflowStore(
    useZustandShallowSelector(["workflow", "initialized", "init", "destroy", "setEnabled", "release", "discard"])
  );
  useEffect(() => {
    // TODO: loading & error
    workflowState.init(workflowId!);

    return () => {
      workflowState.destroy();
    };
  }, [workflowId]);

  const [tabValue, setTabValue] = useState<"orchestration" | "runs">("orchestration");

  const [isRunning, setIsRunning] = useState(false);

  const [allowDiscard, setAllowDiscard] = useState(false);
  const [allowRelease, setAllowRelease] = useState(false);
  const [allowRun, setAllowRun] = useState(false);

  const lastRunStatus = useMemo(() => {
    return workflow.lastRunStatus;
  }, [workflow]);

  useEffect(() => {
    setIsRunning(lastRunStatus == WORKFLOW_RUN_STATUSES.RUNNING);
  }, [lastRunStatus]);

  useEffect(() => {
    if (!!workflowId && isRunning) {
      subscribeWorkflow(workflowId, (e) => {
        if (e.record.lastRunStatus !== WORKFLOW_RUN_STATUSES.RUNNING) {
          setIsRunning(false);
          unsubscribeWorkflow(workflowId);
        }
      });

      return () => {
        unsubscribeWorkflow(workflowId);
      };
    }
  }, [workflowId, isRunning]);

  useEffect(() => {
    const hasReleased = !!workflow.content;
    const hasChanges = workflow.hasDraft! || !isEqual(workflow.draft, workflow.content);
    setAllowDiscard(!isRunning && hasReleased && hasChanges);
    setAllowRelease(!isRunning && hasChanges);
    setAllowRun(hasReleased);
  }, [workflow.content, workflow.draft, workflow.hasDraft, isRunning]);

  const handleEnableChange = async () => {
    if (!workflow.enabled && (!workflow.content || !isAllNodesValidated(workflow.content))) {
      messageApi.warning(t("workflow.action.enable.failed.uncompleted"));
      return;
    }

    try {
      await workflowState.setEnabled(!workflow.enabled);
    } catch (err) {
      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    }
  };

  const handleDeleteClick = () => {
    modalApi.confirm({
      title: t("workflow.action.delete"),
      content: t("workflow.action.delete.confirm"),
      onOk: async () => {
        try {
          const resp = await removeWorkflow(workflow);
          if (resp) {
            navigate("/workflows", { replace: true });
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  const handleDiscardClick = () => {
    modalApi.confirm({
      title: t("workflow.detail.orchestration.action.discard"),
      content: t("workflow.detail.orchestration.action.discard.confirm"),
      onOk: async () => {
        try {
          await workflowState.discard();

          messageApi.success(t("common.text.operation_succeeded"));
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  const handleReleaseClick = () => {
    if (!isAllNodesValidated(workflow.draft!)) {
      messageApi.warning(t("workflow.detail.orchestration.action.release.failed.uncompleted"));
      return;
    }

    modalApi.confirm({
      title: t("workflow.detail.orchestration.action.release"),
      content: t("workflow.detail.orchestration.action.release.confirm"),
      onOk: async () => {
        try {
          await workflowState.release();

          messageApi.success(t("common.text.operation_succeeded"));
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  const handleRunClick = () => {
    const { promise, resolve, reject } = Promise.withResolvers();
    if (workflow.hasDraft) {
      modalApi.confirm({
        title: t("workflow.detail.orchestration.action.run"),
        content: t("workflow.detail.orchestration.action.run.confirm"),
        onOk: () => resolve(void 0),
        onCancel: () => reject(),
      });
    } else {
      resolve(void 0);
    }

    promise.then(async () => {
      let unsubscribeFn: Awaited<ReturnType<typeof subscribeWorkflow>> | undefined = undefined;

      try {
        setIsRunning(true);

        // subscribe before running workflow
        unsubscribeFn = await subscribeWorkflow(workflowId!, (e) => {
          if (e.record.lastRunStatus !== WORKFLOW_RUN_STATUSES.RUNNING) {
            setIsRunning(false);
            unsubscribeFn?.();
          }
        });

        await runWorkflow(workflowId!);

        messageApi.info(t("workflow.detail.orchestration.action.run.prompt"));
      } catch (err) {
        setIsRunning(false);
        unsubscribeFn?.();

        console.error(err);
        messageApi.warning(t("common.text.operation_failed"));
      }
    });
  };

  return (
    <div>
      {MessageContextHolder}
      {ModalContextHolder}
      {NotificationContextHolder}

      <Card styles={{ body: { padding: "0.5rem", paddingBottom: 0 } }}>
        <PageHeader
          style={{ paddingBottom: 0 }}
          title={workflow.name}
          extra={
            initialized
              ? [
                  <WorkflowBaseInfoModal key="edit" trigger={<Button>{t("common.button.edit")}</Button>} />,

                  <Button key="enable" onClick={handleEnableChange}>
                    {workflow.enabled ? t("workflow.action.disable") : t("workflow.action.enable")}
                  </Button>,

                  <Dropdown
                    key="more"
                    menu={{
                      items: [
                        {
                          key: "delete",
                          label: t("workflow.action.delete"),
                          danger: true,
                          icon: <DeleteOutlinedIcon />,
                          onClick: () => {
                            handleDeleteClick();
                          },
                        },
                      ],
                    }}
                    trigger={["click"]}
                  >
                    <Button icon={<DownOutlinedIcon />} iconPosition="end">
                      {t("common.button.more")}
                    </Button>
                  </Dropdown>,
                ]
              : []
          }
        >
          <Typography.Paragraph type="secondary">{workflow.description}</Typography.Paragraph>
          <Tabs
            activeKey={tabValue}
            defaultActiveKey="orchestration"
            items={[
              { key: "orchestration", label: t("workflow.detail.orchestration.tab"), icon: <ApartmentOutlinedIcon /> },
              { key: "runs", label: t("workflow.detail.runs.tab"), icon: <HistoryOutlinedIcon /> },
            ]}
            renderTabBar={(props, DefaultTabBar) => <DefaultTabBar {...props} style={{ margin: 0 }} />}
            tabBarStyle={{ border: "none" }}
            onChange={(key) => setTabValue(key as typeof tabValue)}
          />
        </PageHeader>
      </Card>

      <div className="p-4">
        <Card loading={!initialized}>
          <Show when={tabValue === "orchestration"}>
            <div className="relative">
              <div className="flex items-center justify-between gap-4">
                <div className="flex-1 overflow-hidden">
                  <Show when={workflow.hasDraft!}>
                    <Alert banner message={<div className="truncate">{t("workflow.detail.orchestration.draft.alert")}</div>} type="warning" />
                  </Show>
                </div>
                <div className="flex justify-end">
                  <Space>
                    <Button disabled={!allowRun} icon={<CaretRightOutlinedIcon />} loading={isRunning} type="primary" onClick={handleRunClick}>
                      {t("workflow.detail.orchestration.action.run")}
                    </Button>

                    <Button.Group>
                      <Button color="primary" disabled={!allowRelease} variant="outlined" onClick={handleReleaseClick}>
                        {t("workflow.detail.orchestration.action.release")}
                      </Button>

                      <Dropdown
                        menu={{
                          items: [
                            {
                              key: "discard",
                              disabled: !allowDiscard,
                              label: t("workflow.detail.orchestration.action.discard"),
                              icon: <UndoOutlinedIcon />,
                              onClick: handleDiscardClick,
                            },
                          ],
                        }}
                        trigger={["click"]}
                      >
                        <Button color="primary" disabled={!allowDiscard} icon={<EllipsisOutlinedIcon />} variant="outlined" />
                      </Dropdown>
                    </Button.Group>
                  </Space>
                </div>
              </div>
              <div className="fixed bottom-8 right-8 z-10 flex items-center gap-2 rounded-lg bg-white p-2 shadow-lg">
                <Button icon={<MinusOutlined />} disabled={scale <= 0.5} onClick={() => setScale((s) => Math.max(0.5, s - 0.1))} />
                <Typography.Text className="min-w-[3em] text-center">{Math.round(scale * 100)}%</Typography.Text>
                <Button icon={<PlusCircleOutlined />} disabled={scale >= 2} onClick={() => setScale((s) => Math.min(2, s + 0.1))} />
                <Button icon={<ReloadOutlined />} onClick={() => setScale(1)} />
              </div>

              <div className="size-full origin-top px-12 py-8 transition-transform duration-300 max-md:px-4" style={{ transform: `scale(${scale})` }}>
                <WorkflowElements />
              </div>
            </div>
          </Show>

          <Show when={tabValue === "runs"}>
            <WorkflowRuns workflowId={workflowId!} />
          </Show>
        </Card>
      </div>
    </div>
  );
};

const WorkflowBaseInfoModal = ({ trigger }: { trigger?: React.ReactNode }) => {
  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { workflow, ...workflowState } = useWorkflowStore(useZustandShallowSelector(["workflow", "setBaseInfo"]));

  const formSchema = z.object({
    name: z
      .string({ message: t("workflow.detail.baseinfo.form.name.placeholder") })
      .min(1, t("workflow.detail.baseinfo.form.name.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    description: z
      .string({ message: t("workflow.detail.baseinfo.form.description.placeholder") })
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim()
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const {
    form: formInst,
    formPending,
    formProps,
    submit: submitForm,
  } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: { name: workflow.name, description: workflow.description },
    onSubmit: async (values) => {
      try {
        await workflowState.setBaseInfo(values.name!, values.description!);
      } catch (err) {
        notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

        throw err;
      }
    },
  });

  const handleFormFinish = async () => {
    return submitForm();
  };

  return (
    <>
      {NotificationContextHolder}

      <ModalForm
        disabled={formPending}
        layout="vertical"
        form={formInst}
        modalProps={{ destroyOnClose: true }}
        okText={t("common.button.save")}
        title={t(`workflow.detail.baseinfo.modal.title`)}
        trigger={trigger}
        width={480}
        {...formProps}
        onFinish={handleFormFinish}
      >
        <Form.Item name="name" label={t("workflow.detail.baseinfo.form.name.label")} rules={[formRule]}>
          <Input placeholder={t("workflow.detail.baseinfo.form.name.placeholder")} />
        </Form.Item>

        <Form.Item name="description" label={t("workflow.detail.baseinfo.form.description.label")} rules={[formRule]}>
          <Input placeholder={t("workflow.detail.baseinfo.form.description.placeholder")} />
        </Form.Item>
      </ModalForm>
    </>
  );
};

export default WorkflowDetail;
