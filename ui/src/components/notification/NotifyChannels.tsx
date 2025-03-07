import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareMemo } from "@ant-design/pro-components";
import { Button, Collapse, type CollapseProps, Skeleton, Space, Switch, message, notification } from "antd";

import Show from "@/components/Show";
import { notifyChannelsMap } from "@/domain/settings";
import { useZustandShallowSelector } from "@/hooks";
import { useNotifyChannelsStore } from "@/stores/notify";
import { getErrMsg } from "@/utils/error";

import NotifyChannelEditForm, { type NotifyChannelEditFormInstance } from "./NotifyChannelEditForm";
import NotifyTestButton from "./NotifyTestButton";

type NotifyChannelProps = {
  className?: string;
  style?: React.CSSProperties;
  channel: string;
};

const NotifyChannel = ({ className, style, channel }: NotifyChannelProps) => {
  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { channels, setChannel } = useNotifyChannelsStore(useZustandShallowSelector(["channels", "setChannel"]));

  const channelConfig = useDeepCompareMemo(() => channels[channel], [channels, channel]);
  const channelFormRef = useRef<NotifyChannelEditFormInstance>(null);
  const [channelFormChanged, setChannelFormChanged] = useState(false);

  const handleSubmit = async () => {
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

      <NotifyChannelEditForm
        ref={channelFormRef}
        className="mt-2"
        channel={channel}
        initialValues={channelConfig}
        onValuesChange={() => setChannelFormChanged(true)}
      />

      <Space className="mb-2">
        <Button type="primary" disabled={!channelFormChanged} onClick={handleSubmit}>
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

  const { channels, loadedAtOnce, setChannel, fetchChannels } = useNotifyChannelsStore(
    useZustandShallowSelector(["channels", "loadedAtOnce", "setChannel", "fetchChannels"])
  );
  useEffect(() => {
    fetchChannels();
  }, []);

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
                checkedChildren={t("settings.notification.channel.switch.on")}
                unCheckedChildren={t("settings.notification.channel.switch.off")}
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
      <Show when={loadedAtOnce} fallback={<Skeleton active />}>
        <Collapse className={classNames?.collapse} style={styles?.collapse} accordion bordered={false} items={channelCollapseItems} />
      </Show>
    </div>
  );
};

export default NotifyChannels;
