import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Button, Drawer, Form, Space, type DrawerProps, type FormProps, type ModalProps } from "antd";

import { useAntdForm, useTriggerElement } from "@/hooks";

export interface DrawerFormProps<T extends NonNullable<unknown> = NonNullable<unknown>> extends Omit<FormProps<T>, "title" | "onFinish"> {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  cancelButtonProps?: ModalProps["cancelButtonProps"];
  cancelText?: ModalProps["cancelText"];
  defaultOpen?: boolean;
  drawerProps?: Omit<DrawerProps, "open" | "title" | "width">;
  okButtonProps?: ModalProps["okButtonProps"];
  okText?: ModalProps["okText"];
  open?: boolean;
  title?: React.ReactNode;
  trigger?: React.ReactNode;
  width?: string | number;
  onOpenChange?: (open: boolean) => void;
  onFinish?: (values: T) => void | Promise<unknown>;
}

const DrawerForm = <T extends NonNullable<unknown> = NonNullable<unknown>>({
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

  const triggerDom = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  const {
    form: formInst,
    formPending,
    formProps,
    submit,
  } = useAntdForm({
    form,
    onSubmit: async (values) => {
      const ret = await onFinish?.(values);
      if (ret != null && !ret) return false;
      return true;
    },
  });
  const mergedFormProps = { ...formProps, ...props };

  const handleOkClick = async () => {
    const ret = await submit();
    if (ret != null && !ret) return;

    setOpen(false);
  };

  const handleCancelClick = () => {
    if (formPending) return;

    setOpen(false);
  };

  return (
    <>
      {triggerDom}

      <Drawer
        footer={
          <Space>
            <Button {...cancelButtonProps} onClick={handleCancelClick}>
              {cancelText || t("common.button.cancel")}
            </Button>
            <Button type="primary" loading={formPending} {...okButtonProps} onClick={handleOkClick}>
              {okText || t("common.button.ok")}
            </Button>
          </Space>
        }
        open={open}
        title={title}
        width={width}
        {...drawerProps}
        onClose={() => setOpen(false)}
      >
        <Form className={className} style={style} form={formInst} {...mergedFormProps}>
          {children}
        </Form>
      </Drawer>
    </>
  );
};

export default DrawerForm;
