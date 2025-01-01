import { useTranslation } from "react-i18next";
import { useRequest } from "ahooks";
import { Button, message, notification, type ButtonProps } from "antd";

import { notifyTest } from "@/api/notify";
import { getErrMsg } from "@/utils/error";

export type NotifyTestButtonProps = {
  className?: string;
  style?: React.CSSProperties;
  channel: string;
  disabled?: boolean;
  size?: ButtonProps["size"];
};

const NotifyTestButton = ({ className, style, channel, disabled, size }: NotifyTestButtonProps) => {
  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { loading, run: executeNotifyTest } = useRequest(
    () => {
      return notifyTest(channel);
    },
    {
      refreshDeps: [channel],
      manual: true,
      onSuccess: () => {
        messageApi.success(t("settings.notification.push_test.pushed"));
      },
      onError: (err) => {
        notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
      },
    }
  );

  const handleClick = () => {
    executeNotifyTest();
  };

  return (
    <>
      {MessageContextHolder}
      {NotificationContextHolder}

      <Button className={className} style={style} disabled={disabled} loading={loading} size={size} onClick={handleClick}>
        {t("settings.notification.push_test.button")}
      </Button>
    </>
  );
};

export default NotifyTestButton;
