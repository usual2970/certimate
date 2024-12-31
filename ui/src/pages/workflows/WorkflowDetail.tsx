import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Card, Dropdown, Form, Input, message, Modal, notification, Tabs, Typography } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { PageHeader } from "@ant-design/pro-components";
import { DeleteOutlined as DeleteOutlinedIcon, EllipsisOutlined as EllipsisOutlinedIcon } from "@ant-design/icons";
import { z } from "zod";

import Show from "@/components/Show";
import ModalForm from "@/components/core/ModalForm";
import End from "@/components/workflow/End";
import NodeRender from "@/components/workflow/NodeRender";
import WorkflowRuns from "@/components/workflow/run/WorkflowRuns";
import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { allNodesValidated, type WorkflowModel, type WorkflowNode } from "@/domain/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { remove as removeWorkflow } from "@/repository/workflow";
import { getErrMsg } from "@/utils/error";

const WorkflowDetail = () => {
  const navigate = useNavigate();

  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [modalApi, ModalContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { id: workflowId } = useParams();
  const { workflow, init, setBaseInfo, switchEnable } = useWorkflowStore(useZustandShallowSelector(["workflow", "init", "setBaseInfo", "switchEnable"]));
  useEffect(() => {
    init(workflowId);
  }, [workflowId, init]);

  const [tabValue, setTabValue] = useState<"orchestration" | "runs">("orchestration");

  // const [running, setRunning] = useState(false);

  const elements = useMemo(() => {
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
    if (!workflow.enabled && !allNodesValidated(workflow.content!)) {
      messageApi.warning(t("workflow.detail.action.save.failed.uncompleted"));
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
            navigate("/workflows");
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  // TODO: 发布更改 撤销更改 立即执行
  // const handleWorkflowSaveClick = () => {
  //   if (!allNodesValidated(workflow.draft as WorkflowNode)) {
  //     messageApi.warning(t("workflow.detail.action.save.failed.uncompleted"));
  //     return;
  //   }
  //   save();
  // };

  // const handleRunClick = async () => {
  //   if (running) return;

  //   setRunning(true);
  //   try {
  //     await runWorkflow(workflow.id as string);
  //     messageApi.success(t("workflow.detail.action.run.success"));
  //   } catch (err) {
  //     console.error(err);
  //     messageApi.warning(t("workflow.detail.action.run.failed"));
  //   } finally {
  //     setRunning(false);
  //   }
  // };

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
          <WorkflowProvider>
            <Show when={tabValue === "orchestration"}>
              <div className="flex flex-col items-center">{elements}</div>
            </Show>

            <Show when={tabValue === "runs"}>
              <WorkflowRuns />
            </Show>
          </WorkflowProvider>
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
      .string({ message: t("workflow.baseinfo.form.name.placeholder") })
      .min(1, t("workflow.baseinfo.form.name.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    description: z
      .string({ message: t("workflow.baseinfo.form.description.placeholder") })
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

  const handleFinish = async () => {
    return formApi.submit();
  };

  return (
    <ModalForm
      disabled={formPending}
      layout="vertical"
      form={formInst}
      modalProps={{ destroyOnClose: true }}
      okText={t("common.button.save")}
      title={t(`workflow.baseinfo.modal.title`)}
      trigger={trigger}
      width={480}
      {...formProps}
      onFinish={handleFinish}
    >
      <Form.Item name="name" label={t("workflow.baseinfo.form.name.label")} rules={[formRule]}>
        <Input placeholder={t("workflow.baseinfo.form.name.placeholder")} />
      </Form.Item>

      <Form.Item name="description" label={t("workflow.baseinfo.form.description.label")} rules={[formRule]}>
        <Input placeholder={t("workflow.baseinfo.form.description.placeholder")} />
      </Form.Item>
    </ModalForm>
  );
};
export default WorkflowDetail;
