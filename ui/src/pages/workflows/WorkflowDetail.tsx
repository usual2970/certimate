import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router-dom";
import {
  ApartmentOutlined as ApartmentOutlinedIcon,
  CaretRightOutlined as CaretRightOutlinedIcon,
  DeleteOutlined as DeleteOutlinedIcon,
  DownOutlined as DownOutlinedIcon,
  EllipsisOutlined as EllipsisOutlinedIcon,
  HistoryOutlined as HistoryOutlinedIcon,
  UndoOutlined as UndoOutlinedIcon,
} from "@ant-design/icons";
import { PageHeader } from "@ant-design/pro-components";
import { useDeepCompareEffect } from "ahooks";
import { Button, Card, Dropdown, Form, Input, Modal, Space, Tabs, Typography, message, notification } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { ClientResponseError } from "pocketbase";
import { isEqual } from "radash";
import { z } from "zod";

import { run as runWorkflow } from "@/api/workflow";
import ModalForm from "@/components/ModalForm";
import Show from "@/components/Show";
import WorkflowElements from "@/components/workflow/WorkflowElements";
import WorkflowRuns from "@/components/workflow/WorkflowRuns";
import { type WorkflowModel, isAllNodesValidated } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { remove as removeWorkflow } from "@/repository/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { getErrMsg } from "@/utils/error";

const WorkflowDetail = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [modalApi, ModalContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { id: workflowId } = useParams();
  const { workflow, initialized, init, save, destroy, setBaseInfo, switchEnable } = useWorkflowStore(
    useZustandShallowSelector(["workflow", "initialized", "init", "destroy", "save", "setBaseInfo", "switchEnable"])
  );
  useEffect(() => {
    // TODO: loading & error
    init(workflowId!);

    return () => {
      destroy();
    };
  }, [workflowId]);

  const [tabValue, setTabValue] = useState<"orchestration" | "runs">("orchestration");

  const [isRunning, setIsRunning] = useState(false);

  const [allowDiscard, setAllowDiscard] = useState(false);
  const [allowRelease, setAllowRelease] = useState(false);
  const [allowRun, setAllowRun] = useState(false);
  useDeepCompareEffect(() => {
    const hasReleased = !!workflow.content;
    const hasChanges = workflow.hasDraft! || !isEqual(workflow.draft, workflow.content);
    setAllowDiscard(!isRunning && hasReleased && hasChanges);
    setAllowRelease(!isRunning && hasChanges);
    setAllowRun(hasReleased);
  }, [workflow.content, workflow.draft, workflow.hasDraft, isRunning]);

  const handleBaseInfoFormFinish = async (values: Pick<WorkflowModel, "name" | "description">) => {
    try {
      await setBaseInfo(values.name!, values.description!);
    } catch (err) {
      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
      return false;
    }
  };

  const handleEnableChange = async () => {
    if (!workflow.enabled && (!workflow.content || !isAllNodesValidated(workflow.content))) {
      messageApi.warning(t("workflow.action.enable.failed.uncompleted"));
      return;
    }

    try {
      await switchEnable();
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
          const resp: boolean = await removeWorkflow(workflow);
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
      onOk: () => {
        alert("TODO");
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
          await save();

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

    // TODO: 异步执行
    promise.then(async () => {
      setIsRunning(true);

      try {
        await runWorkflow(workflowId!);

        messageApi.success(t("common.text.operation_succeeded"));
      } catch (err) {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);
        messageApi.warning(t("common.text.operation_failed"));
      } finally {
        setIsRunning(false);
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
                  <WorkflowBaseInfoModalForm
                    key="edit"
                    data={workflow}
                    trigger={<Button>{t("common.button.edit")}</Button>}
                    onFinish={handleBaseInfoFormFinish}
                  />,

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
              <div className="py-12 lg:pr-36 xl:pr-48">
                <WorkflowElements />
              </div>
              <div className="absolute right-0 top-0 z-[1]">
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
          </Show>

          <Show when={tabValue === "runs"}>
            <WorkflowRuns workflowId={workflowId!} />
          </Show>
        </Card>
      </div>
    </div>
  );
};

const WorkflowBaseInfoModalForm = ({
  data,
  trigger,
  onFinish,
}: {
  data: Pick<WorkflowModel, "name" | "description">;
  trigger?: React.ReactNode;
  onFinish?: (values: Pick<WorkflowModel, "name" | "description">) => Promise<void | boolean>;
}) => {
  const { t } = useTranslation();

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
    ...formApi
  } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: data,
    onSubmit: async () => {
      const ret = await onFinish?.(formInst.getFieldsValue(true));
      if (ret != null && !ret) return false;
      return true;
    },
  });

  const handleFormFinish = async () => {
    return formApi.submit();
  };

  return (
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
  );
};

export default WorkflowDetail;
