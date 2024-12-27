import { useControllableValue } from "ahooks";
import { Form, Modal, type FormProps, type ModalProps } from "antd";

import { useAntdForm, useTriggerElement } from "@/hooks";

export interface ModalFormProps<T extends NonNullable<unknown> = any> extends Omit<FormProps<T>, "title" | "onFinish"> {
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
  cancelButtonProps?: ModalProps["cancelButtonProps"];
  cancelText?: ModalProps["cancelText"];
  defaultOpen?: boolean;
  modalProps?: Omit<
    ModalProps,
    | "cancelButtonProps"
    | "cancelText"
    | "confirmLoading"
    | "forceRender"
    | "okButtonProps"
    | "okText"
    | "okType"
    | "open"
    | "title"
    | "width"
    | "onOk"
    | "onCancel"
  >;
  okButtonProps?: ModalProps["okButtonProps"];
  okText?: ModalProps["okText"];
  open?: boolean;
  title?: ModalProps["title"];
  trigger?: React.ReactNode;
  width?: ModalProps["width"];
  onOpenChange?: (open: boolean) => void;
  onFinish?: (values: T) => void | Promise<unknown>;
}

const ModalForm = <T extends NonNullable<unknown> = any>({
  className,
  style,
  children,
  cancelButtonProps,
  cancelText,
  form,
  modalProps,
  okButtonProps,
  okText,
  title,
  trigger,
  width,
  onFinish,
  ...props
}: ModalFormProps<T>) => {
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

      <Modal
        cancelButtonProps={cancelButtonProps}
        cancelText={cancelText}
        confirmLoading={formPending}
        forceRender={true}
        okButtonProps={okButtonProps}
        okText={okText}
        okType="primary"
        open={open}
        title={title}
        width={width}
        {...modalProps}
        onOk={handleOkClick}
        onCancel={handleCancelClick}
      >
        <div className="pt-4 pb-2">
          <Form className={className} style={style} form={formInst} {...mergedFormProps}>
            {children}
          </Form>
        </div>
      </Modal>
    </>
  );
};

export default ModalForm;
