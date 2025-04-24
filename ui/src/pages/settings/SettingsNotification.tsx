import { useTranslation } from "react-i18next";
import { Alert, Card, Divider } from "antd";

import NotifyChannels from "@/components/notification/NotifyChannels";
import NotifyTemplate from "@/components/notification/NotifyTemplate";
import { useZustandShallowSelector } from "@/hooks";
import { useNotifyChannelsStore } from "@/stores/notify";

/**
 * @deprecated
 */
const SettingsNotification = () => {
  const { t } = useTranslation();

  const { loadedAtOnce } = useNotifyChannelsStore(useZustandShallowSelector(["loadedAtOnce"]));

  return (
    <div>
      <Card className="shadow" title={t("settings.notification.template.card.title")}>
        <div className="md:max-w-[40rem]">
          <NotifyTemplate />
        </div>
      </Card>

      <Divider />

      <Card className="shadow" styles={{ body: loadedAtOnce ? { padding: 0 } : {} }} title={t("settings.notification.channels.card.title")}>
        <Alert type="warning" banner message="本页面相关功能即将在后续版本中废弃，请使用「授权管理」页面来管理通知渠道。" />
        <NotifyChannels classNames={{ form: "md:max-w-[40rem]" }} />
      </Card>
    </div>
  );
};

export default SettingsNotification;
