import { cloneElement, memo, useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Button, Card, Dropdown, Form, Input, message, Modal, notification, Tabs, Typography, type FormInstance } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { PageHeader } from "@ant-design/pro-components";
import { z } from "zod";
import { Ellipsis as EllipsisIcon, Trash2 as Trash2Icon } from "lucide-react";

import Show from "@/components/Show";
import End from "@/components/workflow/End";
import NodeRender from "@/components/workflow/NodeRender";
import WorkflowRuns from "@/components/workflow/run/WorkflowRuns";
import WorkflowProvider from "@/components/workflow/WorkflowProvider";
import { useZustandShallowSelector } from "@/hooks";
import { allNodesValidated, type WorkflowModel, type WorkflowNode } from "@/domain/workflow";
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
  const { workflow, init, setBaseInfo, switchEnable, save } = useWorkflowStore(
    useZustandShallowSelector(["workflow", "init", "setBaseInfo", "switchEnable", "save"])
  );
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

  const handleBaseInfoFormFinish = async (fields: Pick<WorkflowModel, "name" | "description">) => {
    try {
      await setBaseInfo(fields.name!, fields.description!);
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
              <WorkflowBaseInfoModalForm model={workflow} trigger={<Button>{t("common.button.edit")}</Button>} onFinish={handleBaseInfoFormFinish} />

              <Button onClick={handleEnableChange}>{workflow.enabled ? t("common.button.disable") : t("common.button.enable")}</Button>

              <Dropdown
                menu={{
                  items: [
                    {
                      key: "delete",
                      label: t("common.button.delete"),
                      danger: true,
                      icon: <Trash2Icon size={14} />,
                      onClick: () => {
                        handleDeleteClick();
                      },
                    },
                  ],
                }}
              >
                <Button icon={<EllipsisIcon size={14} />} />
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

const WorkflowBaseInfoModalForm = memo(
  ({
    model,
    trigger,
    onFinish,
  }: {
    model: Pick<WorkflowModel, "name" | "description">;
    trigger?: React.ReactElement;
    onFinish?: (fields: Pick<WorkflowModel, "name" | "description">) => Promise<void | boolean>;
  }) => {
    const { t } = useTranslation();

    const [open, setOpen] = useState(false);

    const triggerEl = useMemo(() => {
      if (!trigger) {
        return null;
      }

      return cloneElement(trigger, {
        ...trigger.props,
        onClick: () => {
          setOpen(true);
          trigger.props?.onClick?.();
        },
      });
    }, [trigger, setOpen]);

    const formSchema = z.object({
      name: z
        .string({ message: t("workflow.baseinfo.form.name.placeholder") })
        .trim()
        .min(1, t("workflow.baseinfo.form.name.placeholder"))
        .max(64, t("common.errmsg.string_max", { max: 64 })),
      description: z
        .string({ message: t("workflow.baseinfo.form.description.placeholder") })
        .trim()
        .min(0, t("workflow.baseinfo.form.description.placeholder"))
        .max(256, t("common.errmsg.string_max", { max: 256 }))
        .nullish(),
    });
    const formRule = createSchemaFieldRule(formSchema);
    const [form] = Form.useForm<FormInstance<z.infer<typeof formSchema>>>();
    const [formPending, setFormPending] = useState(false);

    const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model as Partial<z.infer<typeof formSchema>>);
    useDeepCompareEffect(() => {
      setInitialValues(model as Partial<z.infer<typeof formSchema>>);
    }, [model]);

    const handleClickOk = async () => {
      setFormPending(true);
      try {
        await form.validateFields();
      } catch (err) {
        setFormPending(false);
        return Promise.reject();
      }

      try {
        const ret = await onFinish?.(form.getFieldsValue(true));
        if (ret != null && !ret) return;

        setOpen(false);
      } finally {
        setFormPending(false);
      }
    };

    const handleClickCancel = () => {
      if (formPending) return Promise.reject();

      setOpen(false);
    };

    return (
      <>
        {triggerEl}

        <Modal
          afterClose={() => setOpen(false)}
          cancelButtonProps={{ disabled: formPending }}
          closable
          confirmLoading={formPending}
          destroyOnClose
          okText={t("common.button.save")}
          open={open}
          title={t(`workflow.baseinfo.modal.title`)}
          width={480}
          onOk={handleClickOk}
          onCancel={handleClickCancel}
        >
          <div className="pt-4 pb-2">
            <Form form={form} initialValues={initialValues} layout="vertical">
              <Form.Item name="name" label={t("workflow.baseinfo.form.name.label")} rules={[formRule]}>
                <Input placeholder={t("workflow.baseinfo.form.name.placeholder")} />
              </Form.Item>

              <Form.Item name="description" label={t("workflow.baseinfo.form.description.label")} rules={[formRule]}>
                <Input placeholder={t("workflow.baseinfo.form.description.placeholder")} />
              </Form.Item>
            </Form>
          </div>
        </Modal>
      </>
    );
  }
);

export default WorkflowDetail;
