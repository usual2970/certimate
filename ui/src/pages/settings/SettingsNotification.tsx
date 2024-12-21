import { useTranslation } from "react-i18next";
import { Card, Divider } from "antd";

import NotifyChannels from "@/components/notification/NotifyChannels";
import NotifyTemplate from "@/components/notification/NotifyTemplate";
import { useNotifyChannelStore } from "@/stores/notify";

const SettingsNotification = () => {
  const { t } = useTranslation();

  const { initialized } = useNotifyChannelStore();

  return (
    <div>
      <Card className="shadow" title={t("settings.notification.template.card.title")}>
        <div className="md:max-w-[40rem]">
          <NotifyTemplate />
        </div>
      </Card>

      <Divider />

      <Card className="shadow" styles={{ body: initialized ? { padding: 0 } : {} }} title={t("settings.notification.channels.card.title")}>
        <NotifyChannels classNames={{ form: "md:max-w-[40rem]" }} />
      </Card>
    </div>
  );
};

export default SettingsNotification;
