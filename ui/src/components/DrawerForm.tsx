import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Button, Drawer, type DrawerProps, Form, type FormProps, type ModalProps, Space } from "antd";

import { useAntdForm, useTriggerElement } from "@/hooks";

export interface DrawerFormProps<T extends NonNullable<unknown> = any> extends Omit<FormProps<T>, "title" | "onFinish"> {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  cancelButtonProps?: ModalProps["cancelButtonProps"];
  cancelText?: ModalProps["cancelText"];
  defaultOpen?: boolean;
  drawerProps?: Omit<DrawerProps, "defaultOpen" | "forceRender" | "open" | "title" | "width" | "onOpenChange">;
  okButtonProps?: ModalProps["okButtonProps"];
  okText?: ModalProps["okText"];
  open?: boolean;
  title?: React.ReactNode;
  trigger?: React.ReactNode;
  width?: string | number;
  onClose?: (e: React.MouseEvent | React.KeyboardEvent) => void | Promise<unknown>;
  onFinish?: (values: T) => unknown | Promise<unknown>;
  onOpenChange?: (open: boolean) => void;
}

const DrawerForm = <T extends NonNullable<unknown> = any>({
  className,
  style,
  children,
  cancelText,
  cancelButtonProps,
  form,
  drawerProps,
  okText,
  okButtonProps,
  title,
  trigger,
  width,
  onFinish,
  ...props
}: DrawerFormProps<T>) => {
  const { t } = useTranslation();

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerEl = useTriggerElement(trigger, {
    onClick: () => {
      setOpen(true);
    },
  });

  const {
    form: formInst,
    formPending,
    formProps,
    submit,
  } = useAntdForm({
    form,
    onSubmit: (values) => {
      return onFinish?.(values);
    },
  });

  const mergedFormProps: FormProps = {
    clearOnDestroy: drawerProps?.destroyOnHidden ? true : undefined,
    ...formProps,
    ...props,
  };

  const mergedDrawerProps: DrawerProps = {
    ...drawerProps,
    afterOpenChange: (open) => {
      if (!open && !mergedFormProps.preserve) {
        formInst.resetFields();
      }

      drawerProps?.afterOpenChange?.(open);
    },
    onClose: async (e) => {
      if (formPending) return;

      // 关闭 Drawer 时 Promise.reject 阻止关闭
      await drawerProps?.onClose?.(e);
      setOpen(false);
    },
  };

  const handleOkClick = async () => {
    // 提交表单返回 Promise.reject 时不关闭 Drawer
    await submit();

    setOpen(false);
  };

  const handleCancelClick = () => {
    if (formPending) return;

    setOpen(false);
  };

  return (
    <>
      {triggerEl}

      <Drawer
        {...mergedDrawerProps}
        footer={
          <Space className="w-full justify-end">
            <Button {...cancelButtonProps} onClick={handleCancelClick}>
              {cancelText ?? t("common.button.cancel")}
            </Button>
            <Button {...okButtonProps} type="primary" loading={formPending} onClick={handleOkClick}>
              {okText ?? t("common.button.ok")}
            </Button>
          </Space>
        }
        forceRender
        open={open}
        title={title}
        width={width}
      >
        <Form className={className} style={style} {...mergedFormProps} form={formInst}>
          {children}
        </Form>
      </Drawer>
    </>
  );
};

export default DrawerForm;
