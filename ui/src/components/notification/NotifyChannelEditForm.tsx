import { forwardRef, useImperativeHandle, useMemo } from "react";
import { Form, type FormInstance } from "antd";

import { NOTIFY_CHANNELS, type NotifyChannelsSettingsContent } from "@/domain/settings";
import { useAntdForm } from "@/hooks";

import NotifyChannelEditFormBarkFields from "./NotifyChannelEditFormBarkFields";
import NotifyChannelEditFormDingTalkFields from "./NotifyChannelEditFormDingTalkFields";
import NotifyChannelEditFormEmailFields from "./NotifyChannelEditFormEmailFields";
import NotifyChannelEditFormGotifyFields from "./NotifyChannelEditFormGotifyFields.tsx";
import NotifyChannelEditFormLarkFields from "./NotifyChannelEditFormLarkFields";
import NotifyChannelEditFormMattermostFields from "./NotifyChannelEditFormMattermostFields.tsx";
import NotifyChannelEditFormPushoverFields from "./NotifyChannelEditFormPushoverFields";
import NotifyChannelEditFormPushPlusFields from "./NotifyChannelEditFormPushPlusFields";
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

export type NotifyChannelEditFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<NotifyChannelEditFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<NotifyChannelEditFormFieldValues>["resetFields"];
  validateFields: FormInstance<NotifyChannelEditFormFieldValues>["validateFields"];
};

const NotifyChannelEditForm = forwardRef<NotifyChannelEditFormInstance, NotifyChannelEditFormProps>(
  ({ className, style, channel, disabled, initialValues, onValuesChange }, ref) => {
    const { form: formInst, formProps } = useAntdForm({
      initialValues: initialValues,
      name: "notifyChannelEditForm",
    });
    const formFieldsEl = useMemo(() => {
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
        case NOTIFY_CHANNELS.GOTIFY:
          return <NotifyChannelEditFormGotifyFields />;
        case NOTIFY_CHANNELS.LARK:
          return <NotifyChannelEditFormLarkFields />;
        case NOTIFY_CHANNELS.MATTERMOST:
          return <NotifyChannelEditFormMattermostFields />;
        case NOTIFY_CHANNELS.PUSHOVER:
          return <NotifyChannelEditFormPushoverFields />;
        case NOTIFY_CHANNELS.PUSHPLUS:
          return <NotifyChannelEditFormPushPlusFields />;
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
        getFieldsValue: () => {
          return formInst.getFieldsValue(true);
        },
        resetFields: (fields) => {
          return formInst.resetFields(fields);
        },
        validateFields: (nameList, config) => {
          return formInst.validateFields(nameList, config);
        },
      } as NotifyChannelEditFormInstance;
    });

    return (
      <Form
        {...formProps}
        className={className}
        style={style}
        form={formInst}
        disabled={disabled}
        layout="vertical"
        scrollToFirstError
        onValuesChange={handleFormChange}
      >
        {formFieldsEl}
      </Form>
    );
  }
);

export default NotifyChannelEditForm;
