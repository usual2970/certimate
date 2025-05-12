import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Button, Drawer, Space, notification } from "antd";

import { type AccessModel } from "@/domain/access";
import { useTriggerElement, useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";
import { getErrMsg } from "@/utils/error";

import AccessForm, { type AccessFormInstance, type AccessFormProps } from "./AccessForm";

export type AccessEditDrawerProps = {
  data?: AccessFormProps["initialValues"];
  loading?: boolean;
  open?: boolean;
  scene: AccessFormProps["scene"];
  trigger?: React.ReactNode;
  usage?: AccessFormProps["usage"];
  onOpenChange?: (open: boolean) => void;
  afterSubmit?: (record: AccessModel) => void;
};

const AccessEditDrawer = ({ data, loading, trigger, scene, usage, afterSubmit, ...props }: AccessEditDrawerProps) => {
  const { t } = useTranslation();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { createAccess, updateAccess } = useAccessesStore(useZustandShallowSelector(["createAccess", "updateAccess"]));

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  const formRef = useRef<AccessFormInstance>(null);
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
      let values: AccessModel = formRef.current!.getFieldsValue();

      if (scene === "add") {
        if (data?.id) {
          throw "Invalid props: `data`";
        }

        values = await createAccess(values);
      } else if (scene === "edit") {
        if (!data?.id) {
          throw "Invalid props: `data`";
        }

        values = await updateAccess({ ...data, ...values });
      } else {
        throw "Invalid props: `scene`";
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

      <Drawer
        afterOpenChange={setOpen}
        closable={!formPending}
        destroyOnHidden
        footer={
          <Space className="w-full justify-end">
            <Button onClick={handleCancelClick}>{t("common.button.cancel")}</Button>
            <Button loading={formPending} type="primary" onClick={handleOkClick}>
              {scene === "edit" ? t("common.button.save") : t("common.button.submit")}
            </Button>
          </Space>
        }
        loading={loading}
        maskClosable={!formPending}
        open={open}
        title={t(`access.action.${scene}`)}
        width={720}
        onClose={() => setOpen(false)}
      >
        <AccessForm ref={formRef} initialValues={data} scene={scene === "add" ? "add" : "edit"} usage={usage} />
      </Drawer>
    </>
  );
};

export default AccessEditDrawer;
