import { forwardRef, useImperativeHandle, useMemo } from "react";
import { useCreation } from "ahooks";
import { Form, type FormInstance } from "antd";

import { useAntdForm } from "@/hooks";
import { NOTIFY_CHANNELS, type NotifyChannelsSettingsContent } from "@/domain/settings";
import NotifyChannelEditFormBarkFields from "./NotifyChannelEditFormBarkFields";
import NotifyChannelEditFormDingTalkFields from "./NotifyChannelEditFormDingTalkFields";
import NotifyChannelEditFormEmailFields from "./NotifyChannelEditFormEmailFields";
import NotifyChannelEditFormLarkFields from "./NotifyChannelEditFormLarkFields";
import NotifyChannelEditFormServerChanFields from "./NotifyChannelEditFormServerChanFields";
import NotifyChannelEditFormTelegramFields from "./NotifyChannelEditFormTelegramFields";
import NotifyChannelEditFormWebhookFields from "./NotifyChannelEditFormWebhookFields";
import NotifyChannelEditFormWeComFields from "./NotifyChannelEditFormWeComFields";

type NotifyChannelEditFormFieldValues = NotifyChannelsSettingsContent[keyof NotifyChannelsSettingsContent];

export type NotifyChannelEditFormProps = {
  className?: string;
  style?: React.CSSProperties;
  channel: string;
  disabled?: boolean;
  initialValues?: NotifyChannelEditFormFieldValues;
  onValuesChange?: (values: NotifyChannelEditFormFieldValues) => void;
};

export type NotifyChannelEditFormInstance = FormInstance;

const NotifyChannelEditForm = forwardRef<NotifyChannelEditFormInstance, NotifyChannelEditFormProps>(
  ({ className, style, channel, disabled, initialValues, onValuesChange }, ref) => {
    const { form: formInst, formProps } = useAntdForm({
      initialValues: initialValues,
    });
    const formName = useCreation(() => `notifyChannelEditForm_${Math.random().toString(36).substring(2, 10)}${new Date().getTime()}`, []);
    const formFieldsComponent = useMemo(() => {
      /*
        注意：如果追加新的子组件，请保持以 ASCII 排序。
        NOTICE: If you add new child component, please keep ASCII order.
       */
      switch (channel) {
        case NOTIFY_CHANNELS.BARK:
          return <NotifyChannelEditFormBarkFields />;
        case NOTIFY_CHANNELS.DINGTALK:
          return <NotifyChannelEditFormDingTalkFields />;
        case NOTIFY_CHANNELS.EMAIL:
          return <NotifyChannelEditFormEmailFields />;
        case NOTIFY_CHANNELS.LARK:
          return <NotifyChannelEditFormLarkFields />;
        case NOTIFY_CHANNELS.SERVERCHAN:
          return <NotifyChannelEditFormServerChanFields />;
        case NOTIFY_CHANNELS.TELEGRAM:
          return <NotifyChannelEditFormTelegramFields />;
        case NOTIFY_CHANNELS.WEBHOOK:
          return <NotifyChannelEditFormWebhookFields />;
        case NOTIFY_CHANNELS.WECOM:
          return <NotifyChannelEditFormWeComFields />;
      }
    }, [channel]);

    const handleFormChange = (_: unknown, values: NotifyChannelEditFormFieldValues) => {
      onValuesChange?.(values);
    };

    useImperativeHandle(ref, () => {
      return {
        getFieldValue: (name) => formInst.getFieldValue(name),
        getFieldsValue: (...args) => {
          if (Array.from(args).length === 0) {
            return formInst.getFieldsValue(true);
          }

          return formInst.getFieldsValue(args[0] as any, args[1] as any);
        },
        getFieldError: (name) => formInst.getFieldError(name),
        getFieldsError: (nameList) => formInst.getFieldsError(nameList),
        getFieldWarning: (name) => formInst.getFieldWarning(name),
        isFieldsTouched: (nameList, allFieldsTouched) => formInst.isFieldsTouched(nameList, allFieldsTouched),
        isFieldTouched: (name) => formInst.isFieldTouched(name),
        isFieldValidating: (name) => formInst.isFieldValidating(name),
        isFieldsValidating: (nameList) => formInst.isFieldsValidating(nameList),
        resetFields: (fields) => formInst.resetFields(fields),
        setFields: (fields) => formInst.setFields(fields),
        setFieldValue: (name, value) => formInst.setFieldValue(name, value),
        setFieldsValue: (values) => formInst.setFieldsValue(values),
        validateFields: (nameList, config) => formInst.validateFields(nameList, config),
        submit: () => formInst.submit(),
      } as NotifyChannelEditFormInstance;
    });

    return (
      <Form
        {...formProps}
        className={className}
        style={style}
        form={formInst}
        disabled={disabled}
        name={formName}
        layout="vertical"
        onValuesChange={handleFormChange}
      >
        {formFieldsComponent}
      </Form>
    );
  }
);

export default NotifyChannelEditForm;
