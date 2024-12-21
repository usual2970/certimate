import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareMemo } from "@ant-design/pro-components";
import { Button, Collapse, message, notification, Skeleton, Space, Switch, type CollapseProps } from "antd";

import NotifyChannelEditForm, { type NotifyChannelEditFormInstance } from "./NotifyChannelEditForm";
import NotifyTestButton from "./NotifyTestButton";
import { notifyChannelsMap } from "@/domain/settings";
import { useNotifyChannelStore } from "@/stores/notify";
import { getErrMsg } from "@/utils/error";

type NotifyChannelProps = {
  className?: string;
  style?: React.CSSProperties;
  channel: string;
};

const NotifyChannel = ({ className, style, channel }: NotifyChannelProps) => {
  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { channels, setChannel } = useNotifyChannelStore();

  const channelConfig = useDeepCompareMemo(() => channels[channel], [channels, channel]);
  const [channelFormChanged, setChannelFormChanged] = useState(false);
  const channelFormRef = useRef<NotifyChannelEditFormInstance>(null);

  const handleClickSubmit = async () => {
    await channelFormRef.current!.validateFields();

    try {
      setChannel(channel, channelFormRef.current!.getFieldsValue());
      setChannelFormChanged(false);

      messageApi.success(t("common.text.operation_succeeded"));
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    }
  };

  return (
    <div className={className} style={style}>
      {MessageContextHolder}
      {NotificationContextHolder}

      <NotifyChannelEditForm ref={channelFormRef} className="mt-2" channel={channel} model={channelConfig} onModelChange={() => setChannelFormChanged(true)} />

      <Space className="mb-2">
        <Button type="primary" disabled={!channelFormChanged} onClick={handleClickSubmit}>
          {t("common.button.save")}
        </Button>

        {channelConfig != null ? <NotifyTestButton channel={channel} disabled={channelFormChanged} /> : null}
      </Space>
    </div>
  );
};

type NotifyChannelsSemanticDOM = "collapse" | "form";

export type NotifyChannelsProps = {
  className?: string;
  classNames?: Partial<Record<NotifyChannelsSemanticDOM, string>>;
  style?: React.CSSProperties;
  styles?: Partial<Record<NotifyChannelsSemanticDOM, React.CSSProperties>>;
};

const NotifyChannels = ({ className, classNames, style, styles }: NotifyChannelsProps) => {
  const { t, i18n } = useTranslation();

  const { initialized, channels, setChannel, fetchChannels } = useNotifyChannelStore();
  useEffect(() => {
    fetchChannels();
  }, [fetchChannels]);

  const channelCollapseItems: CollapseProps["items"] = useDeepCompareMemo(
    () =>
      Array.from(notifyChannelsMap.values()).map((channel) => {
        return {
          key: `channel-${channel.type}`,
          label: <>{t(channel.name)}</>,
          children: <NotifyChannel className={classNames?.form} style={styles?.form} channel={channel.type} />,
          extra: (
            <div onClick={(e) => e.stopPropagation()} onMouseDown={(e) => e.stopPropagation()} onMouseUp={(e) => e.stopPropagation()}>
              <Switch
                defaultChecked={channels[channel.type]?.enabled as boolean}
                disabled={channels[channel.type] == null}
                checkedChildren={t("settings.notification.channel.enabled.on")}
                unCheckedChildren={t("settings.notification.channel.enabled.off")}
                onChange={(checked) => handleSwitchChange(channel.type, checked)}
              />
            </div>
          ),
          forceRender: true,
        };
      }),
    [i18n.language, channels]
  );

  const handleSwitchChange = (channel: string, enabled: boolean) => {
    setChannel(channel, { enabled });
  };

  return (
    <div className={className} style={style}>
      {!initialized ? (
        <Skeleton active />
      ) : (
        <Collapse className={classNames?.collapse} style={styles?.collapse} accordion={true} bordered={false} items={channelCollapseItems} />
      )}
    </div>
  );
};

export default NotifyChannels;
