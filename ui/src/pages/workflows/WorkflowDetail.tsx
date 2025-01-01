import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Button, Card, Dropdown, Form, Input, message, Modal, notification, Space, Tabs, Typography } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { PageHeader } from "@ant-design/pro-components";
import {
  CaretRightOutlined as CaretRightOutlinedIcon,
  DeleteOutlined as DeleteOutlinedIcon,
  EllipsisOutlined as EllipsisOutlinedIcon,
  UndoOutlined as UndoOutlinedIcon,
} from "@ant-design/icons";
import { ClientResponseError } from "pocketbase";
import { isEqual } from "radash";
import { z } from "zod";

import Show from "@/components/Show";
import ModalForm from "@/components/core/ModalForm";
import End from "@/components/workflow/End";
import NodeRender from "@/components/workflow/NodeRender";
import WorkflowRuns from "@/components/workflow/run/WorkflowRuns";
import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { isAllNodesValidated, type WorkflowModel, type WorkflowNode } from "@/domain/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { remove as removeWorkflow } from "@/repository/workflow";
import { run as runWorkflow } from "@/api/workflow";
import { getErrMsg } from "@/utils/error";

const WorkflowDetail = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [modalApi, ModalContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { id: workflowId } = useParams();
  const { workflow, init, save, setBaseInfo, switchEnable } = useWorkflowStore(
    useZustandShallowSelector(["workflow", "init", "save", "setBaseInfo", "switchEnable"])
  );
  useEffect(() => {
    // TODO: loading
    init(workflowId);
  }, [workflowId, init]);

  const [tabValue, setTabValue] = useState<"orchestration" | "runs">("orchestration");

  const workflowNodes = useMemo(() => {
    let current = workflow.draft as WorkflowNode;

    const elements: JSX.Element[] = [];

    while (current) {
      // 处理普通节点
      elements.push(<NodeRender data={current} key={current.id} />);
      current = current.next as WorkflowNode;
    }

    elements.push(<End key="workflow-end" />);

    return elements;
  }, [workflow]);

  const [workflowRunning, setWorkflowRunning] = useState(false);

  const [allowDiscard, setAllowDiscard] = useState(false);
  const [allowRelease, setAllowRelease] = useState(false);
  useDeepCompareEffect(() => {
    const hasChanges = workflow.hasDraft! || !isEqual(workflow.draft, workflow.content);
    setAllowDiscard(hasChanges && !workflowRunning);
    setAllowRelease(hasChanges && !workflowRunning);
  }, [workflow, workflowRunning]);

  const handleBaseInfoFormFinish = async (values: Pick<WorkflowModel, "name" | "description">) => {
    try {
      await setBaseInfo(values.name!, values.description!);
    } catch (err) {
      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
      return false;
    }
  };

  const handleEnableChange = () => {
    if (!workflow.enabled && !isAllNodesValidated(workflow.content!)) {
      messageApi.warning(t("workflow.action.enable.failed.uncompleted"));
      return;
    }

    switchEnable();
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
    alert("TODO");
  };

  const handleReleaseClick = () => {
    if (!isAllNodesValidated(workflow.draft!)) {
      messageApi.warning(t("workflow.action.release.failed.uncompleted"));
      return;
    }

    save();

    messageApi.success(t("common.text.operation_succeeded"));
  };

  const handleRunClick = () => {
    if (!workflow.enabled) {
      alert("TODO: 暂时只支持执行已启用的工作流");
      return;
    }

    const { promise, resolve, reject } = Promise.withResolvers();
    if (workflow.hasDraft) {
      modalApi.confirm({
        title: t("workflow.action.run"),
        content: t("workflow.action.run.confirm"),
        onOk: () => resolve(void 0),
        onCancel: () => reject(),
      });
    } else {
      resolve(void 0);
    }

    // TODO: 异步执行
    promise.then(async () => {
      setWorkflowRunning(true);

      try {
        await runWorkflow(workflowId!);

        messageApi.warning(t("common.text.operation_succeeded"));
      } catch (err) {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);
        messageApi.warning(t("common.text.operation_failed"));
      } finally {
        setWorkflowRunning(false);
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
          extra={[
            <Button.Group key="actions">
              <WorkflowBaseInfoModalForm data={workflow} trigger={<Button>{t("common.button.edit")}</Button>} onFinish={handleBaseInfoFormFinish} />

              <Button onClick={handleEnableChange}>{workflow.enabled ? t("common.button.disable") : t("common.button.enable")}</Button>

              <Dropdown
                menu={{
                  items: [
                    {
                      key: "delete",
                      label: t("common.button.delete"),
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
                <Button icon={<EllipsisOutlinedIcon />} />
              </Dropdown>
            </Button.Group>,
          ]}
        >
          <Typography.Paragraph type="secondary">{workflow.description}</Typography.Paragraph>
          <Tabs
            activeKey={tabValue}
            defaultActiveKey="orchestration"
            items={[
              { key: "orchestration", label: t("workflow.detail.orchestration.tab") },
              { key: "runs", label: t("workflow.detail.runs.tab") },
            ]}
            renderTabBar={(props, DefaultTabBar) => <DefaultTabBar {...props} style={{ margin: 0 }} />}
            tabBarStyle={{ border: "none" }}
            onChange={(key) => setTabValue(key as typeof tabValue)}
          />
        </PageHeader>
      </Card>

      <div className="p-4">
        <Card>
          <Show when={tabValue === "orchestration"}>
            <div className="relative">
              <div className="flex flex-col items-center py-12 pr-48">
                <WorkflowProvider>{workflowNodes}</WorkflowProvider>
              </div>
              <div className="absolute top-0 right-0 z-[1]">
                <Space>
                  <Button icon={<CaretRightOutlinedIcon />} loading={workflowRunning} type="primary" onClick={handleRunClick}>
                    {t("workflow.action.run")}
                  </Button>

                  <Button.Group>
                    <Button color="primary" disabled={!allowRelease} variant="outlined" onClick={handleReleaseClick}>
                      {t("workflow.action.release")}
                    </Button>

                    <Dropdown
                      menu={{
                        items: [
                          {
                            key: "discard",
                            disabled: !allowDiscard,
                            label: t("workflow.action.discard"),
                            icon: <UndoOutlinedIcon />,
                            onClick: handleDiscardClick,
                          },
                        ],
                      }}
                      trigger={["click"]}
                    >
                      <Button color="primary" icon={<EllipsisOutlinedIcon />} variant="outlined" />
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
