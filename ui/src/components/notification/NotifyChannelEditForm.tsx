import { forwardRef, useImperativeHandle, useMemo, useState } from "react";
import { useCreation, useDeepCompareEffect } from "ahooks";
import { Form } from "antd";

import { type NotifyChannelsSettingsContent } from "@/domain/settings";
import NotifyChannelEditFormBarkFields from "./NotifyChannelEditFormBarkFields";
import NotifyChannelEditFormDingTalkFields from "./NotifyChannelEditFormDingTalkFields";
import NotifyChannelEditFormEmailFields from "./NotifyChannelEditFormEmailFields";
import NotifyChannelEditFormLarkFields from "./NotifyChannelEditFormLarkFields";
import NotifyChannelEditFormServerChanFields from "./NotifyChannelEditFormServerChanFields";
import NotifyChannelEditFormTelegramFields from "./NotifyChannelEditFormTelegramFields";
import NotifyChannelEditFormWebhookFields from "./NotifyChannelEditFormWebhookFields";

type NotifyChannelEditFormModelType = NotifyChannelsSettingsContent[keyof NotifyChannelsSettingsContent];

export type NotifyChannelEditFormProps = {
  className?: string;
  style?: React.CSSProperties;
  channel: keyof NotifyChannelsSettingsContent;
  disabled?: boolean;
  loading?: boolean;
  model?: NotifyChannelEditFormModelType;
  onModelChange?: (model: NotifyChannelEditFormModelType) => void;
};

export type NotifyChannelEditFormInstance = {
  getFieldsValue: () => NotifyChannelEditFormModelType;
  resetFields: () => void;
  validateFields: () => Promise<NotifyChannelEditFormModelType>;
};

const NotifyChannelEditForm = forwardRef<NotifyChannelEditFormInstance, NotifyChannelEditFormProps>(
  ({ className, style, channel, disabled, loading, model, onModelChange }, ref) => {
    const [form] = Form.useForm();
    const formName = useCreation(() => `notifyChannelEditForm_${Math.random().toString(36).substring(2, 10)}${new Date().getTime()}`, []);
    const formFieldsComponent = useMemo(() => {
      /*
        注意：如果追加新的子组件，请保持以 ASCII 排序。
        NOTICE: If you add new child component, please keep ASCII order.
       */
      switch (channel) {
        case "bark":
          return <NotifyChannelEditFormBarkFields />;
        case "dingtalk":
          return <NotifyChannelEditFormDingTalkFields />;
        case "email":
          return <NotifyChannelEditFormEmailFields />;
        case "lark":
          return <NotifyChannelEditFormLarkFields />;
        case "serverchan":
          return <NotifyChannelEditFormServerChanFields />;
        case "telegram":
          return <NotifyChannelEditFormTelegramFields />;
        case "webhook":
          return <NotifyChannelEditFormWebhookFields />;
      }
    }, [channel]);

    const [initialValues, setInitialValues] = useState(model);
    useDeepCompareEffect(() => {
      setInitialValues(model);
    }, [model]);

    const handleFormChange = (_: unknown, fields: NotifyChannelEditFormModelType) => {
      onModelChange?.(fields);
    };

    useImperativeHandle(ref, () => ({
      getFieldsValue: () => {
        return form.getFieldsValue(true);
      },
      resetFields: () => {
        return form.resetFields();
      },
      validateFields: () => {
        return form.validateFields();
      },
    }));

    return (
      <Form
        className={className}
        style={style}
        form={form}
        disabled={loading || disabled}
        initialValues={initialValues}
        layout="vertical"
        name={formName}
        onValuesChange={handleFormChange}
      >
        {formFieldsComponent}
      </Form>
    );
  }
);

export default NotifyChannelEditForm;
