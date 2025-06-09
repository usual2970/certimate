import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Modal, notification } from "antd";

import ModalForm from "@/components/ModalForm";
import { useTriggerElement, useZustandShallowSelector } from "@/hooks";
import { getErrMsg } from "@/utils/error";

import WorkflowForm, { type WorkflowFormInstance, type WorkflowFormProps } from "./WorkflowForm";

export type WorkflowEditModalProps = {
  data?: WorkflowFormProps["initialValues"];
  loading?: boolean;
  open?: boolean;
  usage?: WorkflowFormProps["usage"];
  scene: WorkflowFormProps["scene"];
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
  afterSubmit?: (record: WorkflowModel) => void;
};

const WorkflowEditModal = ({ data, loading, trigger, scene, usage, afterSubmit, ...props }: WorkflowEditModalProps) => {
  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { createWorkflow, updateWorkflow } = useWorkflowesStore(useZustandShallowSelector(["createWorkflow", "updateWorkflow"]));

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  const formRef = useRef<WorkflowFormInstance>(null);
  const [formPending, setFormPending] = useState(false);

  const handleOkClick = async () => {
    setFormPending(true);
    try {
      await formRef.current!.validateFields();
    } catch (err) {
      setFormPending(false);
      throw err;
    }

    try {
      let values: WorkflowModel = formRef.current!.getFieldsValue();

      if (scene === "add") {
        if (data?.id) {
          throw "Invalid props: `data`";
        }

        values = await createWorkflow(values);
      } else if (scene === "edit") {
        if (!data?.id) {
          throw "Invalid props: `data`";
        }

        values = await updateWorkflow({ ...data, ...values });
      } else {
        throw "Invalid props: `preset`";
      }

      afterSubmit?.(values);
      setOpen(false);
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
    } finally {
      setFormPending(false);
    }
  };

  const handleCancelClick = () => {
    if (formPending) return;

    setOpen(false);
  };

  return (
    <>
      {NotificationContextHolder}

      {triggerEl}

      <Modal
        styles={{
          content: {
            maxHeight: "calc(80vh - 64px)",
            overflowX: "hidden",
            overflowY: "auto",
          },
        }}
        afterClose={() => setOpen(false)}
        cancelButtonProps={{ disabled: formPending }}
        cancelText={t("common.button.cancel")}
        closable
        confirmLoading={formPending}
        destroyOnHidden
        loading={loading}
        okText={scene === "edit" ? t("common.button.save") : t("common.button.submit")}
        open={open}
        title={t(`access.action.${scene}`)}
        width={480}
        onOk={handleOkClick}
        onCancel={handleCancelClick}
      >
        <div className="pb-2 pt-4">
          <WorkflowForm ref={formRef} initialValues={data} scene={scene === "add" ? "add" : "edit"} usage={usage} />
        </div>
      </Modal>
    </>
  );
};

export default WorkflowEditModal;
