import { forwardRef, useImperativeHandle, useMemo } from "react";
import { useCreation } from "ahooks";
import { Form } from "antd";

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

type NotifyChannelEditFormModelValues = NotifyChannelsSettingsContent[keyof NotifyChannelsSettingsContent];

export type NotifyChannelEditFormProps = {
  className?: string;
  style?: React.CSSProperties;
  channel: string;
  disabled?: boolean;
  model?: NotifyChannelEditFormModelValues;
  onModelChange?: (model: NotifyChannelEditFormModelValues) => void;
};

export type NotifyChannelEditFormInstance = {
  getFieldsValue: () => NotifyChannelEditFormModelValues;
  resetFields: () => void;
  validateFields: () => Promise<NotifyChannelEditFormModelValues>;
};

const NotifyChannelEditForm = forwardRef<NotifyChannelEditFormInstance, NotifyChannelEditFormProps>(
  ({ className, style, channel, disabled, model, onModelChange }, ref) => {
    const { form: formInst, formProps } = useAntdForm({
      initialValues: model,
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

    const handleFormChange = (_: unknown, values: NotifyChannelEditFormModelValues) => {
      onModelChange?.(values);
    };

    useImperativeHandle(ref, () => ({
      getFieldsValue: () => {
        return formInst.getFieldsValue(true);
      },
      resetFields: () => {
        return formInst.resetFields();
      },
      validateFields: () => {
        return formInst.validateFields();
      },
    }));

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
