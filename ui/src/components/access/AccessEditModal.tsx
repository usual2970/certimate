import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Modal, notification } from "antd";

import { type AccessModel } from "@/domain/access";
import { accessProvidersMap } from "@/domain/provider";
import { useTriggerElement, useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { getErrMsg } from "@/utils/error";

import AccessEditForm, { type AccessEditFormInstance, type AccessEditFormProps } from "./AccessEditForm";

export type AccessEditModalProps = {
  data?: AccessEditFormProps["initialValues"];
  loading?: boolean;
  open?: boolean;
  preset: AccessEditFormProps["preset"];
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
  onSubmit?: (record: AccessModel) => void;
};

const AccessEditModal = ({ data, loading, trigger, preset, onSubmit, ...props }: AccessEditModalProps) => {
  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { createAccess, updateAccess } = useAccessesStore(useZustandShallowSelector(["createAccess", "updateAccess"]));

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  const formRef = useRef<AccessEditFormInstance>(null);
  const [formPending, setFormPending] = useState(false);

  const handleClickOk = async () => {
    setFormPending(true);
    try {
      await formRef.current!.validateFields();
    } catch (err) {
      setFormPending(false);
      return Promise.reject(err);
    }

    try {
      let temp: AccessModel = formRef.current!.getFieldsValue();
      temp.usage = accessProvidersMap.get(temp.provider)!.usage;
      if (preset === "add") {
        if (data?.id) {
          throw "Invalid props: `data`";
        }

        temp = await createAccess(temp);
      } else if (preset === "edit") {
        if (!data?.id) {
          throw "Invalid props: `data`";
        }

        temp = await updateAccess({ ...data, ...temp });
      } else {
        throw "Invalid props: `preset`";
      }

      onSubmit?.(temp);
      setOpen(false);
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
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
        <div className="pb-2 pt-4">
          <AccessEditForm ref={formRef} initialValues={data} preset={preset === "add" ? "add" : "edit"} />
        </div>
      </Modal>
    </>
  );
};

export default AccessEditModal;
