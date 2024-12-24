import { cloneElement, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Modal, notification } from "antd";

import { type AccessModel } from "@/domain/access";
import { useAccessStore } from "@/stores/access";
import { getErrMsg } from "@/utils/error";
import AccessEditForm, { type AccessEditFormInstance, type AccessEditFormProps } from "./AccessEditForm";

export type AccessEditModalProps = {
  data?: AccessEditFormProps["model"];
  loading?: boolean;
  open?: boolean;
  preset: AccessEditFormProps["preset"];
  trigger?: React.ReactElement;
  onOpenChange?: (open: boolean) => void;
};

const AccessEditModal = ({ data, loading, trigger, preset, ...props }: AccessEditModalProps) => {
  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { createAccess, updateAccess } = useAccessStore();

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

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

  const formRef = useRef<AccessEditFormInstance>(null);
  const [formPending, setFormPending] = useState(false);

  const handleClickOk = async () => {
    setFormPending(true);
    try {
      await formRef.current!.validateFields();
    } catch (err) {
      setFormPending(false);
      return Promise.reject();
    }

    try {
      if (preset === "add") {
        if (data?.id) {
          throw "Invalid props: `data`";
        }

        await createAccess(formRef.current!.getFieldsValue() as AccessModel);
      } else if (preset === "edit") {
        if (!data?.id) {
          throw "Invalid props: `data`";
        }

        await updateAccess({ ...data, ...formRef.current!.getFieldsValue() } as AccessModel);
      }

      setOpen(false);
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

      throw err;
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
      {NotificationContextHolder}

      {triggerEl}

      <Modal
        afterClose={() => setOpen(false)}
        cancelButtonProps={{ disabled: formPending }}
        closable
        confirmLoading={formPending}
        destroyOnClose
        loading={loading}
        okText={preset === "edit" ? t("common.button.save") : t("common.button.submit")}
        open={open}
        title={t(`access.action.${preset}`)}
        width={480}
        onOk={handleClickOk}
        onCancel={handleClickCancel}
      >
        <div className="pt-4 pb-2">
          <AccessEditForm ref={formRef} preset={preset === "add" ? "add" : "edit"} model={data} />
        </div>
      </Modal>
    </>
  );
};

export default AccessEditModal;
