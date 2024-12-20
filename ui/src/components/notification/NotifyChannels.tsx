import { useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import { useDeepCompareMemo } from "@ant-design/pro-components";
import { Button, Collapse, message, notification, Skeleton, Space, Switch, Tooltip, type CollapseProps } from "antd";

import NotifyChannelEditForm, { type NotifyChannelEditFormInstance } from "./NotifyChannelEditForm";
import NotifyTestButton from "./NotifyTestButton";
import { notifyChannelsMap } from "@/domain/settings";
import { useNotifyChannelStore } from "@/stores/notify";
import { getErrMsg } from "@/utils/error";

type NotifyChannelsSemanticDOM = "collapse" | "form";

export type NotifyChannelsProps = {
  className?: string;
  classNames?: Partial<Record<NotifyChannelsSemanticDOM, string>>;
  style?: React.CSSProperties;
  styles?: Partial<Record<NotifyChannelsSemanticDOM, React.CSSProperties>>;
};

const NotifyChannels = ({ className, classNames, style, styles }: NotifyChannelsProps) => {
  const { t, i18n } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { initialized, channels, setChannel, fetchChannels } = useNotifyChannelStore();
  useEffect(() => {
    fetchChannels();
  }, [fetchChannels]);

  const channelFormRefs = useRef<Array<NotifyChannelEditFormInstance | null>>([]);
  const channelCollapseItems: CollapseProps["items"] = useDeepCompareMemo(
    () =>
      Array.from(notifyChannelsMap.values()).map((channel, index) => {
        return {
          key: `channel-${channel.type}`,
          label: <>{t(channel.name)}</>,
          children: (
            <div className={classNames?.form} style={styles?.form}>
              <NotifyChannelEditForm ref={(ref) => (channelFormRefs.current[index] = ref)} model={channels[channel.type]} channel={channel.type} />

              <Space>
                <Button type="primary" onClick={() => handleClickSubmit(channel.type, index)}>
                  {t("common.button.save")}
                </Button>

                {channels[channel.type] ? (
                  <Tooltip title={t("settings.notification.push_test.tooltip")}>
                    <>
                      <NotifyTestButton channel={channel.type} />
                    </>
                  </Tooltip>
                ) : null}
              </Space>
            </div>
          ),
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

  const handleClickSubmit = async (channel: string, index: number) => {
    const form = channelFormRefs.current[index];
    if (!form) {
      return;
    }

    await form.validateFields();

    try {
      setChannel(channel, form.getFieldsValue());

      messageApi.success(t("common.text.operation_succeeded"));
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    }
  };

  return (
    <div className={className} style={style}>
      {MessageContextHolder}
      {NotificationContextHolder}

      {!initialized ? (
        <Skeleton active />
      ) : (
        <Collapse className={classNames?.collapse} style={styles?.collapse} accordion={true} bordered={false} items={channelCollapseItems} />
      )}
    </div>
  );
};

export default NotifyChannels;
